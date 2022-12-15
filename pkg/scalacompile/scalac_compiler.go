package scalacompile

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sppb "github.com/stackb/scala-gazelle/build/stack/gazelle/scala/parse"
	"github.com/stackb/scala-gazelle/pkg/collections"
	"github.com/stackb/scala-gazelle/pkg/resolver"
)

const debugCompiler = true

// NOT_FOUND is the diagnostic message prefix we expect from the scala compiler
// when it can't resolve types or values.
const NOT_FOUND = "not found: "

var notPackageMemberRe = regexp.MustCompile(`^object ([^ ]+) is not a member of package (.*)$`)

func NewScalacCompiler() *ScalacCompilerService {
	return &ScalacCompilerService{}
}

// ScalacCompilerService implements a scala compiler frontend that communicates
// with a long-running scalac backend process over gRPC.
type ScalacCompilerService struct {
	backendHost        string
	backendPort        int
	backendUrl         string
	backendDialTimeout time.Duration

	// repoRoot is typically the config.Config.RepoRoot
	repoRoot string
	// scalacserverJarPath is the unresolved runfile
	scalacserverJarPath string
	// javaBinPath is the path to the java interpreter
	javaBinPath string

	grpcConn *grpc.ClientConn
	client   sppb.CompilerClient

	// the process
	cmd *exec.Cmd
}

// Name implements part of the provider.SymbolProvider interface.
func (p *ScalacCompilerService) Name() string {
	return "scalac"
}

// RegisterFlags implements part of the provider.SymbolProvider interface.
func (p *ScalacCompilerService) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	fs.StringVar(&p.scalacserverJarPath, "scalac_jar_path", "", "filesystem path to the scala compiler server tool jar")
	fs.StringVar(&p.javaBinPath, "scalac_java_bin_path", "", "filesystem path to the java tool $(location @local_jdk//:bin/java)")
	fs.StringVar(&p.backendHost, "scalac_backend_host", "localhost", "bind host for the backend server")
	fs.IntVar(&p.backendPort, "scalac_backend_port", 0, "bind port for the backend server")
	fs.DurationVar(&p.backendDialTimeout, "scalac_backend_dial_timeout", time.Second*3, "compiler backend dial timeout")
}

// CheckFlags implements part of the Configurer interface.
func (p *ScalacCompilerService) CheckFlags(fs *flag.FlagSet, c *config.Config, scope resolver.Scope) error {
	p.repoRoot = c.RepoRoot
	p.javaBinPath = os.ExpandEnv(p.javaBinPath)
	p.scalacserverJarPath = os.ExpandEnv(p.scalacserverJarPath)

	// start is disabled if the backendPort is already set.
	if p.backendPort == 0 {
		port, err := getFreePort()
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "getting http port: %v", err)
		}
		p.backendPort = port

		if err := p.start(); err != nil {
			return err
		}
	}

	return p.startGrpcClient()
}

// CanProvide implements the resolver.SymbolProvider interface.
func (p *ScalacCompilerService) CanProvide(from label.Label, ruleIndex func(from label.Label) (*rule.Rule, bool)) bool {
	return false
}

func (s *ScalacCompilerService) OnResolve() error {
	return nil
}

func (s *ScalacCompilerService) OnEnd() error {
	if s.cmd != nil {
		if err := s.cmd.Process.Kill(); err != nil {
			return err
		}
	}
	if s.grpcConn != nil {
		return s.grpcConn.Close()
	}
	return nil
}

func (s *ScalacCompilerService) start() error {
	t1 := time.Now()

	//
	// Start the backend process
	//
	cmd := exec.Command(s.javaBinPath,
		fmt.Sprintf("-Dscalac.server.port=%d", s.backendPort),
		"-jar", s.scalacserverJarPath,
	)

	// cmd.Dir = s.repoRoot
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	s.cmd = cmd

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting process %s: %w", s.scalacserverJarPath, err)
	}
	go func() {
		// FIXME(pcj): does it make sense to wait for the process?  We kill it
		// forcefully at the end anyway...
		if err := cmd.Wait(); err != nil {
			if err.Error() != "signal: killed" {
				log.Printf("command wait err: %v", err)
			}
		}
	}()

	t2 := time.Since(t1).Round(1 * time.Millisecond)
	log.Printf("Compiler started (%v)", t2)

	return nil
}

func (s *ScalacCompilerService) startGrpcClient() error {
	s.backendUrl = fmt.Sprintf("grpc://%s:%d", s.backendHost, s.backendPort)

	target := fmt.Sprintf("%s:%d", s.backendHost, s.backendPort)
	conn, err := grpc.Dial(target,
		grpc.WithInsecure(),
		// grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	s.grpcConn = conn
	s.client = sppb.NewCompilerClient(conn)

	if !collections.WaitForConnectionAvailable(s.backendHost, s.backendPort, s.backendDialTimeout) {
		return fmt.Errorf("failed to connect to scalac backend %s in %v", s.backendUrl, s.backendDialTimeout)
	}

	return nil
}

// CompileScalaFiles implements scalacompile.Compiler
func (p *ScalacCompilerService) CompileScalaFiles(from label.Label, dir string, srcs ...string) ([]*sppb.File, error) {
	t1 := time.Now()

	files := make([]*sppb.File, len(srcs))
	fileMap := make(map[string]int)
	for i, src := range srcs {
		files[i] = &sppb.File{Filename: src}
		fileMap[src] = i
	}

	resp, err := p.client.Compile(context.Background(), &sppb.CompileRequest{
		Dir:       dir,
		Filenames: srcs,
	})
	if err != nil {
		return nil, fmt.Errorf("compiler backend error: %w", err)
	}

	seen := make(map[string]bool)
	for _, d := range resp.Diagnostics {
		if d.Source == "" || d.Source == "<no file>" {
			if false {
				log.Printf("skipping diagnostic: %v (no file)", d)
			}
			continue
		}
		// FIXME(pcj): dedup in backend?
		key := fmt.Sprintf("%s:%v:%s", d.Source, d.Severity, d.Message)
		if seen[key] {
			if false {
				log.Printf("skipping diagnostic: %v (seen)", d)
			}
			continue
		}
		seen[key] = true

		log.Printf("diagnostic: %+v", d)
		symbol, ok := symbolFromDiagnostic(d)
		if !ok {
			log.Printf("symbol skipped (not ok): %+v", symbol)

			continue
		}
		log.Printf("symbol: %+v", symbol)
		fileIndex, ok := fileMap[d.Source]
		if !ok {
			log.Panicln("cannot find file", d.Source, "in", srcs)
		}
		file := files[fileIndex]
		file.Symbols = append(file.Symbols, symbol)
	}

	t2 := time.Since(t1).Round(1 * time.Millisecond)
	if false {
		log.Printf("Compile %s (%d files, %v)", from, len(files), t2)
	}
	return files, nil
}

func symbolFromDiagnostic(d *sppb.Diagnostic) (*sppb.Symbol, bool) {
	switch d.Severity {
	case sppb.Severity_ERROR:
		return symbolFromErrorDiagnostic(d)
	default:
		return nil, false
	}
}

func symbolFromErrorDiagnostic(d *sppb.Diagnostic) (*sppb.Symbol, bool) {
	if strings.HasPrefix(d.Message, NOT_FOUND) {
		return symbolFromNotFoundErrorDiagnostic(d.Message[len(NOT_FOUND):])
	} else if match := notPackageMemberRe.FindStringSubmatch(d.Message); match != nil {
		return symbolFromNotPackageMemberErrorDiagnostic(match[1], match[2])
	}
	return nil, false
}

func symbolFromNotFoundErrorDiagnostic(msg string) (*sppb.Symbol, bool) {
	fields := strings.Fields(msg)
	if len(fields) < 2 {
		return nil, false
	}
	return &sppb.Symbol{
		Type: parseSymbolType(fields[0]),
		Name: fields[1],
	}, true
}

func symbolFromNotPackageMemberErrorDiagnostic(obj, pkg string) (*sppb.Symbol, bool) {
	return &sppb.Symbol{
		Type: sppb.SymbolType_SYMBOL_PACKAGE,
		Name: pkg + "?" + obj,
	}, true
	// spec.NotMember = append(spec.NotMember, &NotMemberSymbol{Kind: "object", Name: obj, Package: pkg})
}

func parseSymbolType(val string) sppb.SymbolType {
	switch val {
	case "object":
		return sppb.SymbolType_SYMBOL_OBJECT
	case "type":
		return sppb.SymbolType_SYMBOL_TYPE
	case "value":
		return sppb.SymbolType_SYMBOL_VALUE
	default:
		log.Panicf("unknown symbol type: %q", val)
		return sppb.SymbolType_SYMBOL_TYPE_UNKNOWN
	}
}

// getFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
