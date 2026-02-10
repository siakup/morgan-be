package integrations

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	libtypes "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
	shiftGroupsRepo "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/repository/postgresql"
)

func TestShiftGroupsRepository(t *testing.T) {
	if testPool == nil {
		t.Skip("Skipping integration test: testPool is nil")
	}
	ctx := context.Background()
	repo := shiftGroupsRepo.NewRepository(testPool)

	// Fetch a valid user ID for audit log
	var userID string
	err := testPool.QueryRow(ctx, "SELECT id FROM auth.users LIMIT 1").Scan(&userID)
	if err != nil {

		t.Logf("Warning: Could not fetch user for audit: %v", err)
		userID = "00000000-0000-0000-0000-000000000001" // Fallback dummy
	}

	t.Run("CRUD", func(t *testing.T) {
		// 1. Create
		newGroup := &domain.ShiftGroup{
			Id:        "sg-integration-test-01",
			Name:      "IT",
			Status:    true,
			CreatedAt: time.Now(),
			CreatedBy: &userID,
			UpdatedAt: time.Now(),
			UpdatedBy: &userID,
		}

		err := repo.Store(ctx, newGroup)
		assert.NoError(t, err)

		// 2. Read (FindByID)
		fetchedGroup, err := repo.FindByID(ctx, newGroup.Id)
		assert.NoError(t, err)
		assert.NotNil(t, fetchedGroup)
		assert.Equal(t, newGroup.Id, fetchedGroup.Id)
		assert.Equal(t, newGroup.Name, fetchedGroup.Name)
		assert.Equal(t, newGroup.Status, fetchedGroup.Status)

		// 3. Update
		newGroup.Name = "IT Updated"
		newGroup.Status = false
		newGroup.UpdatedAt = time.Now()

		err = repo.Update(ctx, newGroup)
		assert.NoError(t, err)

		fetchedGroupAfterUpdate, err := repo.FindByID(ctx, newGroup.Id)
		assert.NoError(t, err)
		assert.Equal(t, "IT Updated", fetchedGroupAfterUpdate.Name)
		assert.False(t, fetchedGroupAfterUpdate.Status)

		// 4. Delete
		err = repo.Delete(ctx, newGroup.Id, userID)
		assert.NoError(t, err)

		deletedGroup, err := repo.FindByID(ctx, newGroup.Id)
		assert.NoError(t, err)
		assert.Nil(t, deletedGroup)
	})

	t.Run("FindAll", func(t *testing.T) {
		// Insert some data first to ensure there's something to find
		for i := 0; i < 3; i++ {
			sg := &domain.ShiftGroup{
				Id:        "sg-findall-" + string(rune(i)),
				Name:      "Group " + string(rune(i)),
				Status:    true,
				CreatedAt: time.Now(),
				CreatedBy: &userID,
				UpdatedAt: time.Now(),
				UpdatedBy: &userID,
			}
			_ = repo.Store(ctx, sg)
		}

		filter := domain.ShiftGroupFilter{
			Pagination: libtypes.Pagination{Page: 1, Size: 10},
		}

		groups, total, err := repo.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, groups)
		assert.Greater(t, total, int64(0))
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		group, err := repo.FindByID(ctx, "non-existent-id")
		assert.NoError(t, err) // Based on code, returns nil, nil
		assert.Nil(t, group)
	})
}
