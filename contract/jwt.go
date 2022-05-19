package contract

import "gitlab.com/gocastsian/writino/entity"

type (
	GenerateTokenPair func(secret []byte, user entity.User) (map[string]string, error)
	ParseToken        func(secret []byte, token string) (entity.User, error)
	ParseRefToken     func(secret []byte, token string) (string, error)
)
