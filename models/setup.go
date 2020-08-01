package models

type userWithWallet struct {
	name         string
	location     string
	balance      int
	mobileNumber int32
}

func CreateTables() {
	DB.AutoMigrate(&User{}, &Transaction{}, &Wallet{})
}

func CreateUserAndWallet() {
	userWithWallets := []userWithWallet{
		{
			name:         "Athos",
			location:     "Milky Way",
			balance:      100,
			mobileNumber: 9999,
		},
		{
			name:         "Pothos",
			location:     "Andromeda",
			balance:      100,
			mobileNumber: 8888,
		},
		{
			name:         "Aramis",
			location:     "Whirlpool",
			balance:      100,
			mobileNumber: 7777,
		},
	}
	for _, element := range userWithWallets {
		user := &User{
			Name:     element.name,
			Location: element.location,
		}
		DB.Create(user)
		wallet := &Wallet{
			Balance:      element.balance,
			UserId:       user.ID,
			MobileNumber: element.mobileNumber,
		}
		DB.Create(wallet)
	}
}
