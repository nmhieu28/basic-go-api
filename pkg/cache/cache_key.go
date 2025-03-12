package cache

import (
	"fmt"

	"github.com/google/uuid"
)

func ConfirmAccountKey(userId uuid.UUID) string {
	return fmt.Sprintf("identity:otp_code:%s:email_confirm", userId)
}

func RefreshTokenKey(userId uuid.UUID) string {
	return fmt.Sprintf("identity:refresh_token:%s", userId)
}

func ForgotPasswordKey(userId uuid.UUID) string {
	return fmt.Sprintf("identity:forgot_password:%s", userId)
}
