package repository

import (
	"cms-octo-chat-api/model"

	"gorm.io/gorm"
)

// CreateUser inserts a new user into the database
func (m *BaseRepository) CreateUser(user *model.User) error {
	return m.DB.Create(user).Error
}

func (m *BaseRepository) CeateUserWithMatrix(user model.User, matrix model.UserMatrix) error {
	return m.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		matrix.UserID = user.ID
		if err := tx.Create(&matrix).Error; err != nil {
			return err
		}
		return nil
	})
}

func (m *BaseRepository) CreateUserMatrix(matrix *model.UserMatrix) error {
	return m.DB.FirstOrCreate(matrix).Error
}

// GetUserByID finds a user by ID
func (m *BaseRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	if err := m.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail finds a user by email
func (m *BaseRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := m.DB.Preload("UserMatrix").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user record
func (m *BaseRepository) UpdateUser(user *model.User) error {
	return m.DB.Save(user).Error
}

// DeleteUser removes a user by ID (soft delete)
func (m *BaseRepository) DeleteUser(id uint) error {
	return m.DB.Delete(&model.User{}, id).Error
}
