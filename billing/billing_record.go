package billing

import "time"

type BillingRecord struct {
	SerialNumber        uint64
	WalletID            uint64
	UserID              uint64
	ProductID           uint64
	ProductSerialNumber uint64
	BillingCycle        uint8
	BilledAmount        float64
	BilledAt            time.Time
}

func AddBillingRecord(record BillingRecord) (uint64, error) {
	return addBillingRecord(record)
}

func ListBillingRecordsByWalletID(walletID uint64) ([]BillingRecord, error) {
	return listBillingRecordsByWalletID(walletID)
}

func ListAllBillingRecords() ([]BillingRecord, error) {
	return listAllBillingRecords()
}
