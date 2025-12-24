package resource

import (
	"sort"
	"strconv"
	"strings"

	"github.com/cosi-project/runtime/pkg/resource"
	"github.com/siderolabs/omni/client/pkg/omni/resources/omni"
)

// ToMap converts a resource to a map for display
func ToMap(r resource.Resource) map[string]interface{} {
	result := make(map[string]interface{})

	md := r.Metadata()
	result["id"] = md.ID()
	result["type"] = string(md.Type())
	result["version"] = md.Version().String()
	result["namespace"] = md.Namespace()

	// Add labels
	if labels := md.Labels(); labels != nil {
		result["labels"] = labels.Raw()
	}

	// Add spec based on resource type
	switch r := r.(type) {
	case *omni.Cluster:
		addClusterFields(result, r)
	case *omni.Machine:
		addMachineFields(result, r)
	case *omni.MachineSet:
		addMachineSetFields(result, r)
	case *omni.ClusterMachine:
		addClusterMachineFields(result, r)
	case *omni.MachineStatus:
		addMachineStatusFields(result, r)
	case *omni.MachineClass:
		addMachineClassFields(result, r)
	case *omni.KubernetesVersion:
		addKubernetesVersionFields(result, r)
	case *omni.EtcdBackup:
		addEtcdBackupFields(result, r)
	case *omni.Schematic:
		addSchematicFields(result, r)
	case *omni.OngoingTask:
		addOngoingTaskFields(result, r)
	default:
		// For unknown types, just return basic metadata
	}

	return result
}

func addClusterFields(result map[string]interface{}, c *omni.Cluster) {
	if spec := c.TypedSpec(); spec != nil && spec.Value != nil {
		result["kubernetes_version"] = spec.Value.KubernetesVersion
		result["talos_version"] = spec.Value.TalosVersion
	}
}

func addMachineFields(result map[string]interface{}, m *omni.Machine) {
	if spec := m.TypedSpec(); spec != nil && spec.Value != nil {
		if spec.Value.ManagementAddress != "" {
			result["management_address"] = spec.Value.ManagementAddress
		}
	}
}

func addMachineSetFields(result map[string]interface{}, ms *omni.MachineSet) {
	if spec := ms.TypedSpec(); spec != nil && spec.Value != nil {
		if spec.Value.MachineAllocation != nil {
			result["machine_class"] = spec.Value.MachineAllocation.Name
		} else if spec.Value.MachineClass != nil {
			result["machine_class"] = spec.Value.MachineClass.Name
		}
	}
}

func addClusterMachineFields(result map[string]interface{}, cm *omni.ClusterMachine) {
	if spec := cm.TypedSpec(); spec != nil && spec.Value != nil {
		// ClusterMachine ID is the machine ID
		result["machine_id"] = cm.Metadata().ID()
		result["kubernetes_version"] = spec.Value.KubernetesVersion
	}
}

func addMachineStatusFields(result map[string]interface{}, ms *omni.MachineStatus) {
	if spec := ms.TypedSpec(); spec != nil && spec.Value != nil {
		if spec.Value.Network != nil {
			if spec.Value.Network.Hostname != "" {
				result["hostname"] = spec.Value.Network.Hostname
			}
			if len(spec.Value.Network.Addresses) > 0 {
				result["addresses"] = spec.Value.Network.Addresses
			}
		}
		if spec.Value.PlatformMetadata != nil {
			if spec.Value.PlatformMetadata.Platform != "" {
				result["platform"] = spec.Value.PlatformMetadata.Platform
			}
		}
	}
}

func addPlatformMetadataFields(result map[string]interface{}, pm interface{}) {
	// Use type assertion to safely access fields
	if pmMap, ok := pm.(map[string]interface{}); ok {
		for k, v := range pmMap {
			result[k] = v
		}
		return
	}
	// Try to access common fields via reflection or direct access
	// For now, just store the platform if available
	if pmObj, ok := pm.(interface{ GetPlatform() string }); ok {
		if platform := pmObj.GetPlatform(); platform != "" {
			result["platform"] = platform
		}
	}
}

func addMachineClassFields(result map[string]interface{}, mc *omni.MachineClass) {
	if spec := mc.TypedSpec(); spec != nil && spec.Value != nil {
		if spec.Value.MatchLabels != nil {
			result["match_labels"] = spec.Value.MatchLabels
		}
	}
}

func addKubernetesVersionFields(result map[string]interface{}, kv *omni.KubernetesVersion) {
	if spec := kv.TypedSpec(); spec != nil && spec.Value != nil {
		result["version"] = spec.Value.Version
	}
}

func addEtcdBackupFields(result map[string]interface{}, eb *omni.EtcdBackup) {
	if spec := eb.TypedSpec(); spec != nil && spec.Value != nil {
		if spec.Value.Snapshot != "" {
			result["snapshot"] = spec.Value.Snapshot
		}
		if spec.Value.CreatedAt != nil {
			result["created_at"] = spec.Value.CreatedAt
		}
	}
}

func addSchematicFields(result map[string]interface{}, s *omni.Schematic) {
	if spec := s.TypedSpec(); spec != nil && spec.Value != nil {
		// Use the resource ID as schematic ID
		result["schematic_id"] = s.Metadata().ID()
	}
}

func addOngoingTaskFields(result map[string]interface{}, ot *omni.OngoingTask) {
	if spec := ot.TypedSpec(); spec != nil && spec.Value != nil {
		result["title"] = spec.Value.Title
		result["details"] = spec.Value.Details
	}
}

// CompareVersions compares two version strings (semver-like)
func CompareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		part1 := getPart(parts1, i)
		part2 := getPart(parts2, i)

		cmp := compareVersionPart(part1, part2)
		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

func getPart(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return "0"
}

func compareVersionPart(part1, part2 string) int {
	val1 := parseVersionPart(part1)
	val2 := parseVersionPart(part2)

	if val1 < val2 {
		return -1
	}
	if val1 > val2 {
		return 1
	}
	return 0
}

func parseVersionPart(part string) int {
	numStr := extractNumericPrefix(part)
	if numStr == "" {
		return 0
	}

	val, err := parseInteger(numStr)
	if err != nil {
		return 0
	}

	return val
}

func extractNumericPrefix(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			result.WriteRune(r)
		} else {
			break
		}
	}
	return result.String()
}

func parseInteger(s string) (int, error) {
	return strconv.Atoi(s)
}

// GetVersionString extracts version string from a resource map
func GetVersionString(resourceMap map[string]interface{}) string {
	if version, ok := resourceMap["kubernetes_version"].(string); ok {
		return version
	}
	if version, ok := resourceMap["version"].(string); ok {
		return version
	}
	return ""
}

// SortKubernetesVersions sorts a slice of version strings (largest first)
func SortKubernetesVersions(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		return CompareVersions(versions[i], versions[j]) > 0
	})
}
