package models

type UserRole string

const (
	UserRoleOwner   UserRole = "owner"
	UserRoleManager UserRole = "manager"
	UserRoleFinance UserRole = "finance"
	UserRoleViewer  UserRole = "viewer"
	UserRoleAdmin   UserRole = "admin" // Platform admin (not business admin)
)

func (r UserRole) IsDirectorRole() bool {
	return r == UserRoleOwner || r == UserRoleManager
}

func (r UserRole) CanInitiateTransactions() bool {
	return r == UserRoleOwner || r == UserRoleManager || r == UserRoleFinance
}

func (r UserRole) CanApprove() bool {
	return r == UserRoleOwner || r == UserRoleManager
}
