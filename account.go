package helix

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gocraft/web"
)

type User struct {
	Username string
	Email    string
}

func NewUser() *User {
	return &User{}
}

func (c *Context) GetAccountHandler(w web.ResponseWriter, req *web.Request) {
	if len(c.User) == 0 {
		c.AuthenticationRequired(w, req)
	}
}

func (c *Context) LoginHandler(w web.ResponseWriter, req *web.Request) {
	ok, err := c.verifyPass(req.FormValue("username"), req.FormValue("password"))
	if err != nil {
		logger.Info().Msg("Login error: " + err.Error())
	}
	if !ok {
		c.AuthenticationRequired(w, req)
		return
	}
	user := req.FormValue("username")

	w.Header().Set("User", user)

	c.newAuthzToken(w, req.Request, user)
}

func (c *Context) LogoutHandler(w web.ResponseWriter, req *web.Request) {
	// delete session/cookie
	if len(c.User) == 0 {
		c.AuthenticationRequired(w, req)
		return
	}
}

func (c *Context) DeleteAccountHandler(w web.ResponseWriter, req *web.Request) {
	if len(c.User) == 0 {
		c.AuthenticationRequired(w, req)
		return
	}
	err := c.deleteUser(c.User)
	if err != nil {
		logger.Info().Msg("Error closing account for " + c.User + ":" + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// also clean sessions, etc by logging user out
}

func (c *Context) NewAccountHandler(w web.ResponseWriter, req *web.Request) {
	err := c.addUser(req.FormValue("username"), req.FormValue("password"), req.FormValue("email"))
	if err != nil {
		errMsg := "Error creating account: " + err.Error()
		logger.Info().Msg(errMsg)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errMsg))
		return
	}
	// start new session
	c.User = req.FormValue("username")
	w.Write([]byte("Account created!"))
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

	err := c.Config.BoltDB.View(func(tx *bolt.Tx) error {
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
	err := c.Config.BoltDB.Update(func(tx *bolt.Tx) error {
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

func (c *Context) deleteUser(user string) error {
	err := c.Config.BoltDB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(user))
	})
	return err
}
