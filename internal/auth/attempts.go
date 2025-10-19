package auth

import (
	"fmt"
	"time"
)

func HandleFailedAttempt(attempt int, maxAttempts int) bool {
	if attempt >= maxAttempts {
		fmt.Println("Demasiados intentos fallidos. Conexión bloqueada.")
		return false
	}

	wait := time.Duration(attempt*2) * time.Second
	fmt.Printf("Intento fallido %d/%d. Esperando %v antes de reintentar...\n", attempt, maxAttempts, wait)
	time.Sleep(wait)
	return true
}
