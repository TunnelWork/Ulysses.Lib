package billing

import (
	"encoding/json"
	"errors"
)

type ProductListing struct {
	// Client Viewable
	productID          uint64 // productID is actually the Serial Number of Product Listing
	ProductGroupID     uint64 `json:"product_group_id"`
	ProductName        string `json:"product_name"`
	ProductDescription string `json:"product_description"`

	// Internal BizLogic Related
	ServerType          string                 `json:"server_type"`
	ServerInstanceID    string                 `json:"server_instance_id"`
	ServerConfiguration map[string]interface{} `json:"server_configuration"`  // a big chunk of JSON
	BillingOptions      []BillingOption        `json:"billing_options"`       // Stored as JSON Arr
	UsageBillingFactors UsageBillingFactors    `json:"usage_billing_factors"` // Store as JSON Obj

	// Internal API Behavior Related
	hidden       bool // default: true
	discontinued bool // default: true
}

func (pl *ProductListing) ProductID() uint64 {
	return pl.productID
}

func (pl *ProductListing) Hidden() bool {
	return pl.hidden
}

func (pl *ProductListing) Discontinued() bool {
	return pl.discontinued
}

// For Customer Purchasing. Throwing error for discontinued product
func GetProductListingByID(productID uint64) (*ProductListing, error) {
	return getProductListingByID(productID, false)
}

// For Admin Editing/Customer Viewing Existing. Will show discontinued correctly
func SudoGetProductListingByID(productID uint64) (*ProductListing, error) {
	return getProductListingByID(productID, true)
}

// For Customer. Not showing hidden/disconiued products
func GetProductListingsByGroupID(productGroupID uint64) ([]*ProductListing, error) {
	return getProductListingsByGroupID(productGroupID, false)
}

// For Admin. Include hidden/discontinued products
func SudoGetProductListingsByGroupID(productGroupID uint64) ([]*ProductListing, error) {
	return getProductListingsByGroupID(productGroupID, true)
}

func (pl *ProductListing) Add() (uint64, error) {
	return AddProductListing(pl)
}

func (pl *ProductListing) Save() error {
	return UpdateProductListing(pl)
}

func (pl *ProductListing) Delete() error {
	return DeleteProductListingByID(pl.productID)
}

// Hide() prevents the product from being shown by listing
func (pl *ProductListing) Hide() error {
	if pl.hidden {
		return nil
	} else {
		pl.hidden = true
		return UpdateProductListing(pl)
	}
}

// Unhide() undo the Hide()
func (pl *ProductListing) Unhide() error {
	if pl.BillingOptions == nil || len(pl.BillingOptions) == 0 {
		return errors.New("billing: billing options not set")
	}
	if pl.UsageBillingFactors.PricingPolicy == 0 || len(pl.UsageBillingFactors.Factors) == 0 {
		return errors.New("billing: usage billing factors not set")
	}
	if !pl.hidden {
		return nil
	} else {
		pl.hidden = false
		return UpdateProductListing(pl)
	}
}

// Discontinue() prevents the product from being purchased
func (pl *ProductListing) Discontinue() error {
	if pl.discontinued {
		return nil
	} else {
		pl.discontinued = true
		return UpdateProductListing(pl)
	}
}

// Reactivate() undo the Discontinue()
func (pl *ProductListing) Reactivate() error {
	if pl.BillingOptions == nil || len(pl.BillingOptions) == 0 {
		return errors.New("billing: billing options not set")
	}
	if pl.UsageBillingFactors.PricingPolicy == 0 || len(pl.UsageBillingFactors.Factors) == 0 {
		return errors.New("billing: usage billing factors not set")
	}
	if !pl.discontinued {
		return nil
	} else {
		pl.discontinued = false
		return UpdateProductListing(pl)
	}
}

func (pl *ProductListing) CreateProduct(ownerUserID uint64, ownerAffiliationID uint64, BillingCycle uint8, WalletID uint64) (*Product, error) {
	if pl.discontinued {
		return nil, errors.New("billing: product is discontinued")
	}

	// Find matching BillingOption
	var billingOption BillingOption
	for _, bo := range pl.BillingOptions {
		if bo.BillingCycle == BillingCycle {
			billingOption = bo
			break
		}
	}

	if billingOption.Price == 0 && billingOption.MonthlySpendingCap == 0 {
		return nil, errors.New("billing: matching billing option not found") // if price == 0, it must be usage-based pricing.
	}

	var product Product = Product{
		OwnerUserID:        ownerUserID,
		OwnerAffiliationID: ownerAffiliationID,
		ProductID:          pl.ProductID(),
		WalletID:           WalletID,
		BillingOption:      billingOption,
	}
	return &product, nil
}

func AddProductListing(pl *ProductListing) (uint64, error) {
	// validate all fields set
	if pl.ProductName == "" {
		return 0, errors.New("billing: product name not set")
	}
	if pl.ProductDescription == "" {
		return 0, errors.New("billing: product description not set")
	}
	if pl.ServerType == "" {
		return 0, errors.New("billing: server type not set")
	}
	if pl.ServerInstanceID == "" {
		return 0, errors.New("billing: server instance ID not set")
	}
	if pl.BillingOptions == nil || len(pl.BillingOptions) == 0 {
		pl.hidden = true
		pl.discontinued = true
	}
	if pl.UsageBillingFactors.PricingPolicy == 0 || len(pl.UsageBillingFactors.Factors) == 0 {
		pl.hidden = true
		pl.discontinued = true
	}

	return addProductListing(pl)
}

func UpdateProductListing(pl *ProductListing) error {
	// validate all fields set
	if pl.ProductName == "" {
		return errors.New("billing: product name not set")
	}
	if pl.ProductDescription == "" {
		return errors.New("billing: product description not set")
	}
	if pl.ServerType == "" {
		return errors.New("billing: server type not set")
	}
	if pl.ServerInstanceID == "" {
		return errors.New("billing: server instance ID not set")
	}

	return updateProductListing(pl)
}

func DeleteProductListingByID(productID uint64) error {
	return deleteProductListingByID(productID)
}

type BillingOption struct {
	BillingCycle         uint8   `json:"billing_cycle"` // 0 - Usage-based, 1 - Monthly, 2 - Quarterly, 3 - Annually
	Price                float64 `json:"price"`
	MonthlySpendingCap   float64 `json:"monthly_spending_cap"`   // Only for Usage-based billing
	CurrentMonthSpending float64 `json:"current_month_spending"` // Only for Usage-based billing
}

func BillingOptionsFromJSON(jsonStr string) ([]BillingOption, error) {
	var billingOptions []BillingOption
	err := json.Unmarshal([]byte(jsonStr), &billingOptions)
	if err != nil {
		return nil, err
	}
	return billingOptions, nil
}

func BillingOptionsToJSON(billingOptions []BillingOption) (string, error) {
	jsonBytes, err := json.Marshal(billingOptions)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

type UsageBillingFactors struct {
	PricingPolicy uint8                           `json:"pricing_policy"` // 0 - SUM(GroupPrice...), 1 - MAX(GroupPrice...), 2 - MIN(GroupPrice...)
	Factors       map[uint8][]*UsageBillingFactor `json:"factors"`        // key: BillingGroupID
}

func UsageBillingFactorsFromJSON(jsonStr string) (*UsageBillingFactors, error) {
	var factors UsageBillingFactors
	err := json.Unmarshal([]byte(jsonStr), &factors)
	if err != nil {
		return nil, err
	}
	return &factors, nil
}

func (ubfs *UsageBillingFactors) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(ubfs)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// for resUsageMap, uint64 is the ResourceID, float64 is the usage
func (ubfs *UsageBillingFactors) CalculateTotalPrice(resUsageMap map[uint64]float64) (float64, error) {
	switch ubfs.PricingPolicy {
	case PRICE_SUM:
		return ubfs.sumTotalPrice(resUsageMap)
	case PRICE_MAX:
		return ubfs.maxTotalPrice(resUsageMap)
	case PRICE_MIN:
		return ubfs.minTotalPrice(resUsageMap)
	default:
		return 0, errors.New("billing: invalid pricing policy")
	}
}

func (ubfs *UsageBillingFactors) sumTotalPrice(resUsageMap map[uint64]float64) (float64, error) {
	totalPrice := 0.0
	for _, group := range ubfs.Factors { // iterate through each group
		var groupPrice float64
		for _, factor := range group { // iterate through each factor
			// find the resource usage in the map
			if usage, ok := resUsageMap[factor.ResourceID]; ok {
				groupPrice += factor.CalculatePrice(usage)
			} else {
				return 0, errors.New("billing: resource usage not found")
			}
		}
		totalPrice += groupPrice
	}
	return totalPrice, nil
}

func (ubfs *UsageBillingFactors) maxTotalPrice(resUsageMap map[uint64]float64) (float64, error) {
	totalPrice := 0.0
	for _, group := range ubfs.Factors { // iterate through each group
		var groupPrice float64
		for _, factor := range group { // iterate through each factor
			// find the resource usage in the map
			if usage, ok := resUsageMap[factor.ResourceID]; ok {
				groupPrice += factor.CalculatePrice(usage)
			} else {
				return 0, errors.New("billing: resource usage not found")
			}
		}
		if groupPrice > totalPrice {
			totalPrice = groupPrice
		}
	}
	return totalPrice, nil
}

func (ubfs *UsageBillingFactors) minTotalPrice(resUsageMap map[uint64]float64) (float64, error) {
	if len(ubfs.Factors) == 0 {
		return 0, errors.New("billing: no factors")
	}
	// set totalPrice to the group price of the first group
	var firstGroupPrice float64
	for _, group := range ubfs.Factors { // iterate through first group
		for _, factor := range group {
			// find the resource usage in the map
			if usage, ok := resUsageMap[factor.ResourceID]; ok {
				firstGroupPrice += factor.CalculatePrice(usage)
			} else {
				return 0, errors.New("billing: resource usage not found")
			}
		}
		break
	}
	var totalPrice = firstGroupPrice

	for _, group := range ubfs.Factors { // iterate through each group
		var groupPrice float64
		for _, factor := range group { // iterate through each factor
			// find the resource usage in the map
			if usage, ok := resUsageMap[factor.ResourceID]; ok {
				groupPrice += factor.CalculatePrice(usage)
			} else {
				return 0, errors.New("billing: resource usage not found")
			}
		}
		if groupPrice < totalPrice {
			totalPrice = groupPrice
		}
	}
	return totalPrice, nil
}

type UsageBillingFactor struct {
	BillingGroupID uint8   `json:"billing_group_id"`
	ResourceID     uint64  `json:"resource_id"` // see: Ulysses.Lib/server/resource.go
	UnitPrice      float64 `json:"unit_price"`  // UnitPrice * Used = Price
}

func (ubf *UsageBillingFactor) CalculatePrice(used float64) float64 {
	return ubf.UnitPrice * used
}
