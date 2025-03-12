package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"backend/pkg/logger"
)

const (
	CSRFHeader = "X-CSRF-Token"
	csrfSalt   = "8HHsTtto01InHHcIBo69W8t3qHMkNh8n"
)

func MakeToken(sid string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sid)
	if err != nil {
		logger.Errorf("Make CSRF Token", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}
