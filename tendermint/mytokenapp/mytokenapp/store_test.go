package mytokenapp

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStore(t *testing.T) {
	store := NewStore()
	err := store.SetBalance([]byte("address1"), 100)
	require.NoError(t, err, "set error : %v\n", err)

	err = store.SetBalance([]byte("address2"), 2000)
	require.NoError(t, err, "set error: %v\n", err)

	balance1, err := store.GetBalance([]byte("address1"))
	balance2, err := store.GetBalance([]byte("address2"))
	require.Equal(t, int64(100), balance1, "balance1 not matched")
	require.Equal(t, int64(2000), balance2, "balance2 not matched")

	store.Commit()

	err = store.SetBalance([]byte("address1"), 300)
	require.NoError(t, err, "set error : %v\n", err)

	err = store.SetBalance([]byte("address2"), 100)
	require.NoError(t, err, "set error: %v\n", err)

	store.Commit()

	b1, err := store.GetBalanceVersioned([]byte("address1"), store.LastVersion)
	b2, err := store.GetBalanceVersioned([]byte("address2"), store.LastVersion)
	require.Equal(t, int64(300), b1, "b1 not matched")
	require.Equal(t, int64(100), b2, "b2 not matched")

}
