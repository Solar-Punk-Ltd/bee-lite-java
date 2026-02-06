package mobile

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	beelite "github.com/Solar-Punk-Ltd/bee-lite"
	"github.com/ethersphere/bee/v2/pkg/swarm"
)

const StringSliceDelimiter = "|"

type MobileNode interface {
	BlockchainData() (*BlockchainData, error)
	ConnectedPeerCount() int
	Download(hash string) (*File, error)
	Shutdown() error
	WalletAddress() string
}

type MobileNodeImp struct {
	beeClient *beelite.Beelite
	nodeMode  NodeModeType
}

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

func StartNode(options *MobileNodeOptions, password string, verbosity string) (MobileNode, error) {

	beeliteOptions, err := convert(options)

	fmt.Printf("password: %s\n", password)
	fmt.Printf("%+v\n", beeliteOptions)
	if err != nil {
		return nil, err
	}

	beeClient, err := beelite.Start(beeliteOptions, password, verbosity)
	if err != nil {
		return nil, err
	}

	return &MobileNodeImp{beeClient: beeClient, nodeMode: determineNodeMode(options)}, nil
}

func determineNodeMode(options *MobileNodeOptions) NodeModeType {
	if options.FullNodeMode && !options.BootnodeMode {
		return NodeModeFull
	} else if options.BlockchainRpcEndpoint != "" {
		return NodeModeLight
	}
	return NodeModeUltraLight
}

func convert(options *MobileNodeOptions) (*beelite.LiteOptions, error) {
	validateErr := validate(options)
	if validateErr != nil {
		return nil, validateErr
	}

	bootNodes := []string{}
	if options.Bootnodes != "" {
		bootNodes = strings.Split(options.Bootnodes, StringSliceDelimiter)
	}

	staticNodesOpt := []string{}
	if options.StaticNodes != "" {
		staticNodesOpt = strings.Split(options.StaticNodes, StringSliceDelimiter)
	}

	return &beelite.LiteOptions{
		FullNodeMode:             options.FullNodeMode,
		BootnodeMode:             options.BootnodeMode,
		Bootnodes:                bootNodes,
		StaticNodes:              staticNodesOpt,
		DataDir:                  options.DataDir,
		WelcomeMessage:           options.WelcomeMessage,
		BlockchainRpcEndpoint:    options.BlockchainRpcEndpoint,
		SwapInitialDeposit:       options.SwapInitialDeposit,
		PaymentThreshold:         options.PaymentThreshold,
		SwapEnable:               options.SwapEnable,
		ChequebookEnable:         options.ChequebookEnable,
		UsePostageSnapshot:       options.UsePostageSnapshot,
		Mainnet:                  options.Mainnet,
		NetworkID:                uint64(options.NetworkID),
		NATAddr:                  options.NATAddr,
		CacheCapacity:            uint64(options.CacheCapacity),
		DBOpenFilesLimit:         uint64(options.DBOpenFilesLimit),
		DBWriteBufferSize:        uint64(options.DBWriteBufferSize),
		DBBlockCacheCapacity:     uint64(options.DBBlockCacheCapacity),
		DBDisableSeeksCompaction: options.DBDisableSeeksCompaction,
		RetrievalCaching:         options.RetrievalCaching,
	}, nil
}

func validate(options *MobileNodeOptions) error {
	if options.NetworkID < 0 {
		return errors.New("network ID must be a non-negative integer")
	}

	if options.CacheCapacity < 0 {
		return errors.New("cache capacity must be a non-negative integer")
	}

	if options.DBOpenFilesLimit < 0 {
		return errors.New("cache capacity must be a non-negative integer")
	}

	if options.DBWriteBufferSize < 0 {
		return errors.New("DBWriteBufferSize must be a non-negative integer")
	}

	if options.DBOpenFilesLimit < 0 {
		return errors.New("DBOpenFilesLimit must be a non-negative integer")
	}

	return nil
}

func (bl *MobileNodeImp) Download(hash string) (*File, error) {
	bl.beeClient.GetLogger().Info("downloading: ", "hash", hash)

	var result *File = nil
	if hash == "" {
		e := fmt.Errorf("please enter a hash")
		return nil, e
	}
	dlAddr, err := swarm.ParseHexAddress(hash)
	if err != nil {
		return nil, err
	}

	ref, fileName, err := bl.beeClient.GetBzz(context.Background(), dlAddr, nil, nil, nil)
	if err != nil {
		bl.beeClient.GetLogger().Error(err, "download failed")
		return nil, err
	}

	hash = ""
	data, err := io.ReadAll(ref)
	if err != nil {
		bl.beeClient.GetLogger().Error(err, "convert to bytes failed")
		return nil, err
	}

	bl.beeClient.GetLogger().Info("download succeeded", "fileName", fileName, "size", len(data))
	result = &File{Name: fileName, Data: data}

	return result, nil
}

func (m *MobileNodeImp) WalletAddress() string {
	return m.beeClient.OverlayEthAddress().String()
}

func (m *MobileNodeImp) BlockchainData() (*BlockchainData, error) {
	chequebookBalance, err := m.getChequebookBalance()
	chequebookAddress := m.getChequebookAddr()

	if err != nil {
		return nil, err
	}

	return &BlockchainData{
		WalletAddress:     m.beeClient.OverlayEthAddress().String(),
		ChequebookAddress: chequebookAddress,
		ChequebookBalance: chequebookBalance,
	}, nil
}

func (m *MobileNodeImp) ConnectedPeerCount() int {
	return m.beeClient.ConnectedPeerCount()
}

func (m *MobileNodeImp) Shutdown() error {
	err := m.beeClient.Shutdown()
	if err == nil {
		m.beeClient.GetLogger().Info("shutdown succeeded")
		return nil
	}
	m.beeClient.GetLogger().Error(err, "shutdown failed")
	return err
}

func (m *MobileNodeImp) getChequebookAddr() string {

	if m.nodeMode == NodeModeUltraLight {
		return "N/A"
	}

	return m.beeClient.ChequebookAddr().String()
}

func (m *MobileNodeImp) getChequebookBalance() (string, error) {
	if m.nodeMode == NodeModeUltraLight {
		return "N/A", nil
	}

	chequebookBalance, err := m.beeClient.ChequebookBalance()
	if err != nil {
		m.beeClient.GetLogger().Error(err, "failed to get chequebook balance")
		return "", err
	}

	return chequebookBalance.String(), nil

}
