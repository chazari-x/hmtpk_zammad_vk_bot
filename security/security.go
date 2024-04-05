package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"strings"

	"github.com/chazari-x/hmtpk_zammad_vk_bot/config"
)

type Security struct {
	cfg config.Security
}

func NewSecurity(cfg config.Security) *Security {
	return &Security{cfg: cfg}
}

// CreateHmacSignature Функция для создания подписи HMAC SHA1
func (s *Security) CreateHmacSignature(data []byte) string {
	h := hmac.New(sha1.New, []byte(s.cfg.SecretKey))
	h.Write(data)
	return strings.ReplaceAll(hex.EncodeToString(h.Sum(nil)), "_", "")
}
