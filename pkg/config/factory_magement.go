package config

import (
	ctlharvester "github.com/harvester/harvester/pkg/generated/controllers/harvesterhci.io"
	ctlcore "github.com/rancher/wrangler/v3/pkg/generated/controllers/core"
	"k8s.io/client-go/rest"
	"kubevirt.io/client-go/kubecli"

	ctlnetwork "github.com/harvester/harvester-network-controller/pkg/generated/controllers/network.harvesterhci.io"

	ctldevices "github.com/harvester/pcidevices/pkg/generated/controllers/devices.harvesterhci.io"
	ctlkubevirt "github.com/harvester/pcidevices/pkg/generated/controllers/kubevirt.io"
)

type FactoryManager struct {
	DeviceFactory    *ctldevices.Factory
	HarvesterFactory *ctlharvester.Factory
	CoreFactory      *ctlcore.Factory
	NetworkFactory   *ctlnetwork.Factory
	KubevirtFactory  *ctlkubevirt.Factory

	KubevirtClient kubecli.KubevirtClient
	Cfg            *rest.Config
}

func NewFactoryManager(
	deviceFactory *ctldevices.Factory,
	harvesterFactory *ctlharvester.Factory,
	coreFactory *ctlcore.Factory,
	networkFactory *ctlnetwork.Factory,
	kubevirtFactory *ctlkubevirt.Factory,
	kubevirtClient kubecli.KubevirtClient,
	cfg *rest.Config,
) *FactoryManager {
	return &FactoryManager{
		DeviceFactory:    deviceFactory,
		HarvesterFactory: harvesterFactory,
		CoreFactory:      coreFactory,
		NetworkFactory:   networkFactory,
		KubevirtFactory:  kubevirtFactory,
		KubevirtClient:   kubevirtClient,
		Cfg:              cfg,
	}
}
