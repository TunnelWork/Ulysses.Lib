package billing

import "time"

type Product struct {
	ProductID uint64

	// Ownership
	OwnerID uint64

	// Server-Account binding
	ServerType string
	InstanceID string

	// billing/payment
	DateCreation    time.Time // Accurate to DAY only
	DateTermination time.Time // Accurate to DAY only
	BillingCycle    BillingCycle
	WalletID        uint64

	// //// Only for BillingCycle: PayAsYouGo
	// HourlyRate              float64
	// MonthlySpendingCap      float64 // This is the total amount secured.
	// CurrentMonthExpenditure float64 // This is the amount already spent.

	// //// Only for non-PayAsYouGo BillingCycle:
	// NextBillingDate time.Time
	// RecurringRate   float64
}

// Add entry to database
func AddProduct(product Product) error {
	stmtInsertProduct, err := sqlStatement(`INSERT INTO dbprefix_products 
    (OwnerID, ServerType, InstanceID, DateCreation, BillingCycle, WalletID) 
    VALUE
    (?, ?, ?, CURDATE(), ?, ?)`)
	if err != nil {
		return err
	}
	defer stmtInsertProduct.Close()

	_, err = stmtInsertProduct.Exec(
		product.OwnerID,
		product.ServerType,
		product.InstanceID,
		product.BillingCycle,
		product.WalletID,
	)
	return err
}

// For user viewing. Caller should do local pagination
func ListUserProducts(ownerID uint64) ([]Product, error) {
	stmtListUserProducts, err := sqlStatement(`SELECT ProductID, ServerType, InstanceID, DateCreation, DateTermination, BillingCycle, WalletID FROM dbprefix_products WHERE OwnerID = ?;`)
	if err != nil {
		return nil, err
	}
	defer stmtListUserProducts.Close()

	rows, err := stmtListUserProducts.Query(ownerID)
	if err != nil {
		return nil, err
	}
	var products []Product
	for rows.Next() {
		var product Product = Product{
			OwnerID: ownerID,
		}
		if err := rows.Scan(&product.ProductID, &product.ServerType, &product.InstanceID,
			&product.DateCreation, &product.DateTermination, &product.BillingCycle,
			&product.WalletID); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	return products, err
}

// For admin viewing. Caller should do local pagination
func ListAllProducts() ([]Product, error) {
	stmtListAllProducts, err := sqlStatement(`SELECT ProductID, OwnerID, ServerType, InstanceID, DateCreation, DateTermination, BillingCycle, WalletID FROM dbprefix_products;`)
	if err != nil {
		return nil, err
	}
	defer stmtListAllProducts.Close()

	rows, err := stmtListAllProducts.Query()
	if err != nil {
		return nil, err
	}
	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ProductID, &product.OwnerID, &product.ServerType, &product.InstanceID,
			&product.DateCreation, &product.DateTermination, &product.BillingCycle,
			&product.WalletID); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	return products, err
}

// For billing purposes
func ListActiveProductsByBillingCycle(billingCycle BillingCycle) ([]Product, error) {
	stmtListActiveProducts, err := sqlStatement(`SELECT ProductID, OwnerID, ServerType, InstanceID, DateCreation, BillingCycle, WalletID FROM dbprefix_products WHERE DateTermination = 0;`)
	if err != nil {
		return nil, err
	}
	defer stmtListActiveProducts.Close()

	rows, err := stmtListActiveProducts.Query()
	if err != nil {
		return nil, err
	}
	var products []Product
	for rows.Next() {
		var product Product = Product{
			DateTermination: time.Time{},
		}
		if err := rows.Scan(&product.ProductID, &product.OwnerID, &product.ServerType, &product.InstanceID,
			&product.DateCreation, &product.BillingCycle,
			&product.WalletID); err != nil {
			return products, err
		}
		products = append(products, product)
	}
	err = rows.Err()
	return products, err
}

func SelectProduct(productID uint64) (Product, error) {
	var product Product
	stmtSelectProduct, err := sqlStatement(`SELECT OwnerID, ServerType, InstanceID, DateCreation, DateTermination, BillingCycle, WalletID FROM dbprefix_products WHERE ProductID = ?;`)
	if err != nil {
		return product, err
	}
	defer stmtSelectProduct.Close()

	err = stmtSelectProduct.QueryRow().Scan(&product.OwnerID, &product.ServerType, &product.InstanceID,
		&product.DateCreation, &product.DateTermination, &product.BillingCycle, &product.WalletID)
	return product, err
}

// UpdateProduct() effectively updates OwnerID, ServerType, InstanceID, WalletID
// for the product specified by ProductID
func UpdateProduct(productID uint64, product Product) error {
	stmtUpdateProduct, err := sqlStatement(`UPDATE dbprefix_products SET
    OwnerID = ?, ServerType = ?, InstanceID = ?, WalletID = ?
    WHERE
    ProductID = ?;`)
	if err != nil {
		return err
	}
	defer stmtUpdateProduct.Close()

	_, err = stmtUpdateProduct.Exec(product.OwnerID, product.ServerType, product.InstanceID, product.WalletID, product.ProductID)
	return err
}

// Terminates the product
func TerminateProduct(ProductID uint64) error {
	stmtTerminateProduct, err := sqlStatement(`UPDATE dbprefix_products SET DateTermination = CURDATE() WHERE ProductID = ?;`)
	if err != nil {
		return err
	}
	defer stmtTerminateProduct.Close()

	_, err = stmtTerminateProduct.Exec(ProductID)
	return err
}

type BillingCycle uint8

const (
	PayAsYouGo BillingCycle = iota
	PerMonth
	PerQuarter
	PerYear
)
