// models.go
package rolestrategy

// RoleType represents the type of role
type RoleType string

const (
	GlobalRole  RoleType = "globalRoles"
	ProjectRole RoleType = "projectRoles"
	SlaveRole   RoleType = "slaveRoles"
)

// PermissionTemplate represents a permission template
type PermissionTemplate struct {
	Name          string          `json:"name"`
	PermissionIDs map[string]bool `json:"permissionIds"`
	IsUsed        bool            `json:"isUsed,omitempty"`
	SIDs          []SIDEntry      `json:"sids,omitempty"`
}

// RoleInfo represents a role's details
type RoleInfo struct {
	PermissionIDs map[string]bool `json:"permissionIds"`
	SIDs          []SIDEntry      `json:"sids"`
	Pattern       string          `json:"pattern,omitempty"`
	Template      string          `json:"template,omitempty"`
}

// SIDEntry represents a user or group assignment
type SIDEntry struct {
	Type string `json:"type"` // "USER" or "GROUP"
	SID  string `json:"sid"`
}

// RoleAssignment represents a user/group and their roles
type RoleAssignment struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"` // "USER" or "GROUP"
	Roles []string `json:"roles"`
}

// MatchingItem represents a matched job or agent
type MatchingItem struct {
	Name string `json:"name"`
}
