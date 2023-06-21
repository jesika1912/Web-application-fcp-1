package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory() ([]model.UserTaskCategory, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	// TODO: replace this
	var user model.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// User not found
			return model.User{}, nil
		}
		// Other error occurred
		return model.User{}, result.Error
	}
	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory() ([]model.UserTaskCategory, error) {
	// TODO: replace this
	var userTaskCategories []model.UserTaskCategory
	result := r.db.Table("users").
		Select("users.id, users.fullname, users.email, tasks.title AS task, tasks.deadline, tasks.priority, tasks.status, categories.name AS category").
		Joins("JOIN tasks ON users.id = tasks.user_id").
		Joins("JOIN categories ON tasks.category_id = categories.id").
		Scan(&userTaskCategories)
	if result.Error != nil {
		return nil, result.Error
	}
	return userTaskCategories, nil
}
