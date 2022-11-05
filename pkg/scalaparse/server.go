package scalaparse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/amenzhinsky/go-memexec"
	"github.com/bazelbuild/rules_go/go/tools/bazel"

	sppb "github.com/stackb/scala-gazelle/api/scalaparse"
)

const contentTypeJSON = "application/json"

type ScalaParseClient struct {
}

func NewScalaParseServer() *ScalaParseServer {
	return &ScalaParseServer{}
}

type ScalaParseServer struct {
	sppb.UnimplementedScalaParserServer

	process    *memexec.Exec
	processDir string
	cmd        *exec.Cmd

	httpClient *http.Client
	httpUrl    string

	HttpPort int
}

func (s *ScalaParseServer) Stop() {
	if s.httpClient != nil {
		s.httpClient.CloseIdleConnections()
		s.httpClient = nil
	}
	if s.cmd != nil {
		s.cmd.Process.Kill()
		s.cmd = nil
	}
	if s.process != nil {
		s.process.Close()
		s.process = nil
	}
	if s.processDir != "" {
		os.RemoveAll(s.processDir)
		s.processDir = ""
	}
}

func (s *ScalaParseServer) Start() error {
	//
	// Setup temp process directory and write js files
	//
	processDir, err := bazel.NewTmpDir("")
	if err != nil {
		return err
	}

	scriptPath := filepath.Join(processDir, "sourceindexer.mjs")
	parserPath := filepath.Join(processDir, "node_modules", "scalameta-parsers", "index.js")

	if err := os.MkdirAll(filepath.Dir(parserPath), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(scriptPath, []byte(sourceindexerMjs), os.ModePerm); err != nil {
		return err
	}
	if err := ioutil.WriteFile(parserPath, []byte(scalametaParsersIndexJs), os.ModePerm); err != nil {
		return err
	}

	if debugParse {
		listFiles(".")
	}

	//
	// ensure we have a port
	//
	if s.HttpPort == 0 {
		port, err := getFreePort()
		if err != nil {
			return status.Errorf(codes.FailedPrecondition, "getting http port: %v", err)
		}
		s.HttpPort = port
	}
	s.httpUrl = fmt.Sprintf("http://127.0.0.1:%d", s.HttpPort)

	//
	// Setup the bun process
	//
	exe, err := memexec.New(nodeExe)
	if err != nil {
		return err
	}
	s.process = exe

	//
	// Start the bun process
	//
	cmd := exe.Command("sourceindexer.mjs")
	cmd.Dir = processDir
	cmd.Env = []string{
		"NODE_PATH=" + processDir,
		fmt.Sprintf("PORT=%d", s.HttpPort),
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	s.cmd = cmd

	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		// does it make sense to wait for the process?  We kill it forcefully
		// at the end anyway...
		if err := cmd.Wait(); err != nil {
			if err.Error() != "signal: killed" {
				log.Printf("command wait err: %v", err)
			}
		}
	}()

	waitForConnectionAvailable("localhost", s.HttpPort, 3*time.Second)

	//
	// Setup the http client
	//
	s.httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
	}

	return nil
}

func (s *ScalaParseServer) Parse(ctx context.Context, in *sppb.ScalaParseRequest) (*sppb.ScalaParseResponse, error) {
	req, err := newHttpScalaParseRequest(s.httpUrl, in)
	if err != nil {
		return nil, err
	}
	w, err := s.httpClient.Do(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response error: %v", err)
	}

	if debugParse {
		respDump, err := httputil.DumpResponse(w, true)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("HTTP_RESPONSE:\n%s", string(respDump))
	}

	contentType := w.Header.Get("Content-Type")
	if contentType != contentTypeJSON {
		return nil, status.Errorf(codes.Internal, "response content-type error, want %q, got: %q", contentTypeJSON, contentType)
	}

	data, err := ioutil.ReadAll(w.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "response data error: %v", err)
	}

	if debugParse {
		log.Printf("response body: %s", string(data))
	}
	var response sppb.ScalaParseResponse

	if err := protojson.Unmarshal(data, &response); err != nil {
		return nil, status.Errorf(codes.Internal, "response body error: %v\n%s", err, string(data))
	}

	return &response, nil
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

func newHttpScalaParseRequest(url string, in *sppb.ScalaParseRequest) (*http.Request, error) {
	if url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "request URL is required")
	}
	if in == nil {
		return nil, status.Errorf(codes.InvalidArgument, "ScalaParseRequest is required")
	}
	values := map[string]interface{}{"label": in.Label, "files": in.Filename}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// waitForConnectionAvailable pings a tcp connection every 250 milliseconds
// until it connects and returns true.  If it fails to connect by the timeout
// deadline, returns false.
func waitForConnectionAvailable(host string, port int, timeout time.Duration) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	var wg sync.WaitGroup
	wg.Add(1)
	then := time.Now()

	success := make(chan bool, 1)

	go func() {
		go func() {
			defer wg.Done()
			for {
				_, err := net.Dial("tcp", target)
				if err == nil {
					if debugParse {
						log.Printf("%s is available after %s", target, time.Since(then))
					}
					break
				}
				time.Sleep(250 * time.Millisecond)
			}
		}()
		wg.Wait()
		success <- true
	}()

	select {
	case <-success:
		return true
	case <-time.After(timeout):
		return false
	}
}