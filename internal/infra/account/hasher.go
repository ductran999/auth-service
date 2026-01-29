package account

import "github.com/alexedwards/argon2id"

type hasher struct{}

func NewHasher() *hasher {
	return &hasher{}
}

func (h *hasher) Hash(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func (h *hasher) ComparePasswordAndHash(pass, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pass, hash)
}
