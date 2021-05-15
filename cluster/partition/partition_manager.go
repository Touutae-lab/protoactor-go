package partition

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	clustering "github.com/AsynkronIT/protoactor-go/cluster"
	"github.com/AsynkronIT/protoactor-go/eventstream"
	"github.com/AsynkronIT/protoactor-go/log"
)

const (
	ActorNameIdentity  = "partition"
	ActorNamePlacement = "partition-activator"
)

type PartitionManager struct {
	cluster     *clustering.Cluster
	topologySub *eventstream.Subscription
}

func newPartitionManager(c *clustering.Cluster) *PartitionManager {
	return &PartitionManager{
		cluster: c,
	}
}

func (pm *PartitionManager) Start() {
	plog.Info("Started partition manager")
	system := pm.cluster.ActorSystem

	identityProps := actor.PropsFromProducer(func() actor.Actor { return newIdentityActor(pm.cluster, pm) })
	system.Root.SpawnNamed(identityProps, ActorNameIdentity)
	plog.Info("Started partition identity actor")

	activatorProps := actor.PropsFromProducer(func() actor.Actor { return newPlacementActor(pm.cluster, pm) })
	system.Root.SpawnNamed(activatorProps, ActorNamePlacement)
	plog.Info("Started partition placement actor")

	pm.topologySub = system.EventStream.
		Subscribe(func(ev interface{}) {
			if topology, ok := ev.(*clustering.ClusterTopology); ok {
				pm.onClusterTopology(topology)
			}
		})
}

func (pm *PartitionManager) Stop() {
	system := pm.cluster.ActorSystem
	system.EventStream.Unsubscribe(pm.topologySub)
	plog.Info("Stopped PartitionManager")
}

func (pm *PartitionManager) PidOfIdentityActor(addr string) *actor.PID {
	return actor.NewPID(addr, ActorNameIdentity)
}

func (pm *PartitionManager) PidOfActivatorActor(addr string) *actor.PID {
	return actor.NewPID(addr, ActorNamePlacement)
}

func (pm *PartitionManager) onClusterTopology(tplg *clustering.ClusterTopology) {
	plog.Debug("onClusterTopology", log.Uint64("eventId", tplg.TopologyHash))
	//	system := pm.cluster.ActorSystem
	//TODO: update identity owner lookup

}

func (pm *PartitionManager) Get(identity *clustering.ClusterIdentity) *actor.PID {
	return nil
}
