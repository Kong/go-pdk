package server

import (
	"encoding/json"
	"fmt"
	"github.com/Kong/go-pdk"
	"log"
	"time"
)

type instanceData struct {
	id            int
	startTime     time.Time
	config        interface{}
	handlers      map[string]func(*pdk.PDK)
	lastEventTime time.Time
}

// Configuration data for a new plugin instance.
type PluginConfig struct {
	Name   string // plugin name
	Config []byte // configuration data, as a JSON string
}

type (
	certificater interface{ Certificate(*pdk.PDK) }
	rewriter     interface{ Rewrite(*pdk.PDK) }
	accesser     interface{ Access(*pdk.PDK) }
	responser    interface{ Response(*pdk.PDK) }
	prereader    interface{ Preread(*pdk.PDK) }
	logger       interface{ Log(*pdk.PDK) }
)

func getHandlers(config interface{}) map[string]func(*pdk.PDK) {
	handlers := map[string]func(*pdk.PDK){}

	if h, ok := config.(certificater); ok { handlers["certificate"] = h.Certificate }
	if h, ok := config.(rewriter)    ; ok { handlers["rewrite"]     = h.Rewrite     }
	if h, ok := config.(accesser)    ; ok { handlers["access"]      = h.Access      }
	if h, ok := config.(responser)   ; ok { handlers["response"]    = h.Response    }
	if h, ok := config.(prereader)   ; ok { handlers["preread"]     = h.Preread     }
	if h, ok := config.(logger)      ; ok { handlers["log"]         = h.Log         }

	return handlers
}

func (rh *rpcHandler) addInstance(instance *instanceData) {
	rh.lock.Lock()
	defer rh.lock.Unlock()

	instance.id = rh.nextInstanceId
	rh.nextInstanceId++
	rh.instances[instance.id] = instance
}

// Current state of a plugin instance.  TODO: add some statistics
type InstanceStatus struct {
	Name      string      // plugin name
	Id        int         // instance id
	Config    interface{} // configuration data, decoded
	StartTime int64
}

// StartInstance starts a plugin instance, as required by configuration data.  More than
// one instance can be started for a single plugin.  If the configuration changes,
// a new instance should be started and the old one closed.
//
// RPC exported method
func (rh *rpcHandler) StartInstance(config PluginConfig, status *InstanceStatus) error {
	// TODO: check if config.Name is the one we care

	instanceConfig := rh.constructor()

	if err := json.Unmarshal(config.Config, instanceConfig); err != nil {
		return fmt.Errorf("decoding config: %w", err)
	}

	instance := instanceData{
		startTime: time.Now(),
		config:    instanceConfig,
		handlers:  getHandlers(instanceConfig),
	}

// 	log.Printf("instance: %v", instance)

	rh.addInstance(&instance)

	*status = InstanceStatus{
		Name:      config.Name,
		Id:        instance.id,
		Config:    instance.config,
		StartTime: instance.startTime.Unix(),
	}

// 	log.Printf("Started instance %#v:%v", config.Name, instance.id)

	return nil
}

// InstanceStatus returns a given resource's status (the same given when started)
//
// RPC exported method
func (rh *rpcHandler) InstanceStatus(id int, status *InstanceStatus) error {
	rh.lock.RLock()
	instance, ok := rh.instances[id]
	rh.lock.RUnlock()
	if !ok {
		return fmt.Errorf("no plugin instance %d", id)
	}

	*status = InstanceStatus{
		Name:      "---",
		Id:        instance.id,
		Config:    instance.config,
		StartTime: instance.startTime.Unix(),
	}

	return nil
}

// CloseInstance is used when an instance shouldn't be used anymore.
// Doesn't kill any running event but the instance is no longer accesible,
// so it's not possible to start a new event with it and will be garbage
// collected after the last reference event finishes.
// Returns the status just before closing.
//
// RPC exported method
func (rh *rpcHandler) CloseInstance(id int, status *InstanceStatus) error {
	rh.lock.RLock()
	instance, ok := rh.instances[id]
	rh.lock.RUnlock()
	if !ok {
		return fmt.Errorf("no plugin instance %d", id)
	}

	*status = InstanceStatus{
		Name:      "---",
		Id:        instance.id,
		Config:    instance.config,
		StartTime: instance.startTime.Unix(),
	}

	// kill?

	log.Printf("closed instance %d", instance.id)

	rh.lock.Lock()
	rh.lastCloseInstance = time.Now()
	delete(rh.instances, id)
	rh.lock.Unlock()
	rh.expireInstances()

	return nil
}

func (rh *rpcHandler) expireInstances() {
	const instanceTimeout = 60
	expirationCutoff := time.Now().Add(time.Second * -instanceTimeout)

	rh.lock.Lock()
	oldinstances := []int{}
	for id, inst := range rh.instances {
		if inst.startTime.Before(expirationCutoff) && inst.lastEventTime.Before(expirationCutoff) {
			oldinstances = append(oldinstances, id)
		}
	}

	for _, id := range oldinstances {
		delete(rh.instances, id)
	}
	rh.lock.Unlock()

	for _, id := range oldinstances {
		log.Printf("closed instance %d", id)
	}
}
