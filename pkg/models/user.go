package models

type User struct {
	ID       uint   `gorm:"primary_key" json:"id"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

func (u *User) GetById(id int) error {
	if err := DB.Where("id = ?", id).First(&u).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetByName(name string) error {
	if err := DB.Where("name = ?", name).First(&u).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	if err := DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
