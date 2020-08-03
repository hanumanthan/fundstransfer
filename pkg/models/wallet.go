package models

type Wallet struct {
	ID           uint  `gorm:"primary_key" json:"id"`
	Balance      int   `json:"balance"`
	UserId       uint  `json:"user_id"`
	MobileNumber int32 `json:"mobile_number"`
}

func (w *Wallet) GetWalletForMobileNumber(mobileNumber int32) error {
	if err := DB.Where("mobile_number = ?", mobileNumber).First(&w).Error; err != nil {
		return err
	}
	return nil
}

func (w *Wallet) GetWalletForUser(userId int) error {
	if err := DB.Where("user_id = ?", userId).First(&w).Error; err != nil {
		return err
	}
	return nil
}

func (w *Wallet) Save() error {
	if err := DB.Save(w).Error; err != nil {
		return err
	}
	return nil
}

func GetAllWallets() ([]Wallet, error) {
	var wallets []Wallet
	if err := DB.Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}
