package models

type Credentials struct {
	ID     uint `gorm:"primary_key"`
	UserId uint
	ApiKey string
}

func (c *Credentials) GetApiKeyForUser(userId int) error {
	if err := DB.Where("user_id = ?", userId).First(&c).Error; err != nil {
		return err
	}
	return nil
}
