package utxo

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	. "gopkg.in/check.v1"
)

type BitcoinTemporalStorageTestSuite struct{}

var _ = Suite(
	&BitcoinTemporalStorageTestSuite{},
)

func (s *BitcoinTemporalStorageTestSuite) TestNewTemporalStorage(c *C) {
	memStorage := storage.NewMemStorage()
	db, err := leveldb.Open(memStorage, nil)
	c.Assert(err, IsNil)
	dbTemporalStorage, err := NewTemporalStorage(db, 0)
	c.Assert(err, IsNil)
	c.Assert(dbTemporalStorage, NotNil)
	c.Assert(db.Close(), IsNil)
}
