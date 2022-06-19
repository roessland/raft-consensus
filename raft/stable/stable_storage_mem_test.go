package stable_test

import (
	"github.com/roessland/raft-consensus/raft/stable"
)

// Verify that fakes satisfy interfaces.

var _ stable.IntStore = stable.NewInMemoryIntStore()
var _ stable.NullableIntStore = stable.NewInMemoryNullableIntStore()
var _ stable.LogEntriesStore = stable.NewInMemoryLogEntriesStore()
