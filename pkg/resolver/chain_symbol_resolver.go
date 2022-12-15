package resolver

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
)

// ChainSymbolResolver implements SymbolResolver over a chain of resolvers.
type ChainSymbolResolver struct {
	chain []SymbolResolver
}

func NewChainSymbolResolver(chain ...SymbolResolver) *ChainSymbolResolver {
	return &ChainSymbolResolver{
		chain: chain,
	}
}

// ResolveSymbol implements the SymbolResolver interface
func (r *ChainSymbolResolver) ResolveSymbol(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, imp string) (*Symbol, error) {
	for _, next := range r.chain {
		known, err := next.ResolveSymbol(c, ix, from, lang, imp)
		if err == nil {
			return known, err
		}
		if err == ErrImportNotFound {
			continue
		}
		return nil, err
	}

	return nil, ErrImportNotFound
}