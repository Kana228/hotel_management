package session


// CurrentUser holds the pointer to the currently logged-in user.
// It can store either a *models.Customer or *models.Vendor.
var CurrentUser interface{}

// SetCurrentUser sets the current user (customer or vendor) after login.
func SetCurrentUser(user interface{}) {
	CurrentUser = user
}

// GetCurrentUser returns the current logged-in user.
func GetCurrentUser() interface{} {
	return CurrentUser
}

// ClearCurrentUser clears the current session.
func ClearCurrentUser() {
	CurrentUser = nil
}

