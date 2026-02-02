package ports

// PasswordHasher define el puerto para el servicio de hash de passwords
// Aplica el patrón Strategy y el principio de Inversión de Dependencias
type PasswordHasher interface {
	// Compare verifica si un password en texto plano coincide con un hash
	Compare(hashedPassword, password string) error
}
