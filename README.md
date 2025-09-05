# Go Role Strategy SDK

一个用于管理 Jenkins Role Strategy 插件的 Go SDK，提供完整的角色和权限管理功能。

## 功能特性

- 🔐 **权限模板管理** - 创建、删除和查询权限模板
- 👥 **角色管理** - 支持全局角色、项目角色和节点角色的完整生命周期管理
- 🎯 **用户/组分配** - 灵活的用户和组角色分配功能
- 📊 **查询功能** - 获取角色信息、分配情况和匹配的作业/代理
- 🌐 **多角色类型** - 支持 GlobalRole、ProjectRole 和 SlaveRole

## 安装

```bash
go get github.com/joker-bai/go-role-strategy-sdk
```

## 快速开始

### 初始化客户端

```go
package main

import (
    "fmt"
    "log"
    
    rolestrategy "github.com/joker-bai/go-role-strategy-sdk"
)

func main() {
    // 创建客户端
    client := rolestrategy.NewClient(
        "https://your-jenkins-url",
        "your-username",
        "your-api-token",
    )
}
```

### 权限模板管理

```go
// 添加权限模板
err := client.AddTemplate(
    "developer", 
    "hudson.model.Item.Read,hudson.model.Item.Build", 
    true, // overwrite
)
if err != nil {
    log.Fatal("AddTemplate:", err)
}

// 获取权限模板
template, err := client.GetTemplate("developer")
if err != nil {
    log.Fatal("GetTemplate:", err)
}
fmt.Printf("Template: %+v\n", template)

// 删除权限模板
err = client.RemoveTemplates([]string{"developer"}, false)
if err != nil {
    log.Fatal("RemoveTemplates:", err)
}
```

### 角色管理

```go
// 添加全局角色
err := client.AddRole(
    rolestrategy.GlobalRole, 
    "dev-lead", 
    "hudson.model.Hudson.Administer", 
    true, // overwrite
    "",   // pattern (项目角色使用)
    "",   // template
)
if err != nil {
    log.Fatal("AddRole:", err)
}

// 添加项目角色（带模式匹配）
err = client.AddRole(
    rolestrategy.ProjectRole, 
    "project-dev", 
    "hudson.model.Item.Read,hudson.model.Item.Build", 
    true, 
    "dev-.*", // 匹配以 dev- 开头的项目
    "",
)

// 获取角色信息
role, err := client.GetRole(rolestrategy.GlobalRole, "dev-lead")
if err != nil {
    log.Fatal("GetRole:", err)
}
fmt.Printf("Role: %+v\n", role)

// 删除角色
err = client.RemoveRoles(rolestrategy.GlobalRole, []string{"dev-lead"})
if err != nil {
    log.Fatal("RemoveRoles:", err)
}
```

### 用户和组分配

```go
// 分配用户角色
err := client.AssignUserRole(rolestrategy.GlobalRole, "dev-lead", "alice")
if err != nil {
    log.Fatal("AssignUserRole:", err)
}

// 分配组角色
err = client.AssignGroupRole(rolestrategy.ProjectRole, "project-dev", "developers")
if err != nil {
    log.Fatal("AssignGroupRole:", err)
}

// 取消用户角色分配
err = client.UnassignUserRole(rolestrategy.GlobalRole, "dev-lead", "alice")
if err != nil {
    log.Fatal("UnassignUserRole:", err)
}

// 取消组角色分配
err = client.UnassignGroupRole(rolestrategy.ProjectRole, "project-dev", "developers")
if err != nil {
    log.Fatal("UnassignGroupRole:", err)
}
```

### 查询功能

```go
// 获取所有角色
allRoles, err := client.GetAllRoles(rolestrategy.GlobalRole)
if err != nil {
    log.Fatal("GetAllRoles:", err)
}
fmt.Printf("All Roles: %+v\n", allRoles)

// 获取角色分配情况
assignments, err := client.GetRoleAssignments(rolestrategy.GlobalRole)
if err != nil {
    log.Fatal("GetRoleAssignments:", err)
}
fmt.Printf("Assignments: %+v\n", assignments)

// 获取全局角色名列表
globalRoles, err := client.GetGlobalRoleNames()
if err != nil {
    log.Fatal("GetGlobalRoleNames:", err)
}
fmt.Println("Global Roles:", globalRoles)

// 获取项目角色名列表
projectRoles, err := client.GetProjectRoleNames()
if err != nil {
    log.Fatal("GetProjectRoleNames:", err)
}
fmt.Println("Project Roles:", projectRoles)

// 获取匹配的作业
matchingJobs, err := client.GetMatchingJobs("dev-.*", 10)
if err != nil {
    log.Fatal("GetMatchingJobs:", err)
}
fmt.Printf("Matching Jobs: %+v\n", matchingJobs)

// 获取匹配的代理
matchingAgents, err := client.GetMatchingAgents("agent-.*", 10)
if err != nil {
    log.Fatal("GetMatchingAgents:", err)
}
fmt.Printf("Matching Agents: %+v\n", matchingAgents)
```

## 数据结构

### RoleType

```go
type RoleType string

const (
    GlobalRole  RoleType = "globalRoles"   // 全局角色
    ProjectRole RoleType = "projectRoles"  // 项目角色
    SlaveRole   RoleType = "slaveRoles"    // 节点角色
)
```

### PermissionTemplate

```go
type PermissionTemplate struct {
    Name          string          `json:"name"`
    PermissionIDs map[string]bool `json:"permissionIds"`
    IsUsed        bool            `json:"isUsed,omitempty"`
    SIDs          []SIDEntry      `json:"sids,omitempty"`
}
```

### RoleInfo

```go
type RoleInfo struct {
    PermissionIDs map[string]bool `json:"permissionIds"`
    SIDs          []SIDEntry      `json:"sids"`
    Pattern       string          `json:"pattern,omitempty"`
    Template      string          `json:"template,omitempty"`
}
```

### SIDEntry

```go
type SIDEntry struct {
    Type string `json:"type"` // "USER" 或 "GROUP"
    SID  string `json:"sid"`
}
```

### RoleAssignment

```go
type RoleAssignment struct {
    Name  string   `json:"name"`
    Type  string   `json:"type"` // "USER" 或 "GROUP"
    Roles []string `json:"roles"`
}
```

## API 参考

### 客户端

- `NewClient(baseURL, username, apiToken string) *Client` - 创建新的客户端实例

### 权限模板

- `AddTemplate(name, permissionIds string, overwrite bool) error` - 添加权限模板
- `GetTemplate(name string) (*PermissionTemplate, error)` - 获取权限模板
- `RemoveTemplates(names []string, force bool) error` - 删除权限模板

### 角色管理

- `AddRole(roleType RoleType, roleName, permissionIds string, overwrite bool, pattern, template string) error` - 添加角色
- `GetRole(roleType RoleType, roleName string) (*RoleInfo, error)` - 获取角色信息
- `GetAllRoles(roleType RoleType) (map[string][]SIDEntry, error)` - 获取所有角色
- `RemoveRoles(roleType RoleType, roleNames []string) error` - 删除角色
- `GetGlobalRoleNames() ([]string, error)` - 获取全局角色名列表
- `GetProjectRoleNames() ([]string, error)` - 获取项目角色名列表

### 用户/组分配

- `AssignUserRole(roleType RoleType, roleName, user string) error` - 分配用户角色
- `AssignGroupRole(roleType RoleType, roleName, group string) error` - 分配组角色
- `UnassignUserRole(roleType RoleType, roleName, user string) error` - 取消用户角色分配
- `UnassignGroupRole(roleType RoleType, roleName, group string) error` - 取消组角色分配
- `GetRoleAssignments(roleType RoleType) ([]RoleAssignment, error)` - 获取角色分配情况
- `DeleteUser(roleType RoleType, user string) error` - 删除用户的所有角色分配
- `DeleteGroup(roleType RoleType, group string) error` - 删除组的所有角色分配

### 查询功能

- `GetMatchingJobs(pattern string, maxJobs int) ([]MatchingItem, error)` - 获取匹配的作业
- `GetMatchingAgents(pattern string, maxAgents int) ([]MatchingItem, error)` - 获取匹配的代理

## 前置条件

1. Jenkins 服务器已安装并启用 [Role Strategy Plugin](https://plugins.jenkins.io/role-strategy/)
2. 具有管理员权限的 Jenkins 用户账号
3. 生成的 API Token（在 Jenkins 用户配置中生成）

## 错误处理

所有 API 方法都返回 `error` 类型，建议在生产环境中进行适当的错误处理：

```go
if err := client.AddRole(rolestrategy.GlobalRole, "test-role", "hudson.model.Item.Read", true, "", ""); err != nil {
    log.Printf("Failed to add role: %v", err)
    // 处理错误...
}
```

## 许可证

本项目采用 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 支持

如果您在使用过程中遇到问题，请：

1. 查看 [Jenkins Role Strategy Plugin 文档](https://plugins.jenkins.io/role-strategy/)
2. 提交 [Issue](https://github.com/joker-bai/go-role-strategy-sdk/issues)
3. 参考 `examples/main.go` 中的示例代码