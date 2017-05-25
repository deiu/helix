package helix

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testUser = "alice"
	testPass = "testpass"
)

func Test_Pass_Register(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.registerPass("", "")
	assert.Error(t, err)

	err = ctx.registerPass(testUser, testPass)
	assert.Error(t, err)

	err = ctx.StartBolt()
	assert.NoError(t, err)

	err = ctx.registerPass(testUser, testPass)
	assert.NoError(t, err)

	boltCleanup(ctx)
}

func Test_Pass_Verify(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.StartBolt()
	assert.NoError(t, err)

	ok, err := ctx.verifyPass("", "")
	assert.Error(t, err)
	assert.False(t, ok)

	err = ctx.registerPass(testUser, testPass)
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

func Test_Pass_StoreNoUser(t *testing.T) {
	ctx := NewContext()
	ctx.Config = NewConfig()

	err := ctx.StartBolt()
	assert.NoError(t, err)

	err = ctx.storePass("", "")
	assert.Error(t, err)

	err = ctx.storePass("", testPass)
	assert.Error(t, err)

	boltCleanup(ctx)
}

func boltCleanup(ctx *Context) {
	ctx.BoltDB.Close()
	err := os.Remove(ctx.Config.BoltPath)
	if err != nil {
		panic("Failed to remove " + ctx.Config.BoltPath)
	}
}
