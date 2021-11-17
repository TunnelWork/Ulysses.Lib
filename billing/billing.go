package billing

import (
	"fmt"
	"time"

	"github.com/TunnelWork/Ulysses.Lib/server"
)

// BillingCycle ENUMS
const (
	USAGE_BASED uint8 = iota + 1
	MONTHLY
	QUARTERLY
	ANNUALLY
)

// PricingPolicy ENUMS
const (
	PRICE_SUM uint8 = iota + 1
	PRICE_MAX
	PRICE_MIN
)

// HourlyUsageBilling() should be called at an hourly basis.
// It will bill all active usage-based billing products.
func HourlyUsageBilling() []error {
	monthToday := time.Now().Month()

	// Get all active Usage-based billing products
	usageBasedProducts, err := ListActiveProductsByBillingCycle(USAGE_BASED)
	if err != nil {
		return []error{err}
	}

	var errs []error

	// For each product, calculate new fee for this month.
	for _, product := range usageBasedProducts {
		// Check product listing for ServerType and InstanceID
		productListing, err := GetProductListingByID(product.ProductID)
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, GetProductListingByID() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// Get the server instance using server package
		serverInstance, err := server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, server.NewProvisioningServer() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// Get the current usage for this month
		account, err := serverInstance.GetAccount(product.SerialNumber())
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, serverInstance.GetAccount() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// Get resource usage of current month
		resources, err := account.Resources()
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, account.Resources() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// build resource usage map
		usageMap := make(map[uint64]float64)
		for _, resource := range resources {
			usageMap[resource.ResourceID] = resource.Used
		}

		// Calculate total price
		price, err := productListing.UsageBillingFactors.CalculateTotalPrice(usageMap)
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, productListing.UsageBillingFactors.CalculateTotalPrice() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// Bill the wallet
		err = product.CollectPayment(price)
		if err == ErrInsufficientFunds {
			// When insufficient funds, force billing and TERMINATE the server account.
			err = product.ForceCollectPayment(price)
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, product.ForceCollectPayment() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// TERMINATE the Server Account.
			err = serverInstance.DeleteAccount(product.SerialNumber())
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, serverInstance.DeleteAccount() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// TERMINATE the Product
			err = product.Terminate()
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, product.Terminate() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
			}

			// TODO: Notify admin and product owner.
			continue // Continue to next product
		} else if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, product.CollectPayment() returned error: %s", product.SerialNumber(), err.Error())
			errs = append(errs, err)
			continue
		}

		// For usage-based billing, dateLastBill is used to track the "month"
		// of last usage updates
		if monthToday != product.dateLastBill.Month() {
			// First time of this month.
			// Reset usage for server account to 0.
			err = serverInstance.RefreshAccount(product.SerialNumber())
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, serverInstance.RefreshAccount() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// Reset monthly spending statistics of the product
			product.BillingOption.CurrentMonthSpending = 0
			err = product.Save()
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, product.Save() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// Update dateLastBill
			err = updateDateLastBillProductBySerialNumber(product.SerialNumber())
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, updateDateLastBillProductBySerialNumber() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}
		}
	}
	return errs
}

// DailyRecurringBilling() should be called at (at least) a daily basis.
// It will bill all active recurring billing products that is due.
func DailyRecurringBilling() []error {
	var errs []error

	// Get all active recurring billing products
	monthlyProducts, err := ListActiveProductsByBillingCycle(MONTHLY)
	if err != nil {
		err = fmt.Errorf("billing: ListActiveProductsByBillingCycle(MONTHLY) returned error: %s", err.Error())
		errs = append(errs, err)
	} else {
		// Bill Monthly Products if needed
		monthlyErrs := dailyMonthlyBilling(monthlyProducts)
		errs = append(errs, monthlyErrs...)
	}

	quarterlyProducts, err := ListActiveProductsByBillingCycle(QUARTERLY)
	if err != nil {
		err = fmt.Errorf("billing: ListActiveProductsByBillingCycle(QUARTERLY) returned error: %s", err.Error())
		errs = append(errs, err)
	} else {
		// Bill Quarterly Products if needed
		quarterlyErrs := dailyQuarterlyBilling(quarterlyProducts)
		errs = append(errs, quarterlyErrs...)
	}

	annuallyProducts, err := ListActiveProductsByBillingCycle(ANNUALLY)
	if err != nil {
		err = fmt.Errorf("billing: ListActiveProductsByBillingCycle(ANNUALLY) returned error: %s", err.Error())
		errs = append(errs, err)
	} else {
		// Bill Annually Products if needed
		annuallyErrs := dailyAnnuallyBilling(annuallyProducts)
		errs = append(errs, annuallyErrs...)
	}

	return errs
}

// HourlyProductTermination() should be called at an hourly basis.
// It will pickup all products that is due to be terminated:
// - If usage-based, charge the final amount
// - Terminate Account
// - Terminate Product
func HourlyProductTermination() []error {
	var errs []error

	// Get all active products that is due to be terminated
	products, err := ListProductsToTerminate()
	if err != nil {
		err = fmt.Errorf("billing: ListProductsToTerminate() returned error: %s", err.Error())
		errs = append(errs, err)
	} else {
		// terminate all products in the slice
		for _, product := range products {
			// Check product listing for ServerType and InstanceID
			productListing, err := GetProductListingByID(product.ProductID)
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, GetProductListingByID() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// Get the server instance using server package
			serverInstance, err := server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, server.NewProvisioningServer() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// If the product is Usage-based, charge the final amount before terminating
			if product.BillingOption.BillingCycle == USAGE_BASED {
				// Get the current usage for this month
				account, err := serverInstance.GetAccount(product.SerialNumber())
				if err != nil {
					err = fmt.Errorf("billing: ProductSN==%d, serverInstance.GetAccount() returned error: %s", product.SerialNumber(), err.Error())
					errs = append(errs, err)
					continue
				}

				// Get resource usage of current month
				resources, err := account.Resources()
				if err != nil {
					err = fmt.Errorf("billing: ProductSN==%d, account.Resources() returned error: %s", product.SerialNumber(), err.Error())
					errs = append(errs, err)
					continue
				}

				// build resource usage map
				usageMap := make(map[uint64]float64)
				for _, resource := range resources {
					usageMap[resource.ResourceID] = resource.Used
				}

				// Calculate total price
				price, err := productListing.UsageBillingFactors.CalculateTotalPrice(usageMap)
				if err != nil {
					err = fmt.Errorf("billing: ProductSN==%d, productListing.UsageBillingFactors.CalculateTotalPrice() returned error: %s", product.SerialNumber(), err.Error())
					errs = append(errs, err)
					continue
				}

				err = product.ForceCollectPayment(price)
				if err != nil {
					err = fmt.Errorf("billing: ProductSN==%d, product.ForceCollectPayment() returned error: %s", product.SerialNumber(), err.Error())
					errs = append(errs, err)
					continue
				}
			}

			// TERMINATE the Server Account.
			err = serverInstance.DeleteAccount(product.SerialNumber())
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, serverInstance.DeleteAccount() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
				continue
			}

			// TERMINATE the Product
			err = product.Terminate()
			if err != nil {
				err = fmt.Errorf("billing: ProductSN==%d, product.Terminate() returned error: %s", product.SerialNumber(), err.Error())
				errs = append(errs, err)
			}
		}
	}

	return errs
}

func dailyMonthlyBilling(products []*Product) []error {
	var errs []error

	dateToday := time.Now()
	dateOneMonthAgo := dateToday.AddDate(0, -1, 0)

	for _, product := range products {
		// Check if the product is due for billing
		if product.dateLastBill.Before(dateOneMonthAgo) || product.dateLastBill.Equal(dateOneMonthAgo) { // due for billing
			err := processRecurringBilling(product)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	return errs
}

func dailyQuarterlyBilling(products []*Product) []error {
	var errs []error

	dateToday := time.Now()
	dateOneQuarterAgo := dateToday.AddDate(0, -3, 0)

	for _, product := range products {
		// Check if the product is due for billing
		if product.dateLastBill.Before(dateOneQuarterAgo) || product.dateLastBill.Equal(dateOneQuarterAgo) { // due for billing
			err := processRecurringBilling(product)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	return errs
}

func dailyAnnuallyBilling(products []*Product) []error {
	var errs []error

	dateToday := time.Now()
	dateOneYearAgo := dateToday.AddDate(-1, 0, 0)

	for _, product := range products {
		// Check if the product is due for billing
		if product.dateLastBill.Before(dateOneYearAgo) || product.dateLastBill.Equal(dateOneYearAgo) { // due for billing
			err := processRecurringBilling(product)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}
	}

	return errs
}

func processRecurringBilling(product *Product) error {
	// Check product listing for ServerType and InstanceID
	productListing, err := GetProductListingByID(product.ProductID)
	if err != nil {
		err = fmt.Errorf("billing: ProductSN==%d, GetProductListingByID() returned error: %s", product.SerialNumber(), err.Error())
		return err
	}

	// Get the server instance using server package
	serverInstance, err := server.NewProvisioningServer(productListing.ServerType, productListing.ServerInstanceID, productListing.ServerConfiguration)
	if err != nil {
		err = fmt.Errorf("billing: ProductSN==%d, server.NewProvisioningServer() returned error: %s", product.SerialNumber(), err.Error())
		return err
	}

	// Collect payment from the wallet
	err = product.CollectPayment() // No parameter for recurrings
	if err == ErrInsufficientFunds {
		// If insufficient funds to renew:
		// - Suspend the account
		// - Label the product to be terminated tomorrow

		err = serverInstance.SuspendAccount(product.SerialNumber())
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, serverInstance.DeleteAccount() returned error: %s", product.SerialNumber(), err.Error())
			return err
		}

		err = product.ToTerminateOn(time.Now().AddDate(0, 0, 1))
		if err != nil {
			err = fmt.Errorf("billing: ProductSN==%d, product.Terminate() returned error: %s", product.SerialNumber(), err.Error())
			return err
		}

		// TODO: Notify admin and product owner.
		return nil
	} else if err != nil {
		err = fmt.Errorf("billing: ProductSN==%d, product.CollectPayment() returned error: %s", product.SerialNumber(), err.Error())
		return err
	}
	// Lastly, refresh account for the renewed product
	err = serverInstance.RefreshAccount(product.SerialNumber())
	if err != nil {
		err = fmt.Errorf("billing: ProductSN==%d, serverInstance.RefreshAccount() returned error: %s", product.SerialNumber(), err.Error())
		return err
	}
	return nil
}
