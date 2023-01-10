package models

type Permissions struct {
	Admins      []string
	SpaceAdmins map[string][]string
}

func (p *Permissions) IsAdmin(user *User) bool {
	if user == nil {
		return false
	}
	for _, id := range p.Admins {
		if user.Username == id {
			return true
		}
	}
	return false
}

func (p *Permissions) IsSpaceAdmin(spaceID string, user *User) bool {
	if user == nil {
		return false
	}
	if p.IsAdmin(user) {
		return true
	}
	for _, id := range p.SpaceAdmins[spaceID] {
		if user.Username == id {
			return true
		}
	}
	return false
}
