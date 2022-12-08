package scalarule

import (
	"fmt"
	"sort"
)

// globalProviderRegistry is the default registry singleton.
var globalProviderRegistry = &providerRegistryMap{
	providers: make(map[string]Provider),
}

// GlobalProviderRegistry returns a reference to the global ProviderRegistry
// implementation.
func GlobalProviderRegistry() ProviderRegistry {
	return globalProviderRegistry
}

// providerRegistryMap implements ProviderRegistry using a map.
type providerRegistryMap struct {
	providers map[string]Provider
}

// ProviderNames implements part of the ProviderRegistry interface.
func (p *providerRegistryMap) ProviderNames() []string {
	names := make([]string, 0, len(p.providers))
	for name := range p.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// RegisterProvider implements part of the ProviderRegistry interface.
func (p *providerRegistryMap) RegisterProvider(name string, provider Provider) error {
	_, ok := p.providers[name]
	if ok {
		return fmt.Errorf("provider already registered: %q", name)
	}
	p.providers[name] = provider
	return nil
}

// LookupProvider implements part of the RuleRegistry interface.
func (p *providerRegistryMap) LookupProvider(name string) (Provider, bool) {
	provider, ok := p.providers[name]
	if !ok {
		return nil, false
	}
	return provider, true
}
