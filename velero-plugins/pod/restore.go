package pod

import (
	"encoding/json"
	"fmt"

	"github.com/fusor/ocp-velero-plugin/velero-plugins/common"
	"github.com/heptio/velero/pkg/plugin/velero"
	"github.com/sirupsen/logrus"
	corev1API "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// RestorePlugin is a restore item action plugin for Velero
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to pods
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"pods"},
	}, nil
}

// Execute action for the restore plugin for the pod resource
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[pod-restore] Entering Pod restore plugin")

	pod := corev1API.Pod{}
	itemMarshal, _ := json.Marshal(input.Item)
	json.Unmarshal(itemMarshal, &pod)
	p.Log.Infof("[pod-restore] pod: %s", pod.Name)

	// delete temporary annotations and labels used by mig-controller during backup
	common.DeleteTemporaryKeys(pod.Labels, pod.Annotations)

	if input.Restore.Annotations[common.MigrateCopyPhaseAnnotation] == "stage" {
		common.ConfigureContainerSleep(pod.Spec.Containers, "infinity")
		common.ConfigureContainerSleep(pod.Spec.InitContainers, "0")
		pod.Labels[common.PodStageLabel] = "true"
	} else if input.Restore.Annotations[common.MigrateCopyPhaseAnnotation] != "" {

		registry := pod.Annotations[common.RestoreRegistryHostname]
		backupRegistry := pod.Annotations[common.BackupRegistryHostname]
		if registry == "" {
			return nil, fmt.Errorf("failed to find restore registry annotation")
		}
		common.SwapContainerImageRefs(pod.Spec.Containers, backupRegistry, registry, p.Log)
		common.SwapContainerImageRefs(pod.Spec.InitContainers, backupRegistry, registry, p.Log)

		ownerRefs, err := common.GetOwnerReferences(input.ItemFromBackup)
		if err != nil {
			return nil, err
		}
		// Check if pod has owner Refs and does not have restic backup associated with it
		if len(ownerRefs) > 0 && pod.Annotations[common.ResticBackupAnnotation] == "" {
			p.Log.Infof("[pod-restore] skipping restore of pod %s, has owner references and no restic backup", pod.Name)
			return velero.NewRestoreItemActionExecuteOutput(input.Item).WithoutRestore(), nil
		}
	}

	var out map[string]interface{}
	objrec, _ := json.Marshal(pod)
	json.Unmarshal(objrec, &out)

	return velero.NewRestoreItemActionExecuteOutput(&unstructured.Unstructured{Object: out}), nil
}
