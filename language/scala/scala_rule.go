package scala

import (
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	sppb "github.com/stackb/scala-gazelle/build/stack/gazelle/scala/parse"
	"github.com/stackb/scala-gazelle/pkg/collections"
	"github.com/stackb/scala-gazelle/pkg/resolver"
)

type scalaRule struct {
	// Rule is an embedded struct (FIXME: why embed this?).
	*rule.Rule
	// the parent config
	scalaConfig *scalaConfig
	// from is the label for the rule.
	from label.Label
	// files that are included in the rule.
	files []*sppb.File
	// the import resolver to which we chain to when self-imports are not matched.
	next resolver.KnownImportResolver
	// the registry implementation to which we provide known imports.
	registry resolver.KnownImportRegistry
	// requiredTypes is a mapping from the required type to the symbol that
	// needs it. for example, if 'class Foo requiredTypes Bar', "Bar" is the map
	// key and "Foo" will be the value.
	requiredTypes map[string][]string
	// exports represent symbols that are importable by other rules.
	exports map[string]resolve.ImportSpec
}

func newScalaRule(
	scalaConfig *scalaConfig,
	registry resolver.KnownImportRegistry,
	next resolver.KnownImportResolver,
	r *rule.Rule,
	from label.Label,
	files []*sppb.File,
) *scalaRule {
	scalaRule := &scalaRule{
		Rule:          r,
		scalaConfig:   scalaConfig,
		from:          from,
		files:         files,
		next:          next,
		registry:      registry,
		requiredTypes: make(map[string][]string),
		exports:       make(map[string]resolve.ImportSpec),
	}
	scalaRule.addFiles(files...)
	return scalaRule
}

func (r *scalaRule) addFiles(files ...*sppb.File) {
	for _, file := range files {
		r.addFromFile(file)
	}
}

func (r *scalaRule) addFromFile(file *sppb.File) {
	for _, imp := range file.Classes {
		r.putKnownImport(imp, sppb.ImportType_CLASS)
		r.putExport(imp)
	}
	for _, imp := range file.Objects {
		r.putKnownImport(imp, sppb.ImportType_OBJECT)
		r.putExport(imp)
	}
	for _, imp := range file.Traits {
		r.putKnownImport(imp, sppb.ImportType_TRAIT)
		r.putExport(imp)
	}
	for _, imp := range file.Types {
		r.putKnownImport(imp, sppb.ImportType_TYPE)
		r.putExport(imp)
	}
	for _, imp := range file.Vals {
		r.putKnownImport(imp, sppb.ImportType_VALUE)
		r.putExport(imp)
	}

	for token, extends := range file.Extends {
		r.putExtends(token, extends)
	}
	for _, imp := range file.Imports {
		r.putFileImport(imp)
	}
}

func (r *scalaRule) putFileImport(imp string) {
	// r.imports.Put(imp)
}

func (r *scalaRule) putExport(imp string) {
	r.exports[imp] = resolve.ImportSpec{Imp: imp, Lang: scalaLangName}
}

func (r *scalaRule) putKnownImport(imp string, impType sppb.ImportType) {
	// since we don't need to resolve same-rule symbols to a different label,
	// record all imports as label.NoLabel!
	r.registry.PutKnownImport(resolver.NewKnownImport(impType, imp, "self-import", label.NoLabel))
}

func (r *scalaRule) putExtends(token string, types *sppb.ClassList) {
	parts := strings.SplitN(token, " ", 2)
	if len(parts) != 2 {
		log.Fatalf("invalid extends token: %q: should have form '(class|interface|object) com.foo.Bar' ", token)
	}

	kind := parts[0]
	symbol := parts[1]

	r.putKindExtends(kind, symbol, types)
}

func (r *scalaRule) putKindExtends(kind, symbol string, types *sppb.ClassList) {
	switch kind {
	case "class":
		r.putClassExtends(symbol, types)
	case "interface":
		r.putInterfaceExtends(symbol, types)
	case "object":
		r.putObjectExtends(symbol, types)
	}
}

func (r *scalaRule) putClassExtends(imp string, types *sppb.ClassList) {
	r.putRequiredTypes(imp, types)
}

func (r *scalaRule) putInterfaceExtends(imp string, types *sppb.ClassList) {
	r.putRequiredTypes(imp, types)
}

func (r *scalaRule) putObjectExtends(imp string, types *sppb.ClassList) {
	r.putRequiredTypes(imp, types)
}

func (r *scalaRule) putRequiredTypes(imp string, types *sppb.ClassList) {
	for _, dst := range types.Classes {
		r.putRequiredType(imp, dst)
	}

}

// ResolveKnownImport implements the resolver.KnownImportResolver interface
func (r *scalaRule) ResolveKnownImport(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, imp string) (*resolver.KnownImport, error) {
	if known, ok := r.registry.GetKnownImport(imp); ok {
		return known, nil
	}
	return r.next.ResolveKnownImport(c, ix, from, lang, imp)
}

func (r *scalaRule) putRequiredType(src, dst string) {
	r.requiredTypes[dst] = append(r.requiredTypes[dst], src)
}

// Exports implements part of the scalarule.Rule interface.
func (r *scalaRule) Exports() []resolve.ImportSpec {
	exports := make([]resolve.ImportSpec, 0, len(r.exports))
	for _, v := range r.exports {
		exports = append(exports, v)
	}

	sort.Slice(exports, func(i, j int) bool {
		a := exports[i]
		b := exports[j]
		return a.Imp < b.Imp
	})

	return exports
}

// Imports implements part of the scalarule.Rule interface.
func (r *scalaRule) Imports() resolver.ImportMap {
	imports := resolver.NewImportMap()
	impLang := scalaLangName

	// direct
	for _, file := range r.files {
		for _, imp := range file.Imports {
			imports.Put(resolver.NewDirectImport(imp, file))
		}
	}

	// if this rule has a main_class
	if mainClass := r.AttrString("main_class"); mainClass != "" {
		imports.Put(resolver.NewMainClassImport(mainClass))
	}

	// add import required from extends clauses
	for imp, src := range r.requiredTypes {
		imports.Put(resolver.NewExtendsImport(imp, src[0])) // use first occurrence as source arg
	}

	// Initialize a list of symbols to find implicits for from all known
	// imports. Include all symbols that are defined in the rule too (a
	// gazelle:resolve_with directive should apply to them too).
	required := collections.StringStack(imports.Keys())
	for _, export := range r.Exports() {
		required = append(required, export.Imp)
	}

	// Gather implicit imports transitively.
	for !required.IsEmpty() {
		src, _ := required.Pop()
		for _, dst := range r.scalaConfig.getImplicitImports(impLang, src) {
			required.Push(dst)
			imports.Put(resolver.NewImplicitImport(dst, src))
		}
	}

	return imports
}

// Files implements part of the scalarule.Rule interface.
func (r *scalaRule) Files() []*sppb.File {
	return r.files
}
