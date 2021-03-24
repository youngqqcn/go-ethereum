
代码目录结构

``
$ tree . -d -L 1
.
├── accounts  账户相关
├── build 临时构建目录
├── cmd   主要的进程(命令行)的入口
├── common  公共
├── consensus  共识算法相关
├── console   控制台
├── contracts   合约
├── core   核心
├── crypto  加密相关
├── docs  文档
├── eth    以太坊
├── ethclient   以太坊客户端
├── ethdb   数据库
├── ethstats   状态
├── event   事件
├── graphql   GraphQL支持
├── internal   内部模块
├── les       轻以太坊子协议(Light Ethereum Subprotocol).
├── light    轻节点
├── log    日志
├── metrics   指标参数相关
├── miner  挖矿
├── mobile   移动端支持
├── node   节点
├── p2p    P2P网络相关
├── params   参数
├── rlp      RLP编码
├── rpc      RPC接口
├── signer    签名
├── swarm     Swarm支持
├── tests     测试用例
├── trie      Trie前缀树实现
└── yqq      我自己的临时笔记目录

```









```go
// Node is a container on which services can be registered.
type Node struct {
	eventmux      *event.TypeMux
	config        *Config
	accman        *accounts.Manager
	log           log.Logger
	ephemKeystore string            // if non-empty, the key directory that will be removed by Stop
	dirLock       fileutil.Releaser // prevents concurrent use of instance directory
	stop          chan struct{}     // Channel to wait for termination notifications
	server        *p2p.Server       // Currently running P2P networking layer
	startStopLock sync.Mutex        // Start/Stop are protected by an additional lock
	state         int               // Tracks state of node lifecycle

	lock          sync.Mutex
	
    
    // 此lifecycles 包含了所有的服务(实现了Start  Stop方法)
    lifecycles    []Lifecycle // All registered backends, services, and auxiliary services that have a lifecycle
	
    
    rpcAPIs       []rpc.API   // List of APIs currently provided by the node
	http          *httpServer //
	ws            *httpServer //
	ipc           *ipcServer  // Stores information about the ipc http server
	inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	databases map[*closeTrackingDB]struct{} // All open databases
}
```


```go

package node

// Lifecycle encompasses the behavior of services that can be started and stopped
// on the node. Lifecycle management is delegated to the node, but it is the
// responsibility of the service-specific package to configure and register the
// service on the node using the `RegisterLifecycle` method.
type Lifecycle interface {
	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the service.
	Start() error

	// Stop terminates all goroutines belonging to the service, blocking until they
	// are all terminated.
	Stop() error
}

```


以下结构体实现了 Lifecycle

```go
package eth

// Ethereum implements the Ethereum full node service.
type Ethereum struct {
	config *Config

	// Handlers
	txPool             *core.TxPool
	blockchain         *core.BlockChain
	handler            *handler
	ethDialCandidates  enode.Iterator
	snapDialCandidates enode.Iterator

	// DB interfaces
	chainDb ethdb.Database // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *core.ChainIndexer             // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *EthAPIBackend

	miner     *miner.Miner
	gasPrice  *big.Int
	etherbase common.Address

	networkID     uint64
	netRPCService *ethapi.PublicNetAPI

	p2pServer *p2p.Server

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price and etherbase)
}
```




```go
package ethstats

// Service implements an Ethereum netstats reporting daemon that pushes local
// chain statistics up to a monitoring server.
type Service struct {
	server  *p2p.Server // Peer-to-peer server to retrieve networking infos
	backend backend
	engine  consensus.Engine // Consensus engine to retrieve variadic block fields

	node string // Name of the node to display on the monitoring page
	pass string // Password to authorize access to the monitoring page
	host string // Remote address of the monitoring service

	pongCh chan struct{} // Pong notifications are fed into this channel
	histCh chan []uint64 // History request block numbers are fed into this channel

}
```




```go

// Package les implements the Light Ethereum Subprotocol.
package les

type LightEthereum struct {
	lesCommons

	peers          *serverPeerSet
	reqDist        *requestDistributor
	retriever      *retrieveManager
	odr            *LesOdr
	relay          *lesTxRelay
	handler        *clientHandler
	txPool         *light.TxPool
	blockchain     *light.LightChain
	serverPool     *serverPool
	valueTracker   *lpc.ValueTracker
	dialCandidates enode.Iterator
	pruner         *pruner

	bloomRequests chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer  *core.ChainIndexer             // Bloom indexer operating during block imports

	ApiBackend     *LesApiBackend
	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager
	netRPCService  *ethapi.PublicNetAPI

	p2pServer *p2p.Server
}
```



```go
package les

type LesServer struct {
	lesCommons

	ns          *nodestate.NodeStateMachine
	archiveMode bool // Flag whether the ethereum node runs in archive mode.
	handler     *serverHandler
	broadcaster *broadcaster
	lesTopics   []discv5.Topic
	privateKey  *ecdsa.PrivateKey

	// Flow control and capacity management
	fcManager    *flowcontrol.ClientManager
	costTracker  *costTracker
	defParams    flowcontrol.ServerParams
	servingQueue *servingQueue
	clientPool   *clientPool

	minCapacity, maxCapacity uint64
	threadsIdle              int // Request serving threads count when system is idle.
	threadsBusy              int // Request serving threads count when system is busy(block insertion).

	p2pSrv *p2p.Server
}
```


