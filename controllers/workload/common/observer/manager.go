package observer

import (
	"sync"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"
	"k8s.io/apimachinery/pkg/types"
)

// Manager for a set of observers
type Manager struct {
	observers map[types.NamespacedName]*Observer
	listeners []OnObservation // invoked on each observation event
	lock      sync.RWMutex
	settings  Settings
}

// NewManager returns a new manager
func NewManager(settings Settings) *Manager {
	return &Manager{
		lock:      sync.RWMutex{},
		settings:  settings,
		observers: make(map[types.NamespacedName]*Observer),
	}
}

// Observe gets or create a cluster state observer for the given cluster
// In case something has changed in the given zkClient (eg. different caCert), the observer is recreated accordingly
func (m *Manager) Observe(cluster types.NamespacedName, zkClient zk.BaseClient) *Observer {
	m.lock.RLock()
	observer, exists := m.observers[cluster]
	m.lock.RUnlock()

	switch {
	case !exists:
		return m.createObserver(cluster, zkClient)
		//case exists && !observer.zkClient.Equal(&zkClient):
	case exists && !observer.zkClient.IsAlive(&zkClient):
		log.Info("Replacing observer HTTP client", "namespace", cluster.Namespace, "zk_name", cluster.Name)
		m.StopObserving(cluster)
		return m.createObserver(cluster, zkClient)
	default:
		return observer
	}
}

// createObserver creates a new observer according to the given arguments,
// and create/replace its entry in the observers map
func (m *Manager) createObserver(cluster types.NamespacedName, zkClient zk.BaseClient) *Observer {
	observer := NewObserver(cluster, zkClient, m.settings, m.notifyListeners)
	observer.Start()
	m.lock.Lock()
	m.observers[cluster] = observer
	m.lock.Unlock()
	return observer
}

func (m *Manager) ObservedStateResolver(cluster types.NamespacedName, zkClient zk.BaseClient) State {
	return m.Observe(cluster, zkClient).LastState()
}

// notifyListeners notifies all listeners that an observation occurred.
func (m *Manager) notifyListeners(cluster types.NamespacedName, previousState State, newState State) {
	wg := sync.WaitGroup{}
	m.lock.Lock()
	wg.Add(len(m.listeners))
	// run all listeners in parallel
	for _, l := range m.listeners {
		go func(f OnObservation) {
			defer wg.Done()
			f(cluster, previousState, newState)
		}(l)
	}
	// release the lock asap
	m.lock.Unlock()
	// wait for all listeners to be done
	wg.Wait()
}

// AddObservationListener adds the given listener to the list of listeners notified
// on every observation.
func (m *Manager) AddObservationListener(listener OnObservation) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.listeners = append(m.listeners, listener)
}

func (m *Manager) StopObserving(cluster types.NamespacedName) {
	m.lock.RLock()
	observer, exists := m.observers[cluster]
	m.lock.RUnlock()
	if !exists {
		return
	}
	observer.Stop()
	m.lock.Lock()
	delete(m.observers, cluster)
	m.lock.Unlock()
}
