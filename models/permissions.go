package models

type Permissions struct {
	Admins []string
}

func (p *Permissions) IsAdmin(user *User) bool {
	if user == nil {
		return false
	}
	for _, username := range p.Admins {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (p *Permissions) IsSpaceAdmin(user *User, space *Space) bool {
	if user == nil {
		return false
	}
	if p.IsAdmin(user) {
		return true
	}
	for _, username := range space.Admins {
		if user.Username == username {
			return true
		}
	}
	return false
}
