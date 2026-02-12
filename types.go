package mobile

type MobileNodeOptions struct {
	FullNodeMode             bool
	BootnodeMode             bool
	Bootnodes                string
	StaticNodes              string
	DataDir                  string
	WelcomeMessage           string
	BlockchainRpcEndpoint    string
	SwapInitialDeposit       string
	PaymentThreshold         string
	SwapEnable               bool
	ChequebookEnable         bool
	UsePostageSnapshot       bool
	Mainnet                  bool
	NetworkID                int64
	NATAddr                  string
	CacheCapacity            int64
	DBOpenFilesLimit         int64
	DBWriteBufferSize        int64
	DBBlockCacheCapacity     int64
	DBDisableSeeksCompaction bool
	RetrievalCaching         bool
}

type File struct {
	Name string
	Data []byte
}

type BlockchainData struct {
	WalletAddress     string
	ChequebookAddress string
	ChequebookBalance string
}

type NodeModeType int

const (
	NodeModeUltraLight NodeModeType = iota
	NodeModeLight
	NodeModeFull
)

func (n NodeModeType) String() string {
	switch n {
	case NodeModeUltraLight:
		return "ultra-light"
	case NodeModeLight:
		return "light"
	case NodeModeFull:
		return "full"
	default:
		return "unknown"
	}
}

type StampData struct {
	Label         string
	BatchID       []byte
	BatchAmount   string
	BatchDepth    byte
	BucketDepth   byte
	ImmutableFlag bool
}
