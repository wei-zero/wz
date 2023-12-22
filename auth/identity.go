package auth

type Identity struct {
	UserID   int64
	Roles    []string
	DeviceID string
	Internal string // internal source
}

func (i *Identity) IsAdmin() bool {
	for _, role := range i.Roles {
		if role == "admin" {
			return true
		}
	}
	return false
}
