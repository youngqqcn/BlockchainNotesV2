package mytokenapp

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto"
	db "github.com/tendermint/tm-db"
)

type Store struct {
	tree        *iavl.MutableTree
	LastVersion int64
	LastHash    []byte
}

func NewStore(dirPath string) *Store {

	ldb, err := db.NewGoLevelDB("account", dirPath)
	if err != nil {
		panic(err)
	}

	tree, err := iavl.NewMutableTree(ldb, 1024)
	if tree == nil || err != nil {
		panic(err)
	}

	ver, err := tree.Load()
	if err != nil {
		panic(err)
	}

	hash := tree.Hash()
	return &Store{
		tree:        tree,
		LastVersion: ver,
		LastHash:    hash,
	}
}

func (store *Store) SetBalance(addr crypto.Address, balance int64) error {

	//key, err := codec.MarshalBinaryBare(addr)
	//if err != nil {
	//	return err
	//}

	value, err := codec.MarshalBinaryBare(balance)
	if err != nil {
		return err
	}

	store.tree.Set(addr, value)
	return nil
}

func (store *Store) GetBalance(addr crypto.Address) (balance int64, err error) {
	_, value := store.tree.Get(addr)

	err = codec.UnmarshalBinaryBare(value, &balance)
	if err != nil {
		return
	}
	return
}

func (store *Store) GetBalanceVersioned(addr crypto.Address, version int64) (balance int64, err error) {
	_, bz := store.tree.GetVersioned(addr, version)
	err = codec.UnmarshalBinaryBare(bz, &balance)
	if err != nil {
		return
	}
	return
}

func (store *Store) Commit() {
	hash, ver, err := store.tree.SaveVersion()
	if err != nil {
		panic(err)
	}
	store.LastHash = hash
	store.LastVersion = ver

}
