package billing

import "errors"

type ProductListingGroup struct {
	ProductGroupID          uint64 `json:"product_group_id"`
	ProductGroupName        string `json:"product_group_name"`
	ProductGroupDescription string `json:"product_group_description"`
	Hidden                  bool   `json:"hidden"`
}

func GetProductListingGroupByID(id uint64) (ProductListingGroup, error) {
	return getProductListingGroupByID(id)
}

func NewProductListingGroup(plg ProductListingGroup) (uint64, error) {
	return newProductListingGroup(plg)
}

func (plg ProductListingGroup) Save() error {
	// validate all fields
	if plg.ProductGroupName == "" {
		return errors.New("billing: ProductGroupName is required")
	}
	if plg.ProductGroupDescription == "" {
		return errors.New("billing: ProductGroupDescription is required")
	}

	return updateProductListingGroup(plg)
}

func (plg ProductListingGroup) Delete() error {
	return deleteProductListingGroupByID(plg.ProductGroupID)
}

func DeleteProductListingGroupByID(id uint64) error {
	return deleteProductListingGroupByID(id)
}
