package billing

import (
	"errors"
	"time"
)

var (
	ErrInvalidSerialNumber = errors.New("billing: invalid serial number")
	ErrInvalidOwnerID      = errors.New("billing: invalid owner ID, need OwnerUserID or OwnerAffiliationID")
	ErrInvalidProductID    = errors.New("billing: invalid product ID")
	ErrInvalidWalletID     = errors.New("billing: invalid wallet ID")
)

type Product struct {
	serialNumber uint64 // primary key

	// Ownership. At least one must be set.
	// if both are set, the product is treated
	// as a private owned product.
	OwnerUserID        uint64
	OwnerAffiliationID uint64

	// // Server-Account binding -- moved
	// ServerType string
	// InstanceID string

	// billing/payment
	ProductID       uint64    // identifier for product type, description, pricing, etc.
	dateCreation    time.Time // Accurate to DAY only
	dateTermination time.Time // Accurate to DAY only
	BillingCycle    uint8
	WalletID        uint64
}

func (p *Product) SerialNumber() uint64 {
	return p.serialNumber
}

func (p *Product) DateCreation() time.Time {
	return p.dateCreation
}

func (p *Product) DateTermination() time.Time {
	return p.dateTermination
}

func (p *Product) Add() error {
	return AddProduct(p)
}

func (p *Product) Update() error {
	return UpdateProduct(p)
}

func (p *Product) Terminate() error {
	return TerminateProductBySerialNumber(p.serialNumber)
}

// Add entry to database
func AddProduct(product *Product) error {
	// Verify all fields are valid
	// Belong to either an Affiliation (shared) or a User (owned)
	if product.OwnerUserID == 0 || product.OwnerAffiliationID == 0 {
		return ErrInvalidOwnerID
	}
	if product.ProductID == 0 {
		return ErrInvalidProductID
	}
	if product.WalletID == 0 {
		return ErrInvalidWalletID
	}

	return addProduct(product)
}

func GetProductBySerialNumber(serialNumber uint64) (*Product, error) {
	return getProductBySerialNumber(serialNumber)
}

func UpdateProduct(product *Product) error {
	// Verify all fields are valid
	if product.serialNumber == 0 {
		// Invalid
		return ErrInvalidSerialNumber
	}

	// Belong to either an Affiliation (shared) or a User (owned)
	if product.OwnerUserID == 0 || product.OwnerAffiliationID == 0 {
		return ErrInvalidOwnerID
	}
	if product.ProductID == 0 {
		return ErrInvalidProductID
	}
	if product.WalletID == 0 {
		return ErrInvalidWalletID
	}

	return updateProduct(product)
}

func TerminateProductBySerialNumber(serialNumber uint64) error {
	return terminateProductBySerialNumber(serialNumber)
}

// For user viewing.
// Lists all products owned by the user.
// Client should implement local-pagination to reduce need of repeated query
// API Server should hide dateTermination in reponse, if dateTermination is earlier than dateCreation
func ListUserProducts(ownerUserID uint64) ([]*Product, error) {
	return listUserProducts(ownerUserID)
}

// For affliation user viewing.
// Lists all products owned by the affiliation, including both shared and private owned products.
// Client should implement local-pagination to reduce need of repeated query
// Client should mark private products (ownerUserID != 0) as private
// API Server should hide dateTermination in response, if dateTermination is earlier than dateCreation
func ListAffiliationProducts(ownerAffiliationID uint64) ([]*Product, error) {
	return listAffiliationProducts(ownerAffiliationID)
}

// For admin viewing.
func ListProductsByProductID(productID uint64) ([]*Product, error) {
	return listProductsByProductID(productID)
}

// For automated-billing purposes.
// Designed for Pay-As-You-Go billing. But may be used for other purposes later?
func ListActiveProductsByBillingCycle(billingCycle uint8) ([]*Product, error) {
	return listActiveProductsByBillingCycle(billingCycle)
}

// For admin viewing. Client should implement local-pagination
func ListAllProducts() ([]*Product, error) {
	return listAllProducts()
}
