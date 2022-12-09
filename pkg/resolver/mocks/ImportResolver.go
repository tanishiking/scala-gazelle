// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	config "github.com/bazelbuild/bazel-gazelle/config"
	label "github.com/bazelbuild/bazel-gazelle/label"

	mock "github.com/stretchr/testify/mock"

	resolve "github.com/bazelbuild/bazel-gazelle/resolve"

	resolver "github.com/stackb/scala-gazelle/pkg/resolver"

	rule "github.com/bazelbuild/bazel-gazelle/rule"
)

// ImportResolver is an autogenerated mock type for the ImportResolver type
type ImportResolver struct {
	mock.Mock
}

// AddKnownImportProvider provides a mock function with given fields: provider
func (_m *ImportResolver) AddKnownImportProvider(provider resolver.KnownImportProvider) error {
	ret := _m.Called(provider)

	var r0 error
	if rf, ok := ret.Get(0).(func(resolver.KnownImportProvider) error); ok {
		r0 = rf(provider)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetKnownImport provides a mock function with given fields: imp
func (_m *ImportResolver) GetKnownImport(imp string) (*resolver.KnownImport, bool) {
	ret := _m.Called(imp)

	var r0 *resolver.KnownImport
	if rf, ok := ret.Get(0).(func(string) *resolver.KnownImport); ok {
		r0 = rf(imp)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resolver.KnownImport)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(imp)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// GetKnownImports provides a mock function with given fields: prefix
func (_m *ImportResolver) GetKnownImports(prefix string) []*resolver.KnownImport {
	ret := _m.Called(prefix)

	var r0 []*resolver.KnownImport
	if rf, ok := ret.Get(0).(func(string) []*resolver.KnownImport); ok {
		r0 = rf(prefix)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*resolver.KnownImport)
		}
	}

	return r0
}

// GetKnownRule provides a mock function with given fields: from
func (_m *ImportResolver) GetKnownRule(from label.Label) (*rule.Rule, bool) {
	ret := _m.Called(from)

	var r0 *rule.Rule
	if rf, ok := ret.Get(0).(func(label.Label) *rule.Rule); ok {
		r0 = rf(from)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rule.Rule)
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(label.Label) bool); ok {
		r1 = rf(from)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// KnownImportProviders provides a mock function with given fields:
func (_m *ImportResolver) KnownImportProviders() []resolver.KnownImportProvider {
	ret := _m.Called()

	var r0 []resolver.KnownImportProvider
	if rf, ok := ret.Get(0).(func() []resolver.KnownImportProvider); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]resolver.KnownImportProvider)
		}
	}

	return r0
}

// PutKnownImport provides a mock function with given fields: known
func (_m *ImportResolver) PutKnownImport(known *resolver.KnownImport) error {
	ret := _m.Called(known)

	var r0 error
	if rf, ok := ret.Get(0).(func(*resolver.KnownImport) error); ok {
		r0 = rf(known)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutKnownRule provides a mock function with given fields: from, r
func (_m *ImportResolver) PutKnownRule(from label.Label, r *rule.Rule) error {
	ret := _m.Called(from, r)

	var r0 error
	if rf, ok := ret.Get(0).(func(label.Label, *rule.Rule) error); ok {
		r0 = rf(from, r)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ResolveKnownImport provides a mock function with given fields: c, ix, from, lang, sym
func (_m *ImportResolver) ResolveKnownImport(c *config.Config, ix *resolve.RuleIndex, from label.Label, lang string, sym string) (*resolver.KnownImport, error) {
	ret := _m.Called(c, ix, from, lang, sym)

	var r0 *resolver.KnownImport
	if rf, ok := ret.Get(0).(func(*config.Config, *resolve.RuleIndex, label.Label, string, string) *resolver.KnownImport); ok {
		r0 = rf(c, ix, from, lang, sym)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resolver.KnownImport)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*config.Config, *resolve.RuleIndex, label.Label, string, string) error); ok {
		r1 = rf(c, ix, from, lang, sym)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewImportResolver interface {
	mock.TestingT
	Cleanup(func())
}

// NewImportResolver creates a new instance of ImportResolver. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewImportResolver(t mockConstructorTestingTNewImportResolver) *ImportResolver {
	mock := &ImportResolver{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
