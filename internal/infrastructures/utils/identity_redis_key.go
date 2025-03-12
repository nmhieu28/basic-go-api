package identity_redis_key

import "fmt"

func OtpRegisterKey(userID string) string {
	return fmt.Sprintf("user_otp:%s:register", userID)
}

func ForgotPasswordKey(userID string) string {
	return fmt.Sprintf("user_password_otp:%s:forgot_password", userID)
}
