package types

const (
	// ModuleName defines the module name
	ModuleName = "mars"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_mars"

	// Version defines the current version the IBC module supports
	Version = "mars-1"

	// PortID is the default port id that module binds to
	PortID = "mars"
)

var (
	ParamsKey = []byte("p_mars")
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = KeyPrefix("mars-port-")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
