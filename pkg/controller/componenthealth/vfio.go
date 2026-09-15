package componenthealth

import (
	"context"
	"fmt"
	"reflect"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	harvesterv1 "github.com/harvester/harvester/pkg/apis/harvesterhci.io/v1beta1"
	ctlharvester "github.com/harvester/harvester/pkg/generated/controllers/harvesterhci.io"
	ctlharvesterv1 "github.com/harvester/harvester/pkg/generated/controllers/harvesterhci.io/v1beta1"
	"github.com/harvester/pcidevices/pkg/util/modules"
	"github.com/sirupsen/logrus"
)

const (
	componentName = "pcidevices-controller"
	checkName     = "vfio-modules-loaded"
	checkInterval = 30 * time.Second
)

var requiredVFIOModules = []string{"vfio_pci", "vfio_iommu_type1"}

type Handler struct {
	componentHealths ctlharvesterv1.ComponentHealthClient
	nodeName         string
}

func Register(ctx context.Context, factory *ctlharvester.Factory, nodeName string) error {
	componentHealths := factory.Harvesterhci().V1beta1().ComponentHealth()
	handler := &Handler{
		componentHealths: componentHealths,
		nodeName:         nodeName,
	}

	if err := handler.reconcile(); err != nil {
		logrus.Errorf("failed to update VFIO component health: %v", err)
	}

	go handler.run(ctx)
	return nil
}

func (h *Handler) run(ctx context.Context) {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := h.reconcile(); err != nil {
				logrus.Errorf("failed to update VFIO component health: %v", err)
			}
		}
	}
}

func (h *Handler) reconcile() error {
	loaded, err := modules.Exists("/proc/modules", requiredVFIOModules)
	if err != nil {
		return err
	}

	checks := map[string]harvesterv1.CheckResult{}
	if !loaded {
		checks[checkName] = harvesterv1.CheckResult{
			Severity:      harvesterv1.SeverityError,
			Message:       fmt.Sprintf("required VFIO kernel modules are not loaded: %v", requiredVFIOModules),
			AffectedCount: 1,
			AffectedResources: &harvesterv1.AffectedResources{
				APIVersion: corev1.SchemeGroupVersion.String(),
				Kind:       "Node",
				Names:      map[string]harvesterv1.AffectedResourceDetail{h.nodeName: {}},
			},
		}
	}

	name := fmt.Sprintf("%s-%s", componentName, h.nodeName)
	existing, err := h.componentHealths.Get(name, metav1.GetOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get ComponentHealth %s: %w", name, err)
	}
	if apierrors.IsNotFound(err) {
		_, err = h.componentHealths.Create(&harvesterv1.ComponentHealth{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
				Labels: map[string]string{
					harvesterv1.LabelKeyComponent: componentName,
					harvesterv1.LabelKeyNode:      h.nodeName,
				},
			},
			Status: harvesterv1.ComponentHealthStatus{
				LastCheckedAt: metav1.Now(),
				Checks:        checks,
			},
		})
		return err
	}

	if reflect.DeepEqual(existing.Status.Checks, checks) {
		return nil
	}
	updated := existing.DeepCopy()
	updated.Status.LastCheckedAt = metav1.Now()
	updated.Status.Checks = checks
	if updated.Labels == nil {
		updated.Labels = map[string]string{}
	}
	updated.Labels[harvesterv1.LabelKeyComponent] = componentName
	updated.Labels[harvesterv1.LabelKeyNode] = h.nodeName
	_, err = h.componentHealths.Update(updated)
	return err
}
