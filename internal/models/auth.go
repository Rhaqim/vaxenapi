package models

type UserRole string

const (
	UserRoleOwner   UserRole = "owner"
	UserRoleManager UserRole = "manager"
	UserRoleFinance UserRole = "finance"
	UserRoleViewer  UserRole = "viewer"
)
