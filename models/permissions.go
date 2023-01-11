package models

// TODO move space users to db
type Permissions struct {
	Admins      []string
	SpaceAdmins map[string][]string
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

func (p *Permissions) IsSpaceAdmin(spaceID string, user *User) bool {
	if user == nil {
		return false
	}
	if p.IsAdmin(user) {
		return true
	}
	for _, username := range p.SpaceAdmins[spaceID] {
		if user.Username == username {
			return true
		}
	}
	return false
}

func (p *Permissions) UserSpaces(user *User) []string {
	if user == nil {
		return nil
	}
	var spaceIDs []string
	for spaceID, usernames := range p.SpaceAdmins {
		for _, username := range usernames {
			if user.Username == username {
				spaceIDs = append(spaceIDs, spaceID)
				break
			}
		}
	}
	return spaceIDs
}
