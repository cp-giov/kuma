package reconcile

import (
	"context"

	envoy_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"

	config_core "github.com/kumahq/kuma/pkg/config/core"
	"github.com/kumahq/kuma/pkg/core"
	util_xds_v3 "github.com/kumahq/kuma/pkg/util/xds/v3"
)

var log = core.Log.WithName("kds").WithName("reconcile")

func NewReconciler(hasher envoy_cache.NodeHash, cache util_xds_v3.SnapshotCache, generator SnapshotGenerator, versioner util_xds_v3.SnapshotVersioner, mode config_core.CpMode) Reconciler {
	return &reconciler{
		hasher:    hasher,
		cache:     cache,
		generator: generator,
		versioner: versioner,
		mode:      mode,
	}
}

type reconciler struct {
	hasher    envoy_cache.NodeHash
	cache     util_xds_v3.SnapshotCache
	generator SnapshotGenerator
	versioner util_xds_v3.SnapshotVersioner
	mode      config_core.CpMode
}

func (r *reconciler) Reconcile(ctx context.Context, node *envoy_core.Node) error {
	new, err := r.generator.GenerateSnapshot(ctx, node)
	if err != nil {
		return err
	}
	if err := new.Consistent(); err != nil {
		return err
	}
	id := r.hasher.ID(node)
	old, _ := r.cache.GetSnapshot(id)
	new = r.versioner.Version(new, old)
	r.logChanges(new, old, node)
	return r.cache.SetSnapshot(id, new)
}

func (r *reconciler) logChanges(new util_xds_v3.Snapshot, old util_xds_v3.Snapshot, node *envoy_core.Node) {
	for _, typ := range new.GetSupportedTypes() {
		if old != nil && old.GetVersion(typ) != new.GetVersion(typ) {
			client := node.Id
			if r.mode == config_core.Zone {
				// we need to override client name because Zone is always a client to Global (on gRPC level)
				client = "global"
			}
			log.Info("detected changes in the resources. Sending changes to the client.", "resourceType", typ, "client", client)
		}
	}
}
