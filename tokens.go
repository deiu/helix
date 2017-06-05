package helix

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/deiu/webid-rsa"
	"github.com/gocraft/web"
)

var (
	tokenDuration = time.Hour * 5040
)

func tokenDateIsValid(valid string) error {
	v, err := strconv.ParseInt(valid, 10, 64)
	if err != nil {
		return err
	}
	if time.Now().Local().UnixNano() > v {
		return errors.New("Token has expired!")
	}

	return nil
}

func (ctx *Context) getAuthzUserFromToken(token, host, origin string) (string, error) {
	values, err := ctx.getPersistedToken("Authorization", host, token)
	if err != nil {
		return "", err
	}
	if len(values["webid"]) == 0 || len(values["valid"]) == 0 || len(values["origin"]) == 0 {
		return "", errors.New("Token is missing required values")
	}
	err = tokenDateIsValid(values["valid"])
	if err != nil {
		return "", err
	}
	if origin != values["origin"] {
		return "", errors.New("Cannot authorize user: " + values["webid"] + ". Origin: " + origin + " does not match the origin in the token: " + values["origin"])
	}
	return values["webid"], nil
}

func (c *Context) newAuthzToken(w web.ResponseWriter, req *http.Request, user string) {
	values := map[string]string{
		"webid":  user,
		"origin": webidrsa.GetOrigin(req),
	}
	token, err := c.newPersistedToken("Authorization", req.Host, values)
	if err != nil {
		logger.Info().Msg(err.Error())
		return
	}
	w.Header().Set("Token", token)
}

// newPersistedToken saves an API token to the bolt db. It returns the API token and a possible error
func (ctx *Context) newPersistedToken(tokenType, host string, values map[string]string) (string, error) {
	var token string
	if len(tokenType) == 0 {
		return token, errors.New("Missing token type when trying to generate new token")
	}
	// bucket(host) -> bucket(type) -> values
	err := ctx.Config.BoltDB.Update(func(tx *bolt.Tx) error {
		userBucket, err := tx.CreateBucketIfNotExists([]byte(host))
		if err != nil {
			return err
		}
		bucket, err := userBucket.CreateBucketIfNotExists([]byte(tokenType))
		id, _ := bucket.NextSequence()
		values["id"] = fmt.Sprintf("%d", id)
		// set validity if not alreay set
		if len(values["valid"]) == 0 {
			// age times the duration of 6 month
			values["valid"] = fmt.Sprintf("%d",
				time.Now().Add(time.Duration(ctx.Config.TokenAge)*tokenDuration).UnixNano())
		}
		// marshal values to JSON; this will never error since we only marshal strings
		tokenJson, _ := json.Marshal(values)
		token = fmt.Sprintf("%x", sha256.Sum256(tokenJson))
		bucket.Put([]byte(token), tokenJson)

		return nil
	})

	return token, err
}

func (ctx *Context) getPersistedToken(tokenType, host, token string) (map[string]string, error) {
	tokenValues := map[string]string{}
	if len(tokenType) == 0 || len(host) == 0 || len(token) == 0 {
		return tokenValues, errors.New("Can't retrieve token from db. tokenType, host and token value are requrired.")
	}
	err := ctx.Config.BoltDB.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(host))
		if userBucket == nil {
			return errors.New(host + " bucket not found!")
		}
		bucket := userBucket.Bucket([]byte(tokenType))
		if bucket == nil {
			return errors.New(tokenType + " bucket not found!")
		}

		// unmarshal
		b := bucket.Get([]byte(token))
		err := json.Unmarshal(b, &tokenValues)
		return err
	})
	return tokenValues, err
}

func (ctx *Context) getTokenByOrigin(tokenType, host, origin string) (string, error) {
	token := ""
	if len(tokenType) == 0 || len(host) == 0 || len(origin) == 0 {
		return token, errors.New("Can't retrieve token from db. tokenType, host and token value are requrired.")
	}
	err := ctx.Config.BoltDB.View(func(tx *bolt.Tx) error {
		userBucket := tx.Bucket([]byte(host))
		if userBucket == nil {
			return errors.New(host + " bucket not found!")
		}
		bucket := userBucket.Bucket([]byte(tokenType))
		if bucket == nil {
			return errors.New(tokenType + " bucket not found!")
		}

		// unmarshal
		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := string(k)
			values, err := ctx.getPersistedToken(tokenType, host, key)
			if err == nil && values["origin"] == origin {
				token = key
				break
			}
		}

		return nil
	})
	return token, err
}

func (ctx *Context) deletePersistedToken(tokenType, host, token string) error {
	if len(tokenType) == 0 || len(host) == 0 || len(token) == 0 {
		return errors.New("Can't retrieve token from db. tokenType, host and token value are requrired.")
	}
	err := ctx.Config.BoltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(host))
		if b == nil {
			return errors.New("No bucket for host " + host)
		}
		bucket := b.Bucket([]byte(tokenType))
		if bucket == nil {
			return errors.New("No bucket for token type " + tokenType)
		}

		return bucket.Delete([]byte(token))
	})
	return err
}

func (ctx *Context) getTokensByType(tokenType, host string) (map[string]map[string]string, error) {
	tokens := make(map[string]map[string]string)
	err := ctx.Config.BoltDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(host))
		if b == nil {
			return errors.New("No bucket for host " + host)
		}
		ba := b.Bucket([]byte(tokenType))
		if ba == nil {
			return errors.New("No bucket for type " + tokenType)
		}

		c := ba.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			key := string(k)
			token, err := ctx.getPersistedToken(tokenType, host, key)
			if err == nil {
				tokens[key] = token
			}
		}
		return nil
	})
	return tokens, err
}
