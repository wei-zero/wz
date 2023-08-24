package mailer

type Mailer interface {
	SendByTemplate(from *Email, to *Email, subject string, templateID string, data map[string]string) error
}

// Email holds email name and address info
type Email struct {
	Name    string `json:"name,omitempty"`
	Address string `json:"email,omitempty"`
}

// NewEmail ...
func NewEmail(name string, address string) *Email {
	return &Email{
		Name:    name,
		Address: address,
	}
}
