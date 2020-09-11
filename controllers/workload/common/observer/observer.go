package observer

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/ghostbaby/zookeeper-operator/controllers/workload/common/zk"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("controllers").WithName("observer")

// Settings for the Observer configuration
type Settings struct {
	ObservationInterval time.Duration
	RequestTimeout      time.Duration
}

// Default values:
// - best-case scenario (healthy cluster): a request is performed every 10 seconds
// - worst-case scenario (unhealthy cluster): a request is performed every 70 (60+10) seconds
const (
	DefaultObservationInterval = 10 * time.Second
	DefaultRequestTimeout      = 1 * time.Minute
)

// DefaultSettings is an observer's Params with default values
var DefaultSettings = Settings{
	ObservationInterval: DefaultObservationInterval,
	RequestTimeout:      DefaultRequestTimeout,
}

type Observer struct {
	cluster  types.NamespacedName
	zkClient zk.BaseClient

	settings Settings

	creationTime time.Time

	stopChan chan struct{}
	stopOnce sync.Once

	onObservation OnObservation

	lastState State
	mutex     sync.RWMutex
}

// OnObservation is a function that gets executed when a new state is observed
type OnObservation func(cluster types.NamespacedName, previousState State, newState State)

// NewObserver creates and starts an Observer
func NewObserver(cluster types.NamespacedName, esClient zk.BaseClient, settings Settings, onObservation OnObservation) *Observer {
	observer := Observer{
		cluster:       cluster,
		zkClient:      esClient,
		creationTime:  time.Now(),
		settings:      settings,
		stopChan:      make(chan struct{}),
		stopOnce:      sync.Once{},
		onObservation: onObservation,
		mutex:         sync.RWMutex{},
	}

	log.Info("Creating observer for cluster", "namespace", cluster.Namespace, "zk_name", cluster.Name)
	return &observer
}

// Start the observer in a separate goroutine
func (o *Observer) Start() {
	go o.runUntilStopped()
}

func (o *Observer) Stop() {
	// trigger an async stop, only once
	o.stopOnce.Do(func() {
		go func() {
			close(o.stopChan)
			o.zkClient.Close()
		}()
	})
}

// run the observer main loop, until stopped
func (o *Observer) runUntilStopped() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go o.runPeriodically(ctx)
	<-o.stopChan
}

// runPeriodically triggers a state retrieval every tick,
// until the given context is cancelled
func (o *Observer) runPeriodically(ctx context.Context) {
	o.retrieveState(ctx)
	ticker := time.NewTicker(o.settings.ObservationInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			o.retrieveState(ctx)
		case <-ctx.Done():
			log.Info("Stopping observer for cluster", "namespace", o.cluster.Namespace, "zk_name", o.cluster.Name)
			return
		}
	}
}

// retrieveState retrieves the current ES state, executes onObservation,
// and stores the new state
func (o *Observer) retrieveState(ctx context.Context) {

	timeoutCtx, cancel := context.WithTimeout(ctx, o.settings.RequestTimeout)
	defer cancel()

	newState := RetrieveState(timeoutCtx, o.cluster, o.zkClient)

	if o.onObservation != nil {
		o.onObservation(o.cluster, o.LastState(), newState)
	}

	o.mutex.Lock()
	o.lastState = newState
	o.mutex.Unlock()
}

// RetrieveState returns the current Zookeeper cluster state
func RetrieveState(ctx context.Context, cluster types.NamespacedName, zkClient zk.BaseClient) State {
	// retrieve both cluster state and health in parallel
	clusterStateChan := make(chan *zk.ClusterStats)

	go func() {
		var leaderNode string
		var availableNodes int

		for _, endpoint := range zkClient.Endpoints {

			//ctx, cancel := context.WithCancel(context.Background())
			//timeoutCtx, cancel := context.WithTimeout(ctx, DefaultSettings.RequestTimeout)
			//defer cancel()
			if zkClient.Endpoint != endpoint {
				zkClient.Endpoint = endpoint
			}

			clusterState, err := zkClient.GetClusterStatus(ctx)
			if err != nil {
				// This is expected to happen from time to time
				log.Info("Unable to retrieve cluster state", "error", err, "namespace", cluster.Namespace, "zk_name", cluster.Name)
				clusterStateChan <- nil
				continue
			}

			up, err := zkClient.GetClusterUp(ctx)
			if err != nil {
				// This is expected to happen from time to time
				log.Info("Unable to retrieve cluster state", "error", err, "namespace", cluster.Namespace, "zk_name", cluster.Name)
				clusterStateChan <- nil
				continue
			}

			if clusterState.Mode == 1 {
				leaderNode = strings.Split(strings.Split(endpoint, ":")[1], "//")[1]
			}

			if up {
				availableNodes++
			}
			//clusterStateChan <- &clusterState

		}

		clusterStateChan <- &zk.ClusterStats{
			AvailableNodes: availableNodes,
			LeaderNode:     leaderNode,
		}
	}()

	// return the state when ready, may contain nil values
	return State{
		ClusterStats: <-clusterStateChan,
	}
}

// LastState returns the last observed state
func (o *Observer) LastState() State {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	return o.lastState
}
