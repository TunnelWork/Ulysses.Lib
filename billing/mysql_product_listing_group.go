package billing

const (
	productListingGroupTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_product_listing_group(
		product_group_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        product_group_name VARCHAR(64) NOT NULL,
        product_group_description TEXT NOT NULL,
        hidden BOOLEAN NOT NULL DEFAULT TRUE,
        PRIMARY KEY (product_group_id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`
)

/************ Product Listing Group Database ************/
func getProductListingGroupByID(id uint64) (ProductListingGroup, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_product_listing_group WHERE product_group_id = ?")
	if err != nil {
		return ProductListingGroup{}, err
	}
	defer stmt.Close()

	var productListingGroup ProductListingGroup
	err = stmt.QueryRow(id).Scan(&productListingGroup.ProductGroupID, &productListingGroup.ProductGroupName, &productListingGroup.ProductGroupDescription, &productListingGroup.Hidden)
	if err != nil {
		return ProductListingGroup{}, err
	}

	return productListingGroup, nil
}

func newProductListingGroup(productListingGroup ProductListingGroup) (uint64, error) {
	stmt, err := sqlStatement("INSERT INTO dbprefix_billing_product_listing_group (product_group_name, product_group_description, hidden) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(productListingGroup.ProductGroupName, productListingGroup.ProductGroupDescription, productListingGroup.Hidden)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func updateProductListingGroup(productListingGroup ProductListingGroup) error {
	stmt, err := sqlStatement("UPDATE dbprefix_billing_product_listing_group SET product_group_name = ?, product_group_description = ?, hidden = ? WHERE product_group_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(productListingGroup.ProductGroupName, productListingGroup.ProductGroupDescription, productListingGroup.Hidden, productListingGroup.ProductGroupID)
	if err != nil {
		return err
	}

	return nil
}

func deleteProductListingGroupByID(id uint64) error {
	stmt, err := sqlStatement("DELETE FROM dbprefix_billing_product_listing_group WHERE product_group_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
