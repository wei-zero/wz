package sendcloud

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/wei-zero/wz/mailer"
)

func TestMailer(t *testing.T) {
	godotenv.Load()
	m := NewMailer(os.Getenv("SENDCLOUD_API_USER"), os.Getenv("SENDCLOUD_API_KEY"))
	from := &mailer.Email{
		Name:    "Xiaoma service",
		Address: "service@xiaoma.ai",
	}
	to := &mailer.Email{
		Name:    "39419358",
		Address: "39419358@qq.com",
	}
	err := m.SendByTemplate(from, to, "test subject", "xm_user_register", map[string]string{
		"activate_code": "123456",
	})
	assert.NoError(t, err)
}
