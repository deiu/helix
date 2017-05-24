package helix

import (
	"errors"
	"github.com/boltdb/bolt"
	"gopkg.in/hlandau/passlib.v1"
)

func (c *Context) registerPass(user, pass string) error {
	if len(user) == 0 || len(pass) == 0 {
		return errors.New("The username and password cannot be empty")
	}
	hash, _ := passlib.Hash(pass)
	// store pass hash
	return c.storePass(user, hash)
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

func (c *Context) storePass(user, pass string) error {
	if len(pass) == 0 {
		return errors.New("The password cannot be empty")
	}
	err := c.BoltDB.Update(func(tx *bolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists([]byte(user))
		if err != nil {
			return err
		}
		err = userBucket.Put([]byte("pass"), []byte(pass))
		return err
	})
	return err
}
