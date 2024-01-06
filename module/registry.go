package module

import (
	"log/slog"
	"sync"
)

type Registry struct {
	components map[string]any
	mutex      sync.RWMutex
}

var registry = NewRegistry()

func DefaultRegistry() *Registry {
	return registry
}

func NewRegistry() *Registry {
	return &Registry{
		components: make(map[string]any),
		mutex:      sync.RWMutex{},
	}
}

func (r *Registry) Register(name string, component any) {
	_, exists := r.components[name]
	if exists {
		slog.Warn("component already registered", slog.String("component", name))
		return
	}

	r.components[name] = component
	slog.Info("component registered", slog.String("component", name))
}

func (r *Registry) Unregister(name string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.components, name)
	slog.Info("component unregistered", slog.String("component", name))
}

func (r *Registry) Get(name string) any {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.components[name]
}

func (r *Registry) List() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var l []string
	for k := range r.components {
		l = append(l, k)
	}
	return l
}

type DuplicateComponent string

func (r DuplicateComponent) Error() string {
	return "duplicate component: " + string(r)
}
