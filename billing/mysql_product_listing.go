package billing

import (
	"database/sql"
	"encoding/json"
	"errors"
)

const (
	productListingTblCreation = `CREATE TABLE IF NOT EXISTS dbprefix_billing_product_listing(
        product_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
        product_group_id BIGINT UNSIGNED NOT NULL,
        product_name VARCHAR(64) NOT NULL,
        product_description TEXT NOT NULL,
        server_type VARCHAR(64) NOT NULL,
        server_instance_id VARCHAR(128) NOT NULL,
        server_configuration TEXT NOT NULL,
        billing_options TEXT NOT NULL,
        usage_billing_factors TEXT NOT NULL,
        hidden BOOLEAN NOT NULL DEFAULT TRUE,
        discontinued BOOLEAN NOT NULL DEFAULT TRUE,
        PRIMARY KEY (product_id),
        INDEX (product_group_id),
        CONSTRAINT FOREIGN KEY (product_group_id) REFERENCES dbprefix_billing_product_listing_group(product_group_id) ON DELETE RESTRICT
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;`
)

var (
	ErrProductListingNotFound = errors.New("billing: product listing not found")
)

/************ Product Listing Database ************/
func getProductListingByID(productID uint64, allowDiscontinued bool) (*ProductListing, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_product_listing WHERE product_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var productListing ProductListing
	var serverConfigurationJson string
	var billingOptionsJson string
	var usageBillingFactorsJson string
	err = stmt.QueryRow(productID).Scan(
		&productListing.productID,
		&productListing.ProductGroupID,
		&productListing.ProductName,
		&productListing.ProductDescription,
		&productListing.ServerType,
		&productListing.ServerInstanceID,
		&serverConfigurationJson,
		&billingOptionsJson,
		&usageBillingFactorsJson,
		&productListing.hidden,
		&productListing.discontinued,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProductListingNotFound
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(billingOptionsJson), &productListing.BillingOptions); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(usageBillingFactorsJson), &productListing.UsageBillingFactors); err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(serverConfigurationJson), &productListing.ServerConfiguration); err != nil {
		return nil, err
	}

	// will not show discontinued products unless explicitly allowed
	if !allowDiscontinued && productListing.discontinued {
		return nil, ErrProductListingNotFound
	}

	return &productListing, nil
}

func getProductListingsByGroupID(productGroupID uint64, includeHiddenDiscontinued bool) ([]*ProductListing, error) {
	stmt, err := sqlStatement("SELECT * FROM dbprefix_billing_product_listing WHERE product_group_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var productListings []*ProductListing
	var serverConfigurationJson string
	var billingOptionsJson string
	var usageBillingFactorsJson string
	rows, err := stmt.Query(productGroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var productListing ProductListing
		err = rows.Scan(
			&productListing.productID,
			&productListing.ProductGroupID,
			&productListing.ProductName,
			&productListing.ProductDescription,
			&productListing.ServerType,
			&productListing.ServerInstanceID,
			&serverConfigurationJson,
			&billingOptionsJson,
			&usageBillingFactorsJson,
			&productListing.hidden,
			&productListing.discontinued,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(billingOptionsJson), &productListing.BillingOptions); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(usageBillingFactorsJson), &productListing.UsageBillingFactors); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(serverConfigurationJson), &productListing.ServerConfiguration); err != nil {
			return nil, err
		}

		// will not show discontinued products unless explicitly allowed
		if !includeHiddenDiscontinued && productListing.discontinued {
			continue
		}

		productListings = append(productListings, &productListing)
	}

	return productListings, nil
}

func addProductListing(pl *ProductListing) (uint64, error) {
	stmt, err := sqlStatement(`INSERT INTO dbprefix_billing_product_listing (
        product_group_id, 
        product_name, 
        product_description, 
        server_type, 
        server_instance_id, 
		server_configuration,
        billing_options, 
        usage_billing_factors, 
        hidden, 
        discontinued
    ) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	billingOptionsJson, err := json.Marshal(pl.BillingOptions)
	if err != nil {
		return 0, err
	}

	usageBillingFactorsJson, err := json.Marshal(pl.UsageBillingFactors)
	if err != nil {
		return 0, err
	}

	serverConfigurationJson, err := json.Marshal(pl.ServerConfiguration)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(
		pl.ProductGroupID,
		pl.ProductName,
		pl.ProductDescription,
		pl.ServerType,
		pl.ServerInstanceID,
		string(serverConfigurationJson),
		string(billingOptionsJson),
		string(usageBillingFactorsJson),
		pl.hidden,
		pl.discontinued,
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
}

func updateProductListing(pl *ProductListing) error {
	stmt, err := sqlStatement(`UPDATE dbprefix_billing_product_listing SET
        product_group_id = ?,
        product_name = ?,
        product_description = ?,
        server_type = ?,
        server_instance_id = ?,
		server_configuration = ?,
        billing_options = ?,
        usage_billing_factors = ?,
        hidden = ?,
        discontinued = ?
    WHERE product_id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	billingOptionsJson, err := json.Marshal(pl.BillingOptions)
	if err != nil {
		return err
	}

	usageBillingFactorsJson, err := json.Marshal(pl.UsageBillingFactors)
	if err != nil {
		return err
	}

	serverConfigurationJson, err := json.Marshal(pl.ServerConfiguration)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		pl.ProductGroupID,
		pl.ProductName,
		pl.ProductDescription,
		pl.ServerType,
		pl.ServerInstanceID,
		string(serverConfigurationJson),
		string(billingOptionsJson),
		string(usageBillingFactorsJson),
		pl.hidden,
		pl.discontinued,
		pl.productID,
	)
	return err
}

func deleteProductListingByID(productID uint64) error {
	stmt, err := sqlStatement("DELETE FROM dbprefix_billing_product_listing WHERE product_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(productID)
	return err
}
