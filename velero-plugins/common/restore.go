package common

import (
	"fmt"

	"github.com/heptio/velero/pkg/plugin/velero"
	"github.com/sirupsen/logrus"
)

// RestorePlugin is a restore item action plugin for Heptio Ark.
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to everything.
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{}, nil
}

// Execute sets a custom annotation on the item being restored.
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	p.Log.Info("[common-restore] Entering common restore plugin")

	metadata, annotations, err := getMetadataAndAnnotations(input.Item)
	if err != nil {
		return nil, err
	}
	name := metadata.GetName()
	p.Log.Infof("[common-restore] common restore plugin for %s", name)

	if input.Restore.Annotations[MigrateCopyPhaseAnnotation] != "" {

		version, err := GetServerVersion()
		if err != nil {
			return nil, err
		}

		annotations[RestoreServerVersion] = fmt.Sprintf("%v.%v", version.Major, version.Minor)
		registryHostname, err := GetRegistryInfo(version.Major, version.Minor)
		if err != nil {
			return nil, err
		}
		annotations[RestoreRegistryHostname] = registryHostname
		metadata.SetAnnotations(annotations)
	}

	return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
}
