package aiload

import (
	"hash/maphash"
	"strconv"
)

var hashSeed = maphash.MakeSeed()

func (p *MqTaskId) Unique() uint64 {
	return maphash.String(hashSeed, strconv.FormatInt(int64(p.Id), 10))
}
