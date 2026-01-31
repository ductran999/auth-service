package security

type PasswordHasher interface {
	Hash(password string) (string, error)
	ComparePasswordAndHash(password, hash string) (bool, error)
}
