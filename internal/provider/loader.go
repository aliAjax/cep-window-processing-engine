package provider

type memoryProvider struct {
	values map[string]string
}

func (p *memoryProvider) Get(key string) (string, bool) {
	if p == nil {
		return "", false
	}
	value, ok := p.values[key]
	return value, ok
}

func (p *memoryProvider) Put(key, value string) error {
	p.values[key] = value
	return nil
}

type Loader struct{}

func (Loader) LoadDefaults() *memoryProvider {
	return &memoryProvider{}
}
