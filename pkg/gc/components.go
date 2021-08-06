package gc

import (
	"time"

	"github.com/pkg/errors"

	config_core "github.com/kumahq/kuma/pkg/config/core"
	"github.com/kumahq/kuma/pkg/core/resources/apis/mesh"
	"github.com/kumahq/kuma/pkg/core/resources/apis/system"
	"github.com/kumahq/kuma/pkg/core/resources/model"
	"github.com/kumahq/kuma/pkg/core/runtime"
)

func Setup(rt runtime.Runtime) error {
	if err := setupCollector(rt); err != nil {
		return err
	}
	if err := setupFinalizer(rt); err != nil {
		return err
	}
	return nil
}

func setupCollector(rt runtime.Runtime) error {
	switch rt.Config().Environment {
	// Dataplane GC is run only on Universal because on Kubernetes Dataplanes are bounded by ownership to Pods.
	// Therefore on K8S offline dataplanes are cleaned up quickly enough to not run this.
	case config_core.UniversalEnvironment:
		return rt.Add(
			NewCollector(rt.ResourceManager(), 1*time.Minute, rt.Config().Runtime.Universal.DataplaneCleanupAge),
		)
	default:
		return nil
	}
}

func setupFinalizer(rt runtime.Runtime) error {
	var newTicker func() *time.Ticker
	var resourceTypes []model.ResourceType

	switch rt.Config().Mode {
	case config_core.Standalone:
		newTicker = func() *time.Ticker {
			return time.NewTicker(rt.Config().Metrics.Dataplane.IdleTimeout / 2)
		}
		resourceTypes = []model.ResourceType{
			mesh.DataplaneInsightType,
		}
	case config_core.Zone:
		newTicker = func() *time.Ticker {
			return time.NewTicker(rt.Config().Metrics.Dataplane.IdleTimeout / 2)
		}
		resourceTypes = []model.ResourceType{
			mesh.DataplaneInsightType,
			mesh.ZoneIngressInsightType,
		}
	case config_core.Global:
		newTicker = func() *time.Ticker {
			return time.NewTicker(rt.Config().Metrics.Zone.IdleTimeout / 2)
		}
		resourceTypes = []model.ResourceType{
			system.ZoneInsightType,
		}
	default:
		return errors.Errorf("unknown Kuma CP mode %v", rt.Config().Mode)
	}

	finalizer, err := NewSubscriptionFinalizer(rt.ResourceManager(), newTicker, resourceTypes...)
	if err != nil {
		return err
	}
	return rt.Add(finalizer)
}
