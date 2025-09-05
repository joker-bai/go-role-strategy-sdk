// example/main.go
package main

import (
	"fmt"
	"log"

	rolestrategy "github.com/joker-bai/go-role-strategy-sdk"
)

func main() {
	client := rolestrategy.NewClient(
		"https://dev-sdas.changan.com/jenkins",
		"admin",
		"your-api-token",
	)

	// 1. 添加模板
	err := client.AddTemplate("developer", "hudson.model.Item.Read,hudson.model.Item.Build", true)
	if err != nil {
		log.Fatal("AddTemplate:", err)
	}

	// 2. 添加角色
	err = client.AddRole(rolestrategy.GlobalRole, "dev-lead", "hudson.model.Hudson.Administer", true, "", "")
	if err != nil {
		log.Fatal("AddRole:", err)
	}

	// 3. 分配用户
	err = client.AssignUserRole(rolestrategy.GlobalRole, "dev-lead", "alice")
	if err != nil {
		log.Fatal("AssignUserRole:", err)
	}

	// 4. 查询角色
	role, err := client.GetRole(rolestrategy.GlobalRole, "dev-lead")
	if err != nil {
		log.Fatal("GetRole:", err)
	}
	fmt.Printf("Role: %+v\n", role)

	// 5. 查询分配
	assignments, err := client.GetRoleAssignments(rolestrategy.GlobalRole)
	if err != nil {
		log.Fatal("GetRoleAssignments:", err)
	}
	fmt.Printf("Assignments: %+v\n", assignments)

	// 获取所有全局角色名
	globalRoles, err := client.GetGlobalRoleNames()
	if err != nil {
		log.Fatal("GetGlobalRoleNames failed:", err)
	}
	fmt.Println("Global Roles:", globalRoles) // 输出: [admin developer dev-lead]

	// 获取所有项目角色名
	projectRoles, err := client.GetProjectRoleNames()
	if err != nil {
		log.Fatal("GetProjectRoleNames failed:", err)
	}
	fmt.Println("Project Roles:", projectRoles) // 输出: [staging-reader prod-admin]
}
