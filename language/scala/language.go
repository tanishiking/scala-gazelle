package scala

import (
	"os"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/pcj/mobyprogress"
	"github.com/stackb/rules_proto/pkg/protoc"

	scpb "github.com/stackb/scala-gazelle/build/stack/gazelle/scala/cache"
	"github.com/stackb/scala-gazelle/pkg/crossresolve"
)

const scalaLangName = "scala"

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	packages := make(map[string]*scalaPackage)
	scalaCompiler := newScalaCompiler()

	sourceResolver := crossresolve.NewScalaSourceCrossResolver(scalaLangName)
	protoResolver := crossresolve.NewProtoResolver(scalaLangName, protoc.GlobalResolver().Provided)
	mavenResolver := crossresolve.NewMavenResolver("java")
	jarResolver := crossresolve.NewJarIndexCrossResolver(scalaLangName)

	crossresolve.Resolvers().MustRegisterResolver("source", sourceResolver)
	crossresolve.Resolvers().MustRegisterResolver("maven", mavenResolver)
	crossresolve.Resolvers().MustRegisterResolver("proto", protoResolver)
	crossresolve.Resolvers().MustRegisterResolver("jar", jarResolver)

	return &scalaLang{
		cache:          &scpb.Cache{},
		ruleRegistry:   globalRuleRegistry,
		sourceResolver: sourceResolver,
		scalaCompiler:  scalaCompiler,
		packages:       packages,
		resolvers:      make(map[string]crossresolve.ConfigurableCrossResolver),
		progress:       mobyprogress.NewProgressOutput(mobyprogress.NewOut(os.Stderr)),
		allRules:       make(map[label.Label]*rule.Rule),
	}
}

// scalaLang implements language.Language.
type scalaLang struct {
	// cacheFile is the main cache file, if enabled
	cacheFile string
	// cache is the loaded cache, if configured
	cache *scpb.Cache
	// ruleRegistry is the rule registry implementation.  This holds the rules
	// configured via gazelle directives by the user.
	ruleRegistry RuleRegistry
	// sourceResolver is the source resolver implementation.
	sourceResolver *crossresolve.ScalaSourceCrossResolver
	// scalaCompiler is the compiler implementation.  This is passed to the
	// importRegistry for use during import disambiguation.
	scalaCompiler *scalaCompiler
	// packages is map from the config.Rel to *scalaPackage for the
	// workspace-relative packate name.
	packages map[string]*scalaPackage
	// isResolvePhase is a flag that is tracks if at least one Resolve() call
	// has occurred.  It can be used to determine when the rule indexing phase
	// has completed and deps resolution phase has started (it calls
	// onResolvePhase).
	isResolvePhase bool
	// resolvers is a list of cross resolver implementations named by the -scala_resolvers flag
	resolvers map[string]crossresolve.ConfigurableCrossResolver
	// lastPackage tracks if this is the last generated package
	lastPackage *scalaPackage
	// resolverNames is a comma-separated list of resolver to enable
	resolverNames string
	// remainingRules is a counter that tracks when all rules have been resolved.
	remainingRules int
	// totalRules is used for progress
	totalRules int
	// progress is the progress interface
	progress mobyprogress.Output
	// scalaExistingRules is the value of the scala_existing_rule repeatable flag
	scalaExistingRules stringSliceFlags
	// ruleIndex is a map of all known generated rules
	allRules map[label.Label]*rule.Rule
}

// Name implements part of the language.Language interface
func (sl *scalaLang) Name() string { return scalaLangName }

// KnownDirectives implements part of the language.Language interface
func (*scalaLang) KnownDirectives() []string {
	return []string{
		ruleDirective,
		overrideDirective,
		implicitImportDirective,
		scalaExplainDeps,
		scalaExplainSrcs,
		mapKindImportNameDirective,
	}
}
