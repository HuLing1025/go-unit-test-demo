package repo

import "gorm.io/gorm"

type UserModel struct {
	Id   string `gorm:"column:id" json:"id"`
	Name string `gorm:"column:name" json:"name"`
	Age  int    `gorm:"column:age" json:"age"`
}

type UserRepo struct {
	db *gorm.DB
}

// UserRepoInterface
type UserRepoInterface interface {
	Create(model *UserModel) error
	SelectOne(id string) (UserModel, error)
}

// NewUserRepo
func NewUserRepo(db *gorm.DB) UserRepoInterface {
	return &UserRepo{db: db}
}

func (u *UserModel) TableName() string {
	return "user"
}

func (r *UserRepo) Create(model *UserModel) error {
	return r.db.Create(model).Error
}

func (r *UserRepo) SelectOne(id string) (user UserModel, err error) {
	err = r.db.Model(&UserModel{}).
		Where(`id = ?`, id).
		First(&user).Error
	return
}
