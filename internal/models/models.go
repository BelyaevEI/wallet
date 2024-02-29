package models

// Application constants
const (
	ConfigName string = "app"
	ConfigType string = "env"
)

type (
	// Struct of wallet
	Wallet struct {
		ID     uint32  `json:"id"`
		Amount float64 `json:"amount"`
	}

	// Transfer struct
	Transfer struct {
		WalletIDFrom uint32  `json:"fromid"`
		WalletIDTo   uint32  `json:"toid"`
		Amount       float64 `json:"amount"`
	}
)
