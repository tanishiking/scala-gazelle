package scala

import (
	"flag"
	"log"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	ScalaLangName = "scala"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	return &scalaLang{
		ruleRegistry:          globalRegistry,
		crossResolverRegistry: globalCrossResolverRegistry,
		packages:              make(map[string]*scalaPackage),
	}
}

// scalaLang implements language.Language.
type scalaLang struct {
	// classIndexFile is the filename used for the class index
	classIndexFile string
	// ruleRegistry is the rule registry implementation
	ruleRegistry RuleRegistry
	// crossResolverRegistry is the cross resolver registry implementation
	crossResolverRegistry CrossResolverRegistry
	// packages is map from the config.Rel to *scalaPackage for the
	// workspace-relative packate name.
	packages map[string]*scalaPackage
}

// Name returns the name of the language. This should be a prefix of the kinds
// of rules generated by the language, e.g., "go" for the Go extension since it
// generates "go_library" rules.
func (sl *scalaLang) Name() string { return ScalaLangName }

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (sl *scalaLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	getOrCreateScalaConfig(c) // ignoring return value, only want side-effect

	for _, name := range sl.crossResolverRegistry.CrossResolverNames() {
		if resolver, err := sl.crossResolverRegistry.LookupCrossResolver(name); err == nil {
			resolver.RegisterFlags(fs, cmd, c)
		}
	}
}

func (sl *scalaLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	for _, name := range sl.crossResolverRegistry.CrossResolverNames() {
		if resolver, err := sl.crossResolverRegistry.LookupCrossResolver(name); err == nil {
			if err := resolver.CheckFlags(fs, c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (*scalaLang) KnownDirectives() []string {
	return []string{
		ruleDirective,
	}
}

// Configure implements config.Configurer
func (sl *scalaLang) Configure(c *config.Config, rel string, f *rule.File) {
	if f == nil {
		return
	}
	if err := getOrCreateScalaConfig(c).ParseDirectives(rel, f.Directives); err != nil {
		log.Fatalf("error while parsing rule directives in package %q: %v", rel, err)
	}
}

// Kinds returns a map of maps rule names (kinds) and information on how to
// match and merge attributes that may be found in rules of those kinds. All
// kinds of rules generated for this language may be found here.
func (sl *scalaLang) Kinds() map[string]rule.KindInfo {
	kinds := make(map[string]rule.KindInfo)

	for _, name := range sl.ruleRegistry.RuleNames() {
		rule, err := sl.ruleRegistry.LookupRule(name)
		if err != nil {
			log.Fatal("Kinds:", err)
		}
		kinds[rule.Name()] = rule.KindInfo()
	}

	return kinds
}

// Loads returns .bzl files and symbols they define. Every rule generated by
// GenerateRules, now or in the past, should be loadable from one of these
// files.
func (sl *scalaLang) Loads() []rule.LoadInfo {
	// Merge symbols
	symbolsByLoadName := make(map[string][]string)

	for _, name := range sl.ruleRegistry.RuleNames() {
		rule, err := sl.ruleRegistry.LookupRule(name)
		if err != nil {
			log.Fatal(err)
		}
		load := rule.LoadInfo()
		symbolsByLoadName[load.Name] = append(symbolsByLoadName[load.Name], load.Symbols...)
	}

	// Ensure names are sorted otherwise order of load statements can be
	// non-deterministic
	keys := make([]string, 0)
	for name := range symbolsByLoadName {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// Build final load list
	loads := make([]rule.LoadInfo, 0)
	for _, name := range keys {
		symbols := symbolsByLoadName[name]
		sort.Strings(symbols)
		loads = append(loads, rule.LoadInfo{
			Name:    name,
			Symbols: symbols,
		})
	}
	return loads
}

// Fix repairs deprecated usage of language-specific rules in f. This is called
// before the file is indexed. Unless c.ShouldFix is true, fixes that delete or
// rename rules should not be performed.
func (sl *scalaLang) Fix(c *config.Config, f *rule.File) {
}

// Imports returns a list of ImportSpecs that can be used to import the rule r.
// This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
func (sl *scalaLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	from := label.New("", f.Pkg, r.Name())

	pkg, ok := sl.packages[from.Pkg]
	if !ok {
		// log.Println("scala.Imports(): Unknown package", from.Pkg)
		return nil
	}

	provider := pkg.ruleProvider(r)
	if provider == nil {
		log.Printf("Unknown rule provider for //%s:%s %p", f.Pkg, r.Name(), r)
		return nil
	}

	return provider.Imports(c, r, f)
}

// Embeds returns a list of labels of rules that the given rule embeds. If a
// rule is embedded by another importable rule of the same language, only the
// embedding rule will be indexed. The embedding rule will inherit the imports
// of the embedded rule. Since SkyLark doesn't support embedding this should
// always return nil.
func (*scalaLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each rule
// generated by language.GenerateRules in language.GenerateResult.Imports.
// Resolve generates a "deps" attribute (or the appropriate language-specific
// equivalent) for each import according to language-specific rules and
// heuristics.
func (sl *scalaLang) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	importsRaw interface{},
	from label.Label,
) {
	if pkg, ok := sl.packages[from.Pkg]; ok {
		provider := pkg.ruleProvider(r)
		if provider == nil {
			log.Printf("no known rule provider for %v", from)
		}
		if imports, ok := importsRaw.([]string); ok {
			provider.Resolve(c, ix, r, imports, from)
		} else {
			log.Printf("warning: resolve scala imports: expected []string, got %T", importsRaw)
		}
	} else {
		log.Printf("no known rule package for %v", from.Pkg)
	}
}

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested in
// depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a struct to
// avoid breaking implementations in the future when new fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
func (sl *scalaLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	log.Println("visiting", args.Rel)

	cfg := getOrCreateScalaConfig(args.Config)

	files := make([]*ScalaFile, 0)

	// for _, f := range args.RegularFiles {
	// 	if !isScalaFile(f) {
	// 		continue
	// 	}
	// 	file, err := ParseScalaFile(args.Dir, f)
	// 	if err != nil {
	// 		log.Println("error parsing scala file:", f, err.Error())
	// 		continue
	// 	}
	// 	files = append(files, file)
	// }

	pkg := newScalaPackage(sl.ruleRegistry, args.File, cfg, files...)
	sl.packages[args.Rel] = pkg

	rules := pkg.Rules()
	empty := pkg.Empty()

	imports := make([]interface{}, len(rules))
	for i, r := range rules {
		imports[i] = r.PrivateAttr(config.GazelleImportsKey)
		// internalLabel := label.New("", args.Rel, r.Name())
		// protoc.GlobalRuleIndex().Put(internalLabel, r)
	}

	if args.Rel == "" {
		for _, name := range sl.crossResolverRegistry.CrossResolverNames() {
			// log.Println("cross resolve", name, lang, imp.Imp)
			if resolver, err := sl.crossResolverRegistry.LookupCrossResolver(name); err == nil {
				if ssr, ok := resolver.(*scalaSourceIndexResolver); ok {
					if err := ssr.DumpIndex("/tmp/ssr.json"); err != nil {
						log.Println("dump index error:", err)
					}
				}
			}
		}
	}

	return language.GenerateResult{
		Gen:     rules,
		Empty:   empty,
		Imports: imports,
	}
}

// CrossResolve calls all known resolvers and returns the first non-empty result.
func (sl *scalaLang) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	for _, name := range sl.crossResolverRegistry.CrossResolverNames() {
		// log.Println("cross resolve", name, lang, imp.Imp)
		if resolver, err := sl.crossResolverRegistry.LookupCrossResolver(name); err == nil {
			if result := resolver.CrossResolve(c, ix, imp, lang); len(result) > 0 {
				return result
			}
		}
	}
	return nil
}

func fullyQualifiedRuleKind(loads []*rule.Load, ruleKind string) string {
	for _, load := range loads {
		for _, pair := range load.SymbolPairs() {
			if pair.To == ruleKind {
				return load.Name() + "%" + pair.To
			}
		}
	}
	return ruleKind
}
