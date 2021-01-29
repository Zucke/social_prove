package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptPassword(t *testing.T) {
	u := User{
		Password: "123456",
	}

	assert.Nil(t, u.HashPassword)

	err := u.EncryptPassword()

	assert.NoError(t, err)

	assert.NotNil(t, u.HashPassword)
}

func TestComparePassword(t *testing.T) {
	const password = "123456"
	u := User{
		Password: password,
	}

	assert.Nil(t, u.HashPassword)

	err := u.EncryptPassword()

	assert.NoError(t, err)

	assert.NotNil(t, u.HashPassword)

	assert.True(t, u.ComparePassword(password))

	assert.False(t, u.ComparePassword("error"))
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		user  User
		match bool
	}{
		{
			user:  User{Email: "orlmicron@gmail.com"},
			match: true,
		},
		{
			user:  User{Email: "orlmicrongmail.com"},
			match: false,
		},
		{
			user:  User{Email: "test@example.es"},
			match: true,
		},
		{
			user:  User{Email: "orlando@outlook.com"},
			match: true,
		},
		{
			user:  User{Email: "orlmicron"},
			match: false,
		},
	}

	assert := assert.New(t)

	for _, test := range tests {
		t.Run(test.user.Email, func(t *testing.T) {
			assert.Equal(test.match, test.user.ValidateEmail())
		})
	}
}
