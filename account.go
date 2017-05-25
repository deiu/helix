package helix

import (
	"encoding/json"
	"errors"

	"github.com/boltdb/bolt"
	"gopkg.in/hlandau/passlib.v1"
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

func (c *Context) verifyPass(user, pass string) (bool, error) {
	if len(user) == 0 || len(pass) == 0 {
		return false, errors.New("The username and password cannot be empty")
	}
	hash, err := c.getPass(user)
	if err != nil {
		return false, err
	}

	_, err = passlib.Verify(pass, hash)
	if err != nil {
		// incorrect password, malformed hash, etc.
		// either way, reject
		return false, err
	}

	// TODO: the context has decided, as per its policy, that
	// the hash which was used to validate the password
	// should be changed. It has upgraded the hash using
	// the verified password.
	// if newHash != "" {
	// 	c.storePass(user, newHash)
	// }

	return true, nil
}

func (c *Context) getPass(user string) (string, error) {
	hash := ""
	err := c.BoltDB.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(user))
		if userBucket == nil {
			return errors.New("Could not find a user bucket for " + user)
		}
		hash = string(userBucket.Get([]byte("pass")))
		return nil
	})
	return hash, err
}

func (c *Context) savePass(user, pass string) error {
	if len(pass) == 0 {
		return errors.New("The password cannot be empty")
	}

	hash, _ := passlib.Hash(pass)

	err := c.BoltDB.Update(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(user))
		if userBucket == nil {
			return errors.New("Could not find user bucket for " + user)
		}
		err := userBucket.Put([]byte("pass"), []byte(hash))
		return err
	})
	return err
}
