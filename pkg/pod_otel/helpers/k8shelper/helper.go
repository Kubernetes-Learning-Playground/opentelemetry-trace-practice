package k8shelper

import (
	"fmt"
	core "k8s.io/api/core/v1"
	flowcontrol "k8s.io/api/flowcontrol/v1beta3"
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"strings"
)

const (
	NodeUnreachablePodReason = "NodeLost"
	// NodeUnreachablePodMessage is the message on a pod when its state cannot be confirmed as kubelet is unresponsive
	// on the node it is (was) running.
	NodeUnreachablePodMessage = "Node %v which was running pod %v is unresponsive"

	APIGroupAll    = "*"
	ResourceAll    = "*"
	VerbAll        = "*"
	NonResourceAll = "*"

	GroupKind          = "Group"
	ServiceAccountKind = "ServiceAccount"
	UserKind           = "User"

	// AutoUpdateAnnotationKey is the name of an annotation which prevents reconciliation if set to "false"
	AutoUpdateAnnotationKey = "rbac.authorization.kubernetes.io/autoupdate"
)

func ContainsAccessMode(modes []core.PersistentVolumeAccessMode, mode core.PersistentVolumeAccessMode) bool {
	for _, m := range modes {
		if m == mode {
			return true
		}
	}
	return false
}

func removeDuplicateAccessModes(modes []core.PersistentVolumeAccessMode) []core.PersistentVolumeAccessMode {
	accessModes := []core.PersistentVolumeAccessMode{}
	for _, m := range modes {
		if !ContainsAccessMode(accessModes, m) {
			accessModes = append(accessModes, m)
		}
	}
	return accessModes
}

func GetAccessModesAsString(modes []core.PersistentVolumeAccessMode) string {

	modes = removeDuplicateAccessModes(modes)
	modesStr := []string{}
	if ContainsAccessMode(modes, core.ReadWriteOnce) {
		modesStr = append(modesStr, "RWO")
	}
	if ContainsAccessMode(modes, core.ReadOnlyMany) {
		modesStr = append(modesStr, "ROX")
	}
	if ContainsAccessMode(modes, core.ReadWriteMany) {
		modesStr = append(modesStr, "RWX")
	}
	if ContainsAccessMode(modes, core.ReadWriteOncePod) {
		modesStr = append(modesStr, "RWOP")
	}
	return strings.Join(modesStr, ",")
}

func GetPersistentVolumeClass(volume *core.PersistentVolume) string {
	// Use beta annotation first
	if class, found := volume.Annotations[core.BetaStorageClassAnnotation]; found {
		return class
	}

	return volume.Spec.StorageClassName
}

const IsDefaultStorageClassAnnotation = "storageclass.kubernetes.io/is-default-class"

// BetaIsDefaultStorageClassAnnotation is the beta version of IsDefaultStorageClassAnnotation.
// TODO: remove Beta when no longer used
const BetaIsDefaultStorageClassAnnotation = "storageclass.beta.kubernetes.io/is-default-class"

func IsDefaultAnnotation(obj metav1.ObjectMeta) bool {
	if obj.Annotations[IsDefaultStorageClassAnnotation] == "true" {
		return true
	}
	if obj.Annotations[BetaIsDefaultStorageClassAnnotation] == "true" {
		return true
	}

	return false
}

func IsDefaultAnnotationText(obj metav1.ObjectMeta) string {
	if obj.Annotations[IsDefaultStorageClassAnnotation] == "true" {
		return "Yes"
	}
	if obj.Annotations[BetaIsDefaultStorageClassAnnotation] == "true" {
		return "Yes"
	}

	return "No"
}
func SubjectsStrings(subjects []rbac.Subject) ([]string, []string, []string, []string) {
	users := []string{}
	groups := []string{}
	sas := []string{}
	others := []string{}

	for _, subject := range subjects {
		switch subject.Kind {
		case ServiceAccountKind:
			sas = append(sas, fmt.Sprintf("%s/%s", subject.Namespace, subject.Name))

		case UserKind:
			users = append(users, subject.Name)

		case GroupKind:
			groups = append(groups, subject.Name)

		default:
			others = append(others, fmt.Sprintf("%s/%s/%s", subject.Kind, subject.Namespace, subject.Name))
		}
	}

	return users, groups, sas, others
}

func GetPersistentVolumeClaimClass(claim *core.PersistentVolumeClaim) string {
	// Use beta annotation first
	if class, found := claim.Annotations[core.BetaStorageClassAnnotation]; found {
		return class
	}

	if claim.Spec.StorageClassName != nil {
		return *claim.Spec.StorageClassName
	}

	return ""
}

var _ sort.Interface = FlowSchemaSequence{}

// FlowSchemaSequence holds sorted set of pointers to FlowSchema objects.
// FlowSchemaSequence implements `sort.Interface`
type FlowSchemaSequence []*flowcontrol.FlowSchema

func (s FlowSchemaSequence) Len() int {
	return len(s)
}

func (s FlowSchemaSequence) Less(i, j int) bool {
	// the flow-schema w/ lower matching-precedence is prior
	if ip, jp := s[i].Spec.MatchingPrecedence, s[j].Spec.MatchingPrecedence; ip != jp {
		return ip < jp
	}
	// sort alphabetically
	return s[i].Name < s[j].Name
}

func (s FlowSchemaSequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
