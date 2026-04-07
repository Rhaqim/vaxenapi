package services_test

import (
	"testing"

	"vaxen/api/internal/services"
	"vaxen/api/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrganizationService_Get(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)
	org := testutil.SeedOrganization(t, db)

	result, err := svc.Get(org.ID)
	require.NoError(t, err)
	assert.Equal(t, org.ID, result.ID)
	assert.Equal(t, "Test Corp", result.Name)
}

func TestOrganizationService_Get_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)

	_, err := svc.Get("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestOrganizationService_Update(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)
	org := testutil.SeedOrganization(t, db)

	newName := "Updated Corp"
	newLegal := "Updated Corporation Ltd"
	updated, err := svc.Update(org.ID, services.UpdateOrganizationInput{
		Name:      &newName,
		LegalName: &newLegal,
	})
	require.NoError(t, err)
	assert.Equal(t, "Updated Corp", updated.Name)
	assert.Equal(t, "Updated Corporation Ltd", updated.LegalName)
	assert.Equal(t, org.Country, updated.Country, "untouched fields should stay the same")
}

func TestOrganizationService_Update_NoFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)
	org := testutil.SeedOrganization(t, db)

	updated, err := svc.Update(org.ID, services.UpdateOrganizationInput{})
	require.NoError(t, err)
	assert.Equal(t, org.Name, updated.Name, "should return unchanged org")
}

func TestOrganizationService_Update_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)

	name := "test"
	_, err := svc.Update("nonexistent", services.UpdateOrganizationInput{Name: &name})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestOrganizationService_GetUsers(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)
	org := testutil.SeedOrganization(t, db)
	testutil.SeedUser(t, db, org.ID, "owner", true)
	testutil.SeedUser(t, db, org.ID, "manager", true)

	users, err := svc.GetUsers(org.ID)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestOrganizationService_GetUsers_Empty(t *testing.T) {
	db := testutil.SetupTestDB(t)
	svc := services.NewOrganizationService(db)
	org := testutil.SeedOrganization(t, db)

	users, err := svc.GetUsers(org.ID)
	require.NoError(t, err)
	assert.Len(t, users, 0)
}
