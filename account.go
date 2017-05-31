package helix

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
)

type User struct {
	Username string
	Email    string
}

func NewUser() *User {
	return &User{}
}

func (c *Context) addUser(user, pass, email string) error {
	if len(user) == 0 || len(pass) == 0 || len(email) == 0 {
		return errors.New("The username and password cannot be empty")
	}
	u := NewUser()
	u.Username = user
	u.Email = email

	// store the new user
	err := c.saveUser(u)
	if err != nil {
		return err
	}
	// store the new pass
	return c.savePass(user, pass)
}

func (c *Context) getUser(username string) (*User, error) {
	user := NewUser()

	err := c.BoltDB.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(username))
		if userBucket == nil {
			return errors.New("Could not find a user bucket for " + username)
		}
		err := json.Unmarshal(userBucket.Get([]byte("user")), user)
		return err
	})
	return user, err
}

func (c *Context) saveUser(user *User) error {
	err := c.BoltDB.Update(func(tx *bolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists([]byte(user.Username))
		if err != nil {
			return err
		}
		// No need to handle error since we only have strings in the user struct
		buf, _ := json.Marshal(user)
		err = userBucket.Put([]byte("user"), buf)
		return err
	})
	return err
}
