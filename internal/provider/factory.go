package provider

type Provider interface {
	Get(string) (string, bool)
	Put(string, string) error
}

type Factory struct {
	Loader Loader
}

func (f Factory) Build(enabled bool) Provider {
	if !enabled {
		var provider *memoryProvider
		return provider
	}
	return f.Loader.LoadDefaults()
}
