package helix

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUser  = "alice"
	testPass  = "testpass"
	testEmail = "foo@bar.baz"
)

func Test_AddUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.addUser("", "", "")
	assert.Error(t, err)

	err = ctx.addUser(testUser, "", "")
	assert.Error(t, err)

	err = ctx.addUser(testUser, testPass, "")
	assert.Error(t, err)

	err = ctx.addUser("", testPass, "")
	assert.Error(t, err)

	err = ctx.addUser("", "", testEmail)
	assert.Error(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.Error(t, err)

	err = ctx.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	boltCleanup(ctx)
}

func Test_SaveUserFail(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.StartBolt()
	assert.NoError(t, err)

	user := NewUser()

	err = ctx.saveUser(user)
	assert.Error(t, err)

	boltCleanup(ctx)
}

func Test_GetUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.StartBolt()
	assert.NoError(t, err)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	_, err = ctx.getUser("foo")
	assert.Error(t, err)

	user, err := ctx.getUser(testUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser, user.Username)
	assert.Equal(t, testEmail, user.Email)

	boltCleanup(ctx)
}

func Test_SavePassFail(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.savePass(testUser, "")
	assert.Error(t, err)

	err = ctx.StartBolt()
	assert.NoError(t, err)

	err = ctx.savePass("foo", testPass)
	assert.Error(t, err)

	boltCleanup(ctx)
}

func Test_VerifyPass(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.StartBolt()
	assert.NoError(t, err)

	ok, err := ctx.verifyPass("", "")
	assert.Error(t, err)
	assert.False(t, ok)

	err = ctx.addUser(testUser, testPass, testEmail)
	assert.NoError(t, err)

	ok, err = ctx.verifyPass("bob", "foo")
	assert.Error(t, err)
	assert.False(t, ok)

	ok, err = ctx.verifyPass(testUser, "foo")
	assert.Error(t, err)
	assert.False(t, ok)

	ok, err = ctx.verifyPass(testUser, testPass)
	assert.NoError(t, err)
	assert.True(t, ok)

	boltCleanup(ctx)
}

func boltCleanup(ctx *Context) {
	ctx.BoltDB.Close()
	err := os.Remove(ctx.Config.BoltPath)
	if err != nil {
		panic("Failed to remove " + ctx.Config.BoltPath)
	}
}
