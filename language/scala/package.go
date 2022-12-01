package scala

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/scala-gazelle/pkg/scalaparse"
)

const (
	ruleProviderKey = "_scala_rule_provider"
)

type ScalaPackage interface {
	// Rel returns the relative path to this package
	Rel() string
	// Dir returns the absolute path to the worksace
	Dir() string
	// File returns the BUILD file for the package
	File() *rule.File
	// ScalaParser returns the parser instance to use.
	ScalaParser() scalaparse.Parser
}

// scalaPackage provides a set of proto_library derived rules for the package.
type scalaPackage struct {
	// parser is the file parser
	parser scalaparse.Parser
	// rel is the package (args.Rel)
	rel string
	// the registry to use
	ruleRegistry RuleRegistry
	// the build file
	file *rule.File
	// the config for this package
	cfg *scalaConfig
	// the generated and empty rule providers
	gen, empty []RuleProvider
	// parent is the parent package
	parent *scalaPackage
	// resolved is the final state of generated rules, by name.
	rules map[string]*rule.Rule
}

// newScalaPackage constructs a Package given a list of scala files.
func newScalaPackage(ruleRegistry RuleRegistry, parser scalaparse.Parser, rel string, file *rule.File, cfg *scalaConfig) *scalaPackage {
	s := &scalaPackage{
		parser:       parser,
		rel:          rel,
		ruleRegistry: ruleRegistry,
		file:         file,
		cfg:          cfg,
		rules:        make(map[string]*rule.Rule),
	}
	s.gen = s.generateRules(true)
	// s.empty = s.generateRules(false)

	return s
}

// Config returns the the underlying config.
func (s *scalaPackage) Config() *scalaConfig {
	return s.cfg
}

// getRule returns the named rule, if it exists
func (s *scalaPackage) getRule(name string) (*rule.Rule, bool) {
	got, ok := s.rules[name]
	return got, ok
}

// ruleProvider returns the provider of a rule or nil if not known.
func (s *scalaPackage) ruleProvider(r *rule.Rule) RuleProvider {
	if provider, ok := r.PrivateAttr(ruleProviderKey).(RuleProvider); ok {
		return provider
	}
	return nil
}

func (s *scalaPackage) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	importsRaw interface{},
	from label.Label,
) {
	provider := s.ruleProvider(r)
	if provider == nil {
		log.Printf("no known rule provider for %v", from)
		return
	}
	provider.Resolve(c, ix, r, importsRaw, from)
	s.rules[r.Name()] = r
}

// generateRules constructs a list of rules based on the configured set of rule
// configurations.
func (s *scalaPackage) generateRules(enabled bool) []RuleProvider {
	rules := make([]RuleProvider, 0)

	existingRulesByFQN := make(map[string][]*rule.Rule)
	if s.file != nil {
		for _, r := range s.file.Rules {
			fqn := FullyQualifiedName(s.file.Loads, r.Kind())
			existingRulesByFQN[fqn] = append(existingRulesByFQN[fqn], r)
		}
	}

	for _, rc := range s.cfg.configuredRules() {
		// if enabled != rc.Enabled {
		if !rc.Enabled {
			// log.Printf("%s: skipping rule config %s (not enabled)", s.rel, rc.Name)
			continue
		}
		rule := s.provideRule(rc)
		if rule != nil {
			rules = append(rules, rule)
		}
		existing := existingRulesByFQN[rc.Implementation]
		if len(existing) > 0 {
			for _, r := range existing {
				rule := s.resolveRule(rc, r)
				if rule != nil {
					rules = append(rules, rule)
				}
			}
		}
		delete(existingRulesByFQN, rc.Implementation)
	}

	return rules
}

func (s *scalaPackage) provideRule(rc *RuleConfig) RuleProvider {
	impl, err := s.ruleRegistry.LookupRule(rc.Implementation)
	if err == ErrUnknownRule {
		log.Fatalf(
			"%s: rule not registered: %q (available: %v)",
			s.Rel(),
			rc.Implementation,
			s.ruleRegistry.RuleNames(),
		)
	}
	rc.Impl = impl

	return impl.ProvideRule(rc, s)
}

func (s *scalaPackage) resolveRule(rc *RuleConfig, r *rule.Rule) RuleProvider {
	impl, err := s.ruleRegistry.LookupRule(rc.Implementation)
	if err == ErrUnknownRule {
		log.Fatalf(
			"%s: rule not registered: %q (available: %v)",
			s.Rel(),
			rc.Implementation,
			globalRuleRegistry.RuleNames(),
		)
	}
	rc.Impl = impl

	if rr, ok := impl.(RuleResolver); ok {
		return rr.ResolveRule(rc, s, r)
	}

	return nil
}

// ScalaParser implements part of the ScalaPackage interface.
func (s *scalaPackage) ScalaParser() scalaparse.Parser {
	return s.parser
}

// File implements part of the ScalaPackage interface.
func (s *scalaPackage) File() *rule.File {
	return s.file
}

// Rel implements part of the ScalaPackage interface.
func (s *scalaPackage) Rel() string {
	return s.rel
}

// Dir implements part of the ScalaPackage interface.
func (s *scalaPackage) Dir() string {
	return s.cfg.config.RepoRoot
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
			// log.Println("provided rule %s %s", r.Kind(), r.Name())
			s.rules[r.Name()] = r
		}
		rules = append(rules, r)
	}
	return rules
}

func FullyQualifiedName(loads []*rule.Load, kind string) string {
	for _, load := range loads {
		for _, pair := range load.SymbolPairs() {
			// when there is no aliasing, pair.From == pair.To, so this covers
			// both cases (aliases and not).
			if pair.From == pair.To && pair.From == kind {
				return load.Name() + "%" + pair.From
			}
			if pair.To == kind {
				return load.Name() + "%" + pair.From
			}
		}
	}
	// no match, just return the kind (e.g. native.java_binary)
	return kind
}
