package scala

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ruleProviderKey = "_scala_rule_provider"
)

type ScalaPackage interface {
	Rel() string
	File() *rule.File
	// Files() []*ScalaFile
}

// scalaPackage provides a set of proto_library derived rules for the package.
type scalaPackage struct {
	// the registry to use
	ruleRegistry RuleRegistry
	// the build file
	file *rule.File
	// the config for this package
	cfg *scalaConfig
	// the list of '.scala' files
	files []*ScalaFile
	// the generated and empty rule providers
	gen, empty []RuleProvider
}

// newScalaPackage constructs a Package given a list of scala files.
func newScalaPackage(ruleRegistry RuleRegistry, file *rule.File, cfg *scalaConfig, files ...*ScalaFile) *scalaPackage {
	s := &scalaPackage{
		ruleRegistry: ruleRegistry,
		file:         file,
		cfg:          cfg,
		files:        files,
	}
	s.gen = s.generateRules(true)
	s.empty = s.generateRules(false)

	return s
}

// ruleProvider returns the provider of a rule or nil if not known.
func (s *scalaPackage) ruleProvider(r *rule.Rule) RuleProvider {
	if provider, ok := r.PrivateAttr(ruleProviderKey).(RuleProvider); ok {
		return provider
	}
	return nil
}

// generateRules constructs a list of rules based on the configured set of
// languages.
func (s *scalaPackage) generateRules(enabled bool) []RuleProvider {
	rules := make([]RuleProvider, 0)

	if s.file != nil {
		for _, r := range s.file.Rules {
			// fqrk := fullyQualifiedRuleKind(args.File.Loads, r.Kind())
			rc, ok := s.cfg.GetConfiguredRule(r.Kind())
			if !ok {
				continue
			}

			log.Println("matched resolver rule", r.Kind(), r.Name())

			rule := s.resolveRule(rc, r)
			if rule == nil {
				continue
			}
			rules = append(rules, rule)

			// if strings.HasPrefix(r.Kind(), "scala_") {
			// 	log.Printf("Existing rule %s %s", r.Kind(), r.Name())
			// }
		}
	}

	for _, rc := range s.cfg.configuredRules() {
		if enabled != rc.Enabled {
			continue
		}
		rule := s.provideRule(rc)
		if rule == nil {
			continue
		}
		rules = append(rules, rule)
	}

	return rules
}

func (s *scalaPackage) provideRule(rc *RuleConfig) RuleProvider {
	impl, err := globalRegistry.LookupRule(rc.Implementation)
	if err == ErrUnknownRule {
		log.Fatalf(
			"%s: rule not registered: %q (available: %v)",
			s.Rel(),
			rc.Implementation,
			globalRegistry.RuleNames(),
		)
	}
	rc.Impl = impl

	rule := impl.ProvideRule(rc, s)
	if rule == nil {
		return nil
	}

	return rule
}

func (s *scalaPackage) resolveRule(rc *RuleConfig, r *rule.Rule) RuleProvider {
	impl, err := globalRegistry.LookupRule(rc.Implementation)
	if err == ErrUnknownRule {
		log.Fatalf(
			"%s: rule not registered: %q (available: %v)",
			s.Rel(),
			rc.Implementation,
			globalRegistry.RuleNames(),
		)
	}
	rc.Impl = impl

	if rr, ok := impl.(RuleResolver); ok {
		return rr.ResolveRule(rc, s, r)
	}

	return nil
}

// File implements part of the ScalaPackage interface.
func (s *scalaPackage) File() *rule.File {
	return s.file
}

// Rel implements part of the ScalaPackage interface.
func (s *scalaPackage) Rel() string {
	var rel string
	if s.file != nil {
		rel = s.file.Pkg
	}
	return rel
}

// Files implements part of the ScalaPackage interface.
func (s *scalaPackage) Files() []*ScalaFile {
	return s.files
}

// Rules provides the aggregated rule list for the package.
func (s *scalaPackage) Rules() []*rule.Rule {
	return s.getProvidedRules(s.gen, true)
}

// Empty names the rules that can be deleted.
func (s *scalaPackage) Empty() []*rule.Rule {
	// it's a bit sad that we construct the full rules only for their kind and
	// name, but that's how it is right now.
	rules := s.getProvidedRules(s.empty, false)

	empty := make([]*rule.Rule, len(rules))
	for i, r := range rules {
		empty[i] = rule.NewRule(r.Kind(), r.Name())
	}

	return empty
}

func (s *scalaPackage) getProvidedRules(providers []RuleProvider, shouldResolve bool) []*rule.Rule {
	rules := make([]*rule.Rule, 0)
	for _, p := range providers {
		r := p.Rule()
		if r == nil {
			continue
		}

		if shouldResolve {
			// record the association of the rule provider here for the resolver.
			r.SetPrivateAttr(ruleProviderKey, p)

			// imports := r.PrivateAttr(config.GazelleImportsKey)
			// if imports == nil {
			// 	lib := s.ruleLibs[p]
			// 	r.SetPrivateAttr(ProtoLibraryKey, lib)
			// }

			// NOTE: this is a bit of a hack: it would be preferable to populate
			// the global resolver with import specs during the .Imports()
			// function.  One would think that the RuleProvider could be set as
			// a PrivateAttr to be retrieved in the Imports() function. However,
			// the rule ref seems to have changed by that time, the PrivateAttr
			// is removed.  Maybe this is due to rule merges?  Very difficult to
			// track down bug that cost me days.
			from := label.New("", s.Rel(), r.Name())
			file := rule.EmptyFile("", s.Rel())
			provideResolverImportSpecs(s.cfg.config, p, r, file, from)
		}

		rules = append(rules, r)
	}
	return rules
}

func provideResolverImportSpecs(c *config.Config, provider RuleProvider, r *rule.Rule, f *rule.File, from label.Label) {
	for _, imp := range provider.Imports(c, r, f) {
		protoc.GlobalResolver().Provide(
			"scala",
			imp.Lang,
			imp.Imp,
			from,
		)
	}
}
