package memory

type MemoryInfo struct {
	Total          uint64
	Free           uint64
	Available      uint64
	Buffers        uint64
	Cached         uint64
	CachedActive   uint64
	CachedInactive uint64
	SwapTotal      uint64
	SwapFree       uint64
}

type IDs struct {
	Real, Effective, Saved, FS uint
}

type Process struct {
	Name  string
	User  string
	Group string
	Pid   uint64
	VmRSS uint64
	State string
	Uid   IDs
	Gid   IDs
}
