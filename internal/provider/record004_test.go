package provider

import (
	"reflect"
	"testing"
)

func TestDefaultProviderInitializesState(t *testing.T) {
	provider := (Loader{}).LoadDefaults()
	if provider == nil || provider.values == nil {
		t.Fatal("default provider storage was not initialized")
	}
	if err := provider.Put("threshold", "42"); err != nil {
		t.Fatal(err)
	}
}

func TestFactoryReturnsTrueNil(t *testing.T) {
	provider := (Factory{}).Build(false)
	if provider != nil {
		t.Fatalf("disabled provider must be a nil interface, got %T", provider)
	}
}

func TestValidatorRejectsTypedNil(t *testing.T) {
	var concrete *memoryProvider
	var provider Provider = concrete
	if reflect.ValueOf(provider).IsNil() && (Validator{}).Validate(provider) == nil {
		t.Fatal("typed-nil provider passed validation")
	}
}

func TestServiceWritesDefaultProvider(t *testing.T) {
	service := Service{Provider: (Loader{}).LoadDefaults()}
	if err := service.PutVariable("window", "5m"); err != nil {
		t.Fatal(err)
	}
	if value, ok := service.Provider.Get("window"); !ok || value != "5m" {
		t.Fatalf("stored variable = %q, %v", value, ok)
	}
}
