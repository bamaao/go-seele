/**
*  @file
*  @copyright defined in go-seele/LICENSE
 */

package light

import (
	"github.com/seeleteam/go-seele/common"
	"github.com/seeleteam/go-seele/database"
	"github.com/seeleteam/go-seele/trie"
)

type odrDatabase struct {
	kvs map[string][]byte
}

func newOdrDatabase() *odrDatabase {
	return &odrDatabase{make(map[string][]byte)}
}

// Get implements the trie.Database interface to store trie node key-value pairs.
func (db *odrDatabase) Get(key []byte) ([]byte, error) {
	return db.kvs[string(key)], nil
}

type odrTrie struct {
	odr      *odrBackend
	root     common.Hash
	db       *odrDatabase
	dbPrefix []byte
	trie     *trie.Trie
}

func newOdrTrie(odr *odrBackend, root common.Hash, dbPrefix []byte) *odrTrie {
	return &odrTrie{
		odr:      odr,
		root:     root,
		db:       newOdrDatabase(),
		dbPrefix: dbPrefix,
	}
}

func (t *odrTrie) Hash() common.Hash {
	panic("unsupported")
}

func (t *odrTrie) Commit(batch database.Batch) common.Hash {
	panic("unsupported")
}

func (t *odrTrie) Get(key []byte) ([]byte, bool) {
	request := &odrTriePoof{
		Root: t.root,
		Key:  key,
	}

	var err error

	// send ODR request to get trie proof.
	if err = t.odr.sendRequest(request); err != nil {
		// @todo refactor the trie struct to return error for Get method.
		// In full node, the Get method should return error for any levelDB error.
		return nil, false
	}

	if err = request.getError(); err != nil {
		return nil, false
	}

	// insert the trie proof in databse.
	for k, v := range request.Proof {
		t.db.kvs[k] = v
	}

	// construct the MPT for the first time.
	if t.trie == nil {
		t.trie, err = trie.NewTrie(t.root, t.dbPrefix, t.db)
		if err != nil {
			return nil, false
		}
	}

	return t.trie.Get(key)
}

func (t *odrTrie) Put(key, value []byte) error {
	panic("unsupported")
}

func (t *odrTrie) DeletePrefix(prefix []byte) bool {
	panic("unsupported")
}

func (t *odrTrie) GetProof(key []byte) (map[string][]byte, error) {
	panic("unsupported")
}
