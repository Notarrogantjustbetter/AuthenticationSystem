package database

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
	ErrUsernameTaken = errors.New("username is taken")
)

type User struct {
	id int64
}

type UserInterface interface {
	Authenticate(password string) error
	GetId() (int64, error)
	GetNickname() (string, error)
	GetHash() ([]byte, error)
}

func NewUser(nickname string, hash []byte) (*User, error) {
	usernameExists, _ := client.HExists("user:by-username", nickname).Result()
	if usernameExists {
		return nil, ErrUsernameTaken
	}
	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline()
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "nickname", nickname)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", nickname, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &User{id}, nil
}

func (u *User) GetId() (int64, error) {
	return u.id, nil
}

func (u *User) GetNickname() (string, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return client.HGet(key, "nickname").Result()
}

func (u *User) GetHash() ([]byte, error) {
	key := fmt.Sprintf("user:%d", u.id)
	return client.HGet(key, "hash").Bytes()
}

func (u *User) Authenticate(password string) error {
	hash, err := u.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}

func GetUserById(id int64) (*User, error) {
	return &User{id}, nil
}

func GetUserByNickname(nickname string) (*User, error) {
	id, err := client.HGet("user:by-username", nickname).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	return GetUserById(id)
}

func LoginUser(nickname, password string) (*User, error) {
	user, err := GetUserByNickname(nickname)
	if err != nil {
		return nil, err
	}
	return user, user.Authenticate(password)
}

func RegisterUser(username, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}
