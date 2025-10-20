// api.go
package rolestrategy

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const strategyPrefix = "role-strategy/strategy"

// AddTemplate adds a new permission template
func (c *Client) AddTemplate(name, permissionIds string, overwrite bool) error {
	v := url.Values{}
	v.Set("name", name)
	v.Set("permissionIds", permissionIds)
	v.Set("overwrite", fmt.Sprintf("%t", overwrite))

	req, err := c.newRequest("POST", strategyPrefix+"/addTemplate", v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add template failed: %d %s: %s", resp.StatusCode, resp.Status, string(body))
	}
	return nil
}

// RemoveTemplates removes one or more templates
func (c *Client) RemoveTemplates(names []string, force bool) error {
	v := url.Values{}
	v.Set("names", strings.Join(names, ","))
	v.Set("force", fmt.Sprintf("%t", force))

	req, err := c.newRequest("POST", strategyPrefix+"/removeTemplates", v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove templates failed: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// AddRole adds a new role
func (c *Client) AddRole(roleType RoleType, roleName, permissionIds string, overwrite bool, pattern, template string) error {
	v := url.Values{}
	v.Set("type", string(roleType))
	v.Set("roleName", roleName)
	v.Set("permissionIds", permissionIds)
	v.Set("overwrite", fmt.Sprintf("%t", overwrite))
	if pattern != "" {
		v.Set("pattern", pattern)
	}
	if template != "" {
		v.Set("template", template)
	}

	req, err := c.newRequest("POST", strategyPrefix+"/addRole", v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add role failed: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// RemoveRoles removes one or more roles
func (c *Client) RemoveRoles(roleType RoleType, roleNames []string) error {
	v := url.Values{}
	v.Set("type", string(roleType))
	v.Set("roleNames", strings.Join(roleNames, ","))

	req, err := c.newRequest("POST", strategyPrefix+"/removeRoles", v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remove roles failed: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// AssignUserRole assigns a user to a role
func (c *Client) AssignUserRole(roleType RoleType, roleName, user string) error {
	return c.assignRole(roleType, roleName, user, "user")
}

// AssignGroupRole assigns a group to a role
func (c *Client) AssignGroupRole(roleType RoleType, roleName, group string) error {
	return c.assignRole(roleType, roleName, group, "group")
}

func (c *Client) assignRole(roleType RoleType, roleName, sid, param string) error {
	v := url.Values{}
	v.Set("type", string(roleType))
	v.Set("roleName", roleName)
	v.Set(param, sid)

	endpoint := "assignUserRole"
	if param == "group" {
		endpoint = "assignGroupRole"
	}

	req, err := c.newRequest("POST", strategyPrefix+"/"+endpoint, v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("assign %s failed: %d %s", param, resp.StatusCode, resp.Status)
	}
	return nil
}

// UnassignUserRole removes a user from a role
func (c *Client) UnassignUserRole(roleType RoleType, roleName, user string) error {
	return c.unassignRole(roleType, roleName, user, "user")
}

// UnassignGroupRole removes a group from a role
func (c *Client) UnassignGroupRole(roleType RoleType, roleName, group string) error {
	return c.unassignRole(roleType, roleName, group, "group")
}

func (c *Client) unassignRole(roleType RoleType, roleName, sid, param string) error {
	v := url.Values{}
	v.Set("type", string(roleType))
	v.Set("roleName", roleName)
	v.Set(param, sid)

	endpoint := "unassignUserRole"
	if param == "group" {
		endpoint = "unassignGroupRole"
	}

	req, err := c.newRequest("POST", strategyPrefix+"/"+endpoint, v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unassign %s failed: %d %s", param, resp.StatusCode, resp.Status)
	}
	return nil
}

// DeleteUser removes a user from all roles
func (c *Client) DeleteUser(roleType RoleType, user string) error {
	return c.deleteSid(roleType, user, "deleteUser")
}

// DeleteGroup removes a group from all roles
func (c *Client) DeleteGroup(roleType RoleType, group string) error {
	return c.deleteSid(roleType, group, "deleteGroup")
}

func (c *Client) deleteSid(roleType RoleType, sid, endpoint string) error {
	v := url.Values{}
	v.Set("type", string(roleType))
	v.Set("user", sid) // 注意：deleteUser 和 deleteGroup 都用 'user' 参数

	req, err := c.newRequest("POST", strategyPrefix+"/"+endpoint, v.Encode())
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s failed: %d %s", endpoint, resp.StatusCode, resp.Status)
	}
	return nil
}

// GetTemplate gets a permission template
func (c *Client) GetTemplate(name string) (*PermissionTemplate, error) {
	u := fmt.Sprintf("%s/getTemplate?name=%s", strategyPrefix, url.QueryEscape(name))
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get template failed: %d %s", resp.StatusCode, resp.Status)
	}

	var tmpl PermissionTemplate
	if err := json.NewDecoder(resp.Body).Decode(&tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// GetRole gets a role's details
func (c *Client) GetRole(roleType RoleType, roleName string) (*RoleInfo, error) {
	u := fmt.Sprintf("%s/getRole?type=%s&roleName=%s", strategyPrefix,
		url.QueryEscape(string(roleType)), url.QueryEscape(roleName))
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get role failed: %d %s", resp.StatusCode, resp.Status)
	}

	var role RoleInfo
	if err := json.NewDecoder(resp.Body).Decode(&role); err != nil {
		return nil, err
	}
	return &role, nil
}

// GetAllRoles gets all roles for a type
func (c *Client) GetAllRoles(roleType RoleType) (map[string][]interface{}, error) {
	u := fmt.Sprintf("%s/getAllRoles?type=%s", strategyPrefix, url.QueryEscape(string(roleType)))
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get all roles failed: %d %s", resp.StatusCode, resp.Status)
	}

	var roles map[string][]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleAssignments gets all SIDs and their assigned roles
func (c *Client) GetRoleAssignments(roleType RoleType) ([]RoleAssignment, error) {
	u := fmt.Sprintf("%s/getRoleAssignments?type=%s", strategyPrefix, url.QueryEscape(string(roleType)))
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get role assignments failed: %d %s", resp.StatusCode, resp.Status)
	}

	var assignments []RoleAssignment
	if err := json.NewDecoder(resp.Body).Decode(&assignments); err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetMatchingJobs gets jobs matching a pattern
func (c *Client) GetMatchingJobs(pattern string, maxJobs int) ([]MatchingItem, error) {
	u := fmt.Sprintf("%s/getMatchingJobs?pattern=%s&maxJobs=%d", strategyPrefix,
		url.QueryEscape(pattern), maxJobs)
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get matching jobs failed: %d %s", resp.StatusCode, resp.Status)
	}

	var items []MatchingItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetMatchingAgents gets agents matching a pattern
func (c *Client) GetMatchingAgents(pattern string, maxAgents int) ([]MatchingItem, error) {
	u := fmt.Sprintf("%s/getMatchingAgents?pattern=%s&maxAgents=%d", strategyPrefix,
		url.QueryEscape(pattern), maxAgents)
	req, err := c.newRequest("GET", u, "")
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get matching agents failed: %d %s", resp.StatusCode, resp.Status)
	}

	var items []MatchingItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

// GetGlobalRoleNames 获取所有全局角色的名称列表
func (c *Client) GetGlobalRoleNames() ([]string, error) {
	return c.getRoleNames(GlobalRole)
}

// GetProjectRoleNames 获取所有项目（局部）角色的名称列表
func (c *Client) GetProjectRoleNames() ([]string, error) {
	return c.getRoleNames(ProjectRole)
}

// getRoleNames 是一个私有方法，用于获取指定类型的角色名列表
func (c *Client) getRoleNames(roleType RoleType) ([]string, error) {
	// 调用底层的 GetAllRoles 方法
	rolesMap, err := c.GetAllRoles(roleType)
	if err != nil {
		return nil, err
	}

	// 从 map 的 key 中提取出所有的角色名
	var roleNames []string
	for roleName := range rolesMap {
		roleNames = append(roleNames, roleName)
	}

	return roleNames, nil
}

