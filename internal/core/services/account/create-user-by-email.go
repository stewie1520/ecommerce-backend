package account

import (
	"time"

	pb_forms "github.com/pocketbase/pocketbase/forms"
	pb_models "github.com/pocketbase/pocketbase/models"
	pb_types "github.com/pocketbase/pocketbase/tools/types"
	"github.com/stewie1520/ecommerce-backend/internal/core"
)

func (as *AccountService) CreateUserByEmail(rs core.RequestScope, email string, password string, name string, birthday string, verify bool) (string, error) {
	passwordHash, err := as.hash.HashPassword(password)
	if err != nil {
		return "", err
	}

	verificationAt := pb_types.DateTime{}
	verificationAt.Scan(time.Now())

	pbApp := rs.PbApp()
	user := &pb_models.User{
		BaseAccount: pb_models.BaseAccount{
			Email:        email,
			PasswordHash: passwordHash,
		},
		Verified:               true,
		LastVerificationSentAt: verificationAt,
	}

	createUserForm := pb_forms.NewUserUpsert(pbApp, user)
	createUserForm.Password = password
	createUserForm.PasswordConfirm = password

	err = createUserForm.Submit()
	if err != nil {
		return "", err
	}

	user.Profile.SetDataValue("name", name)
	user.Profile.SetDataValue("birthday", birthday)
	err = pbApp.Dao().SaveRecord(user.Profile)
	if err != nil {
		return "", err
	}

	return user.Id, nil
}
