package models

import "time"

type GroupRole string

const (
	Leader  GroupRole = "leader"  // 组长
	Manager GroupRole = "manager" // 管理员
	member  GroupRole = "member"  // 成员
)

type GroupMember struct {
	Member   User
	Role     GroupRole
	JoinedAt time.Time
}
