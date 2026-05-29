package principal

func (p *Principal) HasRole(
	role string,
) bool {

	for _, r := range p.Roles {
		if r == role {
			return true
		}
	}

	return false
}

func (p *Principal) HasPermission(
	permission string,
) bool {

	for _, p := range p.Permissions {
		if p == permission {
			return true
		}
	}

	return false
}
