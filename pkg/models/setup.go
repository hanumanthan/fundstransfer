package models

type userWithWallet struct {
	name         string
	location     string
	balance      int
	mobileNumber int32
	apiKey       string
}

func CreateTables() {
	DB.AutoMigrate(&User{}, &Transaction{}, &Wallet{}, &Credentials{})
}

func SetupData() {
	userWithWallets := []userWithWallet{
		{
			name:         "Athos",
			location:     "Milky Way",
			balance:      100,
			mobileNumber: 9999,
			apiKey:       "abc@123",
		},
		{
			name:         "Pothos",
			location:     "Andromeda",
			balance:      100,
			mobileNumber: 8888,
			apiKey:       "def@456",
		},
		{
			name:         "Aramis",
			location:     "Whirlpool",
			balance:      100,
			mobileNumber: 7777,
			apiKey:       "ghi@789",
		},
		{
			name:     "admin",
			location: "Earth",
			apiKey:   "admin@123",
		},
	}
	for _, element := range userWithWallets {
		user := &User{
			Name:     element.name,
			Location: element.location,
		}
		DB.Create(user)

		credentials := &Credentials{
			UserId: user.ID,
			ApiKey: element.apiKey,
		}
		DB.Create(credentials)

		if element.name != "admin" {
			wallet := &Wallet{
				Balance:      element.balance,
				UserId:       user.ID,
				MobileNumber: element.mobileNumber,
			}
			DB.Create(wallet)
		}
	}
}
