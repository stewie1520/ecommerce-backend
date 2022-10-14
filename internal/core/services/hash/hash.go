package hash

import "golang.org/x/crypto/bcrypt"

type IHashService interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword string, password string) error
}

type BcryptHashService struct {
}

func NewBcryptHashService() *BcryptHashService {
	return &BcryptHashService{}
}

func (bhs BcryptHashService) HashPassword(password string) (string, error) {
	val, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(val), nil
}

func (bhs BcryptHashService) CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
