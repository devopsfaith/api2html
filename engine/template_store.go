package engine

import "sync"

func NewTemplateStore() *TemplateStore {
	store := &TemplateStore{
		&templateStore{
			map[string]Renderer{},
			map[string]*sync.RWMutex{},
		},
		make(chan Subscription),
		&sync.Map{},
	}
	go store.subscribe()
	return store
}

type TemplateStore struct {
	*templateStore
	Subscribe chan Subscription
	observers *sync.Map
}

func (p *TemplateStore) subscribe() {
	for {
		subscription := <-p.Subscribe
		actual, loaded := p.observers.LoadOrStore(subscription.Name, []chan Renderer{subscription.In})
		if loaded {
			chans := actual.([]chan Renderer)
			p.observers.Store(subscription.Name, append(chans, subscription.In))
		}
	}
}

func (p *TemplateStore) Set(name string, tmpl Renderer) error {
	if err := p.templateStore.Set(name, tmpl); err != nil {
		return err
	}

	if actual, ok := p.observers.Load(name); ok {
		r := p.data[name]
		chans := actual.([]chan Renderer)
		for _, out := range chans {
			out <- r
		}
	}

	p.observers.Store(name, []chan Renderer{})
	return nil
}

type templateStore struct {
	data  map[string]Renderer
	mutex map[string]*sync.RWMutex
}

func (p *templateStore) Get(name string) (Renderer, bool) {
	m := p.getMutex(name)
	m.RLock()
	defer m.RUnlock()
	t, ok := p.data[name]
	return t, ok
}

func (p *templateStore) Set(name string, tmpl Renderer) error {
	m := p.getMutex(name)
	m.Lock()
	p.data[name] = tmpl
	m.Unlock()
	return nil
}

func (p *templateStore) getMutex(name string) *sync.RWMutex {
	m, ok := p.mutex[name]
	if !ok {
		m = &sync.RWMutex{}
		p.mutex[name] = m
	}
	return m
}
