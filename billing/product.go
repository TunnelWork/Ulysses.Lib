package billing

import (
	"errors"
	"time"

	"github.com/TunnelWork/Ulysses.Lib/server"
)

var (
	ErrInvalidSerialNumber = errors.New("billing: invalid serial number")
	ErrInvalidOwnerID      = errors.New("billing: invalid owner ID, need OwnerUserID or OwnerAffiliationID")
	ErrInvalidProductID    = errors.New("billing: invalid product ID")
	ErrInvalidWalletID     = errors.New("billing: invalid wallet ID")
)

type Product struct {
	serialNumber uint64 // primary key, the auto-incremented id for product "instances"

	// Ownership. At least one must be set.
	// if both are set, the product is treated
	// as a private owned product.
	OwnerUserID        uint64
	OwnerAffiliationID uint64

	// billing/payment
	ProductID       uint64    // identifier for product type, description, pricing, etc. productID is actually the Serial Number of Product Listing
	dateCreation    time.Time // Accurate to DAY only
	dateLastBill    time.Time // Accurate to DAY only, used for recurring billing
	dateTermination time.Time // Accurate to DAY only
	terminated      bool
	WalletID        uint64
	BillingOption   BillingOption
}

func (p *Product) SerialNumber() uint64 {
	return p.serialNumber
}

func (p *Product) DateCreation() time.Time {
	return p.dateCreation
}

func (p *Product) DateLastBill() time.Time {
	return p.dateLastBill
}

func (p *Product) DateTermination() time.Time {
	return p.dateTermination
}

func (p *Product) Terminated() bool {
	return p.terminated
}

func (p *Product) Add() (uint64, error) {
	return AddProduct(p)
}

func (p *Product) Save() error {
	return UpdateProduct(p)
}

func (p *Product) Terminate() error {
	return TerminateProductBySerialNumber(p.serialNumber)
}

func (p *Product) ToTerminateOn(terminationDate time.Time) error {
	return ToTerminateProductOn(p, terminationDate)
}

// CollectPayment() expects
func (p *Product) CollectPayment(v ...interface{}) error {
	if p.BillingOption.BillingCycle == USAGE_BASED {
		// For USAGE_BASED products, expect 1 argument: total
		var ok bool
		var total float64
		if len(v) == 1 {
			total, ok = v[0].(float64)
		}
		if ok {
			return CollectUsageBasedPayment(p, total)
		}
	} else {
		// For RECURRING products, expect no arguments
		if len(v) == 0 {
			return CollectRecurringPayment(p)
		}
	}
	return errors.New("billing: CollectPayment() received invalid arguments")
}

func (p *Product) ForceCollectPayment(v ...interface{}) error {
	if p.BillingOption.BillingCycle == USAGE_BASED {
		// For USAGE_BASED products, expect 1 argument: total
		var ok bool
		var total float64
		if len(v) == 1 {
			total, ok = v[0].(float64)
		}
		if ok {
			return ForceCollectUsageBasedPayment(p, total)
		}
	} else {
		// For RECURRING products, expect no arguments
		if len(v) == 0 {
			return ForceCollectRecurringPayment(p)
		}
	}
	return errors.New("billing: ForceCollectPayment() received invalid arguments")
}

// Add entry to database
func AddProduct(product *Product) (uint64, error) {
	// Verify all fields are valid
	// Belong to either an Affiliation (shared) or a User (owned)
	if product.OwnerUserID == 0 || product.OwnerAffiliationID == 0 {
		return 0, ErrInvalidOwnerID
	}
	if product.ProductID == 0 {
		return 0, ErrInvalidProductID
	}
	if product.WalletID == 0 {
		return 0, ErrInvalidWalletID
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

func ToTerminateProductOn(product *Product, terminationDate time.Time) error {
	// if Usage-Based, must be today
	if BeginningOfDay(terminationDate).Before(BeginningOfDay(time.Now())) {
		return errors.New("billing: you cannot rewrite the history")
	}

	if product.BillingOption.BillingCycle == USAGE_BASED {
		if !BeginningOfDay(terminationDate).Equal(BeginningOfDay(time.Now())) {
			return errors.New("billing: terminationDate for ToTerminateProductOn() on usage-based products must be today")
		}
		// suspend the server account immediately to prevent any further usage
		productListing, err := GetProductListingByID(product.ProductID)
		if err != nil {
			return err
		}

		// Get the server instance using server package
		serverInstance, err := server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
		if err != nil {
			return err
		}

		serverInstance.SuspendAccount(product.SerialNumber())
	}

	return toTerminateProductOn(product.SerialNumber(), terminationDate.Format("2006-01-02"))
}

func CollectRecurringPayment(product *Product) error {
	// First: Check if product is due (and active)
	var dueDate time.Time
	if product.BillingOption.BillingCycle == MONTHLY {
		dueDate = product.dateLastBill.AddDate(0, 1, 0)
	} else if product.BillingOption.BillingCycle == QUARTERLY {
		dueDate = product.dateLastBill.AddDate(0, 3, 0)
	} else if product.BillingOption.BillingCycle == ANNUALLY {
		dueDate = product.dateLastBill.AddDate(1, 0, 0)
	} else {
		return errors.New("billing: unrecognized billing cycle, maybe not for recurring")
	}
	// Check if due today or before today
	paymentIsDue := time.Now().After(dueDate) || time.Now().Equal(dueDate)

	if paymentIsDue {
		wallet, err := GetWalletByID(product.WalletID)
		if err != nil {
			return err
		}

		// Spend from bound wallet
		err = wallet.TrySpend(product.BillingOption.Price)
		if err != nil {
			return err
		}

		// Update last bill date
		err = updateDateLastBillProductBySerialNumber(product.serialNumber)
		return err
	} else {
		return nil // Not due yet, so no error
	}
}

// THIS FUNCTION IS NOT EXPECTED TO BE USED IN ANYWHERE FOR Pre-V2 builds.
// IT DOES NOT CHECK IF IT IS DUE OR NOT.
func ForceCollectRecurringPayment(product *Product) error {
	wallet, err := GetWalletByID(product.WalletID)
	if err != nil {
		return err
	}

	err = wallet.Spend(product.BillingOption.Price)

	if err != nil {
		return err
	}

	// Update last bill date
	err = updateDateLastBillProductBySerialNumber(product.serialNumber)
	return err
}

func CollectUsageBasedPayment(product *Product, total float64) error {
	// Cap the amount
	amount, err := usageBasedAmount(product, total)
	if err != nil {
		return err
	}
	if amount > 0 { // only charge when needed
		wallet, err := GetWalletByID(product.WalletID)
		if err != nil {
			return err
		}

		err = wallet.TrySpend(amount)
		if err != nil {
			return err
		}

		// Update spending
		product.BillingOption.CurrentMonthSpending += amount
		err = product.Save()
		return err
	}

	return nil
}

// Be sure to terminate this product after this.
func ForceCollectUsageBasedPayment(product *Product, total float64) error {
	// Cap the amount
	amount, err := usageBasedAmount(product, total)
	if err != nil {
		return err
	}
	if amount > 0 { // only charge when needed
		wallet, err := GetWalletByID(product.WalletID)
		if err != nil {
			return err
		}

		err = wallet.Spend(amount)
		if err != nil {
			return err
		}

		// Update spending
		product.BillingOption.CurrentMonthSpending += amount
		err = product.Save()
		return err
	}

	return nil
}

func usageBasedAmount(product *Product, total float64) (float64, error) {
	// First calculate how much to charge based on total(real-time) and current month spending
	if total < product.BillingOption.CurrentMonthSpending {
		return 0, errors.New("billing: usage based billing total less than current month spending")
	}
	var amount float64                                    // The amount to be charged
	if total > product.BillingOption.MonthlySpendingCap { // Overspent. Need to cap the amount
		amount = product.BillingOption.MonthlySpendingCap - product.BillingOption.CurrentMonthSpending
	} else {
		amount = total - product.BillingOption.CurrentMonthSpending
	}

	return amount, nil
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

// For auto-billing.
func ListProductsToTerminate() ([]*Product, error) {
	return listProductsToTerminate()
}
