package common

type routingConfig struct {
	Subdomain string `json:"subdomain"`
}

type imagePolicyConfig struct {
	InternalRegistryHostname string `json:"internalRegistryHostname"`
}

// APIServerConfig stores configuration information about the current cluster
type APIServerConfig struct {
	ImagePolicyConfig imagePolicyConfig `json:"imagePolicyConfig"`
	RoutingConfig     routingConfig     `json:"routingConfig"`
}

const BackupServerVersion string = "openshift.io/backup-server-version"
const RestoreServerVersion string = "openshift.io/restore-server-version"

const BackupRegistryHostname string = "openshift.io/backup-registry-hostname"
const RestoreRegistryHostname string = "openshift.io/restore-registry-hostname"

const MigrationRegistry string = "openshift.io/migration-registry"

const SwingPVAnnotation string = "openshift.io/swing-pv"

// copy, swing, TODO: others (snapshot, custom, etc.)
const MigrateTypeAnnotation string = "openshift.io/migrate-type"

//stage, final. Only valid for copy type.
const MigrateCopyPhaseAnnotation string = "openshift.io/migrate-copy-phase"

const MigrateQuiesceAnnotation string = "openshift.io/migrate-quiesce-pods"

const PodStageLabel string = "migration-stage-pod"

// Restic annotations
const ResticRestoreAnnotationPrefix string = "snapshot.velero.io"
const ResticBackupAnnotation string = "backup.velero.io/backup-volumes"
