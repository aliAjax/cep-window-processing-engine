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
		return nil
	}
	return f.Loader.LoadDefaults()
}
