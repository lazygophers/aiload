package aiload

func (p UserRole) Accessible(role UserRole) bool {
	switch p {
	case UserRole_Public:
		return true

	case UserRole_User:
		return role == UserRole_User || role == UserRole_Admin

	case UserRole_Admin:
		return role == UserRole_Admin

	}
	return false
}
