package memory

type MemoryInfo struct {
	Total uint64
	Free uint64
	Available uint64
	Buffers uint64
	Cached uint64
	SwapTotal uint64
	SwapFree uint64
}
