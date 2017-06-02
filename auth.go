package helix

import (
	"errors"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/deiu/webid-rsa"
	"github.com/gocraft/web"
	"gopkg.in/hlandau/passlib.v1"
)

func (c *Context) AuthenticationRequired(w web.ResponseWriter, req *web.Request) {
	if len(reqUser(req)) == 0 {
		authn := webidrsa.NewAuthenticateHeader(req.Request)
		w.Header().Set("WWW-Authenticate", authn)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func (c *Context) Authentication(w web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	errMsg := ""
	user := ""

	if len(req.Header.Get("Authorization")) > 0 {
		authz, err := webidrsa.ParseAuthorizationHeader(req.Header.Get("Authorization"))
		if err != nil {
			errMsg = err.Error()
		}
		switch authz.Type {
		case "WebID-RSA":
			user, err = webidrsa.Authenticate(req.Request)
		case "Bearer":
			token, _ := ParseBearerAuthorizationHeader(req.Header.Get("Authorization"))
			user, err = c.getAuthzUserFromToken(token, req.Host, webidrsa.GetOrigin(req.Request))
		}

		if err != nil {
			errMsg = "Could not authenticate user: " + err.Error()
		}
	}
	req.Header.Set("User", user)
	logger = zerolog.New(os.Stderr).With().
		Timestamp().
		Str("Method", req.Method).
		Str("Path", req.Request.URL.String()).
		Str("User", user).
		Logger()

	logger.Info().Msg(errMsg)

	next(w, req)
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

func ParseBearerAuthorizationHeader(header string) (string, error) {
	if len(header) == 0 {
		return "", errors.New("Cannot parse Authorization header: no header present")
	}

	parts := strings.SplitN(header, " ", 2)
	if parts[0] != "Bearer" {
		return "", errors.New("Not a Bearer header. Got: " + parts[0])
	}
	return url.QueryUnescape(parts[1])
}
