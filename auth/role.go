package auth

// All accesses to Ulysses are role-based. An account must be assigned relevant roles to gain access.
//
// 		- All roles are uniquely and positively defined, meaning that a user may only gain
// 		access when they have a new role.
// 		- On the other hand, if a user does not have the relevant role, they don't have
//		access to the resource even if they are authorized for a higher permission.
// 		i.e., An ADMIN may not have READ access if they don't have USER role.

type Role uint32

// Known roles as unambiguous binary flags allowing cascading
const (
	ROLELESS Role = 0

	/************ Global Role ************/
	GLOBAL_EVALUATION_USER Role = 1 << (iota - 1) // EVALUATION_USER is a global role. In principle it is mutual exclusive against PRODUCTION_USER.
	GLOBAL_PRODUCTION_USER                        // PRODUCTION_USER is a global role. In principle it is mutual exclusive against EVALUATION_USER.
	GLOBAL_INTERNAL_USER                          // INTERNAL_USER may order products free of charge
	GLOBAL_ADMIN                                  // ADMIN owns all access to management interface

	/************ Exemptional Role ************/
	EXEMPT_MARKETING_CONTACT // User won't be contacted for marketing purposes
	EXEMPT_BILLING_CONTACT   // User won't be notified for billing updates
	EXEMPT_SUPPORT_CONTACT   // User won't be notified for supporting case updates

	/************ Affiliation Role ************/
	// Affiliations (enterprises) may purchase products and set them
	// to be shared by users
	AFFILIATION_ACCOUNT_USER  // ACCOUNT_USER is a user belong to an enterprise
	AFFILIATION_ACCOUNT_ADMIN // ACCOUNT_ADMIN may create users and manage users (assigning roles, etc)

	AFFILIATION_PRODUCT_USER  // PRODUCT_USER may only view(and use) products
	AFFILIATION_PRODUCT_ADMIN // PRODUCT_ADMIN may create and edit shared products

	AFFILIATION_BILLING_USER  // BILLING_USER may purchase products with Affiliation-owned wallet
	AFFILIATION_BILLING_ADMIN // BILLING_ADMIN may deposit funds into Affiliation-owned wallet and view/manage associated products
)

// Roles() merge input roles into one single role.
// repeated entry will be ignored.
func Roles(roles ...Role) Role {
	var role Role
	for _, r := range roles {
		if role.Includes(r) { // never add repeated role
			continue
		}
		role |= r
	}
	return role
}

// Includes() checks if the input role is included in current role.
func (r Role) Includes(other Role) bool {
	return r&other == other
}

// AddRole() add a role to the current role.
// repeated entry will be ignored.
func (r Role) AddRole(role Role) Role {
	return r | role
}

// RemoveRole() remove a role from the current role.
func (r Role) RemoveRole(role Role) Role {
	return r &^ role
}
