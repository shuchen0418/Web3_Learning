# Web3

## 问题

### 1.交易被⽤⼾使⽤eth_sendRawTransaction接⼝发送给执⾏层客⼾端，交易会如何被保存？

在以太坊的执行层客户端（例如 Go Ethereum，Geth）中，当用户通过 `eth_sendRawTransaction` 接口发送交易时，交易首先被发送到交易池（transaction pool）进行管理。Geth 中负责处理这些逻辑的核心代码位于 `core/tx_pool.go` 文件中。

具体步骤如下：

1. **`eth_sendRawTransaction` API 调用**：
   用户通过 `eth_sendRawTransaction` 发送签名好的交易时，客户端会通过 JSON-RPC 接口接收这个交易。可以在 `eth/api.go` 文件中找到这个 API 的处理代码。

   ```go
   // eth/api.go
   func (s *PublicTransactionPoolAPI) SendRawTransaction(ctx context.Context, encodedTx hexutil.Bytes) (common.Hash, error) {
       // 解码交易
       tx := new(types.Transaction)
       if err := rlp.DecodeBytes(encodedTx, tx); err != nil {
           return common.Hash{}, err
       }
       // 将交易发送到交易池
       return SubmitTransaction(ctx, s.b, tx)
   }
   ```

2. **提交交易到交易池**：
   在 `SubmitTransaction` 函数中，交易被提交给交易池进行管理。交易 池会验证交易的有效性并将其保存。

   ```go
   func SubmitTransaction(ctx context.Context, b Backend, tx *types.Transaction) (common.Hash, error) {
       // 提交交易到交易池
       if err := b.TxPool().AddLocal(tx); err != nil {
           return common.Hash{}, err
       }
       return tx.Hash(), nil
   }
   ```

3. **交易池（Transaction Pool）管理**：
   交易被发送到交易池的 `AddLocal` 函数中，该函数会进行一系列检查，如交易的有效性、nonce、gas 价格等。如果交易有效，会将其保存到交易池中。

   主要代码位于 `core/tx_pool.go` 文件中：

   ```go
   // tx_pool.go
   func (pool *TxPool) AddLocal(tx *types.Transaction) error {
       // 检查交易有效性
       if err := pool.validateTx(tx, local); err != nil {
           return err
       }
       // 将交易添加到池中
       pool.addTx(tx)
       return nil
   }
   ```

4. **交易存储**：
   交易池在 `addTx` 方法中将交易按照不同的优先级和 nonce 存储在池中，以便稍后矿工或验证者从池中挑选并打包交易。

   ```go
   // tx_pool.go
   func (pool *TxPool) addTx(tx *types.Transaction) {
       // 根据 sender 和 nonce 组织交易池
       from, _ := types.Sender(pool.signer, tx)
       pool.pending[from] = append(pool.pending[from], tx)
   }
   ```

这些是 Go Ethereum 中 `eth_sendRawTransaction` 接口到交易池处理和存储交易的主要代码路径。

### 2. 交易被保存后，如何被选中？

在执行层客户端中，矿工或验证者从交易池中选择优先级较高的交易进行打包。Geth 中交易选择的逻辑由交易池根据 **gas 价格** 和 **nonce 顺序** 来决定。

对应代码位于 `core/tx_pool.go` 文件中的 `Pending` 方法，它返回可以立即打包进区块的交易。

```go
// tx_pool.go
func (pool *TxPool) Pending() (map[common.Address]types.Transactions, error) {
    // 获取所有 pending 状态的交易
    pool.mu.Lock()
    defer pool.mu.Unlock()
    return pool.pending, nil
}
```

矿工在创建区块时会调用 `core/state_processor.go` 文件中的 `ApplyTransaction` 方法，来将选定的交易应用到当前状态。

```go
// state_processor.go
func ApplyTransaction(...) error {
    // 执行交易并更新状态
    ...
    return nil
}
```

### 3. 一个交易消耗的 gas 如何计算？是什么时候从 sender 账户扣除？

以太坊中交易的 gas 计算在执行交易时逐步累加。每个操作（如计算、存储写入等）都有一个固定的 gas 消耗量，详见 `EVM` 虚拟机的相关代码。Gas 是在交易执行过程中逐步扣除的。

主要代码在 `core/vm/evm.go` 文件中：

```go
// evm.go
func (evm *EVM) Call(...) ([]byte, error) {
    // 根据操作指令扣除 gas
    for step := range steps {
        evm.StateDB.SubBalance(sender, gasCost)
    }
    ...
}
```

Gas 在执行结束后，将总消耗的 gas 从 sender 账户扣除。

### 4. 创建合约交易和普通交易处理方式的区别？

创建合约交易和普通交易在处理方式上不同，创建合约的交易 `to` 字段为空，`data` 字段包含合约字节码。执行时，会在区块链上分配一个新的地址，并将合约字节码存储在该地址。

主要代码在 `core/state_processor.go` 中的 `ApplyTransaction` 中处理普通交易和合约创建交易的区别：

```go
// state_processor.go
func ApplyTransaction(...) error {
    if tx.To() == nil {
        // 创建合约
        contractAddr := crypto.CreateAddress(sender, nonce)
        evm.Create(...)
    } else {
        // 普通交易
        evm.Call(...)
    }
    return nil
}
```

### 5. 创建合约有什么限制？

以太坊对合约创建有以下限制：

1. **Gas 上限**：创建合约需要足够的 gas 来支付字节码存储和执行费用。
2. **合约大小限制**：合约的字节码大小不能超过 24 KB（EIP-170）。

相关检查在 `core/vm/evm.go` 的 `Create` 方法中进行：

```go
// evm.go
func (evm *EVM) Create(...) ([]byte, common.Address, error) {
    // 检查合约大小是否超限
    if len(code) > MaxCodeSize {
        return nil, common.Address{}, ErrMaxCodeSizeExceeded
    }
    ...
}
```

### 6. 交易是怎么被广播给其他执行层客户端的？

交易在接收到后会通过 P2P 网络进行广播。Geth 使用的是一个基于 Kademlia 的 P2P 网络，节点之间会互相同步交易。相关代码位于 `eth/handler.go` 中：

```go
// handler.go
func (pm *ProtocolManager) BroadcastTx(tx *types.Transaction) {
    // 向其他节点广播交易
    pm.peers.BroadcastTx(tx)
}
```

当交易被验证有效后，它会被广播给与该节点连接的其他节点，直到整个网络中所有节点都接收到该交易。

### 7. PoW（Proof of Work）与 PoS（Proof of Stake）？

- **PoW**：在工作量证明（PoW）中，矿工通过计算哈希值来竞争找到一个有效的 nonce。核心代码逻辑位于 `consensus/ethash` 包中。

```go
// ethash.go
func (ethash *Ethash) Seal(...) {
    // 执行哈希计算，找到有效的 nonce
    ...
}
```

- **PoS**：权益证明（PoS）机制下，验证者根据质押的 ETH 数量来获得出块权。相关逻辑位于 Prysm 客户端中，Geth 负责执行层的交易处理。

### 8. PoW 时期，矿工如何赚取收益？收益什么时候被添加到账户？

在 PoW 中，矿工通过成功挖出一个区块来获得区块奖励和交易手续费。收益会在区块被打包并成功广播到网络后添加到账户中。相关逻辑在 `core/blockchain.go` 中：

```go
// blockchain.go
func (bc *BlockChain) FinalizeBlock(...) error {
    // 给矿工账户添加奖励
    bc.State.AddBalance(coinbase, reward)
    ...
}
```

### 9. PoW 时期，如何产生一个新块？并把区块广播出去？

在 PoW 中，矿工通过不断尝试不同的 nonce 来寻找一个符合目标难度的区块哈希值。找到有效哈希值后，会生成一个新区块并通过 P2P 网络进行广播。

主要代码在 `consensus/ethash/worker.go` 中：

```go
// worker.go
func (w *worker) commitNewWork() {
    // 生成新块并广播
    block := w.makeCurrentBlock()
    w.broadcastNewBlock(block)
}
```

### 10. 合并到 PoS 后，新区块如何同步？

在 PoS 中，共识层验证者提议并验证新区块。新区块会通过共识层客户端（如 Prysm）进行同步和验证。Geth 作为执行层客户端，处理交易执行和状态更新。共识层通过 `Engine API` 与执行层进行通信，相关代码在 `consensus/engine.go` 中。

```go
// engine.go
func (eng *ConsensusEngine) NewBlock(...) {
    // 通过共识层同步新区块
    eng.executeBlock()
}
```

## 知识点

以太坊分为：执行层（Execution Layer）和共识层（Consensus Layer） 

​					执行层客户端为go-ethereum（go语言）   共识层客户端为Prysm（go语言）

L1：主流的链（BTC，ETH，SUI，SOL） 

L2：各种业务都要提交到L1上 ，独立不完整，主要是ETH生态居多（optimism，Arbtrum，Base，Linea，PolygonZk，ZkStark）

### P2P网络相关知识：

Bucket最多存储维护的节点：17*24（理论上） go-ethereum还有个启动参数max-peer参数，限制最大维护的peer的数量

想要连入以太坊，需要找可以通信的节点（可以自己指定），然后发送（find_node）拉取引导节点的Bucket的部分节点，引导节点返回消息，用（DiscV5 DiscV4解析 消息优先给v5处理 v5处理不了就给v4），后添加到replacements里，等entries有变化（Ping不通或者消息回复慢），再从replacements里面拿一个到entries里，作为常驻节点。每隔一段时间选取三个节点去执行（find_node）维护数据。

执行层通过RLPx协议广播交易 用udp发送消息

共识层通过gossip协议 建立tcp连接广播消息

DiscV5 DiscV4 是特有的P2P实现，在go-etherrum中仅用于节点发现

gossip是 标准的p2p实现

### 默克尔树与MPT树：

默克尔树：叶子节点是数据块的hash值，非叶子节点存储的是其左右子节点hash值的组合hash，树的顶点是根节点（Merkle Root）

MPT树：结合默克尔树、前缀树（Trie）和Patricia树，如果两个叶子节点有相同的key前缀，就共享这部分的路径节点，每个节点的hash值与默克尔树相同

插入操作：当有一个值拆入，判断new key与已存在的key有无相同的前缀，会有三种种情况：1.相同的部分new key小于 key：会在新的相同的部分产生一个FullNode，然后在这个FullNode上插入新的ShortNode和之前的Node；2.相同的部分new key等于key：替换之前key的value；3.没有相同的部分：直接产生一个FullNode把当前key所在的FullNode插入新的FullNode，new Key成为一个ShortKey

查找：通过相同前缀匹配路径

删除：通过相同前缀匹配路径，找到后删除，修改MPT树的结构，如果删除后FullNode只剩下一个叶子节点，那么这个FullNode节点变成一个Shot节点

以太坊2.0有四个root：StateRoot TransactionsRoot ReceiptRoot WithdrawRoot也就是四个Tree

### 以太坊L2原理：

Layer1平台被称为链上（以太坊）

Layer2指的是基于L1的链下网络、系统或技术，是解决拓展性和吞吐性的主流方案之一

侧链：与L1区块链平行，有**自己的区块生成规则和共识机制**，可以和主链进行跨链操作 ，可以和主网有不同的共识机制比如主网用POW 侧用POS，交易也可以用**不同的币**

RollUp：将L2的一组交易打包成一个交易存储在L1上，减少L1的交易数量和计算负担

主流的两种RollUp方案：Optimistic RollUp 和 Zero-knowledge rollup(ZK rollup)

Optimistic RollUp：乐观假设所有交易都是有效的，可以在没有初始证明的情况下提交批次到主网上，仍何人都可以在挑战期内，检测和证明数据的State Root

(Optimism Arbitrum Base)

Zero-knowledge rollup(ZK rollup)：除了交易数据，post state root和previous state root之外，还要携带一个"有效性证明"，任何人可以用这个证明验证交易的正确性。省略了人的验证工作，提交同时完成验证

(Manta Network、Polygon zkEVM、Linea、Starknet、Scroll、zkSync )

### Optimism的CL客户端

Sequencer：中心化节点，用于对交易排序

Verifer：用于验证Sequencer排序后的tx，验证Sequencer发布的新区块是否正确

### Go语法

#### 数字

整型:int8 int16 int32 int64 uint8 uint16 uint32 uint64 uintptr

浮点数:float32 float 64 复数complex64 complex128 real()获取实部 imag()获取虚部

byte类型: byte==uint8 字符串可以直接被转换成[]byte 相反也是一样 rune类型 rune可以当作uint32 但是不一定相等 var v1 rune = '号'

字符串:都是UTF-8 可以用双引号和反引号 string byte rune可以相互转换，但是string的长度和byte一样 rune的长度是根据符号数量

bool: true false

零值:数字类型0 字符串类型:空字符串 布尔类型:false

#### 数组

申明一个数组

```
var b = [5]int{1, 2, 3, 4, 5}

var a[5]  = [5] int {1,2,3,4,5}

a := [5] int {1,2,3,4,5}

a := [...]int {1,2,3,4,5}

```

修改就是a[索引] = ?  方法中作为参数传入时，修改只在方法内生效，要想全局改变原来的数组值就要传入指针

### 性能分析

#### benchmark 基准测试

 总结：

1. 进行性能测试时，尽可能保持测试环境的稳定
2. 实现 benchmark 测试
   • 位于 `_test.go` 文件中
   • 函数名以 `Benchmark` 开头
   • 参数为 `b *testing.B`
   • `b.ResetTimer()` 可重置定时器
   • `b.StopTimer()` 暂停计时
   • `b.StartTimer()` 开始计时
3. 执行 benchmark 测试
   • `go test -bench .` 执行当前测试
   • `b.N` 决定用例需要执行的次数
   • `-bench` 可传入正则，匹配用例
   • `-cpu` 可改变 CPU 核数
   • `-benchtime` 可指定执行时间或具体次数
   • `-count` 可设置 benchmark 轮数
   • `-benchmem` 可查看内存分配量和分配次数

#### pprof 性能分析

总结：

1. 性能分析类型

   - CPU 性能分析，runtime 每隔 10 ms 中断一次，记录此时正在运行的 goroutines 的堆栈信息
   - 内存性能分析，记录堆内存分配时的堆栈信息，忽略栈内存分配信息，默认每 1000 次采样 1 次
   - 阻塞性能分析，GO 中独有的，记录一个协程等待一个共享资源花费的时间
   - 锁性能分析，记录因为锁竞争导致的等待或延时

2. CPU 性能分析

   - 使用原生 `runtime/pprof` 包，通过在 main 函数中添加代码运行可生成性能分析报告：

     ```
     pprof.StartCPUProfile(os.Stdout)
     defer pprof.StopCPUProfile()
     ```

   - 可通过 `go tool pprof -http=:9999 cpu.pprof` 在 web 页面查看分析数据

   - 可通过 `go tool pprof cpu.prof` 交互模式查看分析数据，可使用 `help` 查看支持的命令和选项

3. 内存性能分析

   - 使用 `pkg/profile` 库，通过在 main 函数中添加代码运行可生成性能分析报告：

     ```
     defer profile.Start(profile.MemProfile, profile.MemProfileRate(1)).Stop()
     ```

   - 同样可通过 web 页面或交互模式查看分析数据

4. benchmark 生成 profile

- 可通过在 `go test ` 中添加参数 `-cpuprofile=$FILE,-memprofile=$FILE,-blockprofile=$FILE` 生成相应的 profile 文件
- 生成的 profile 文件同样可通过 web 页面或交互模式查看分析数据

### 常用数据结构

#### 字符串拼接性能及原理

总结：

1. 字符串最高效的拼接方式是结合预分配内存方式 `Grow` 使用 `string.Builder`
2. 当使用 `+` 拼接字符串时，生成新字符串，需要开辟新的空间
3. 当使用 `strings.Builder`，`bytes.Buffer` 或 `[]byte` 的内存是按倍数申请的，在原基础上不断增加
4. `strings.Builder` 比 `bytes.Buffer` 性能更快，一个重要区别在于 `bytes.Buffer` 转化为字符串重新申请了一块空间存放生成的字符串变量；而 `strings.Builder` 直接将底层的 `[]byte` 转换成字符串类型返回

#### 切片(slice)性能及陷阱

总结：

1. GO 中的数组变量属于值类型，当数组变量被赋值或传递时，实际上会复制整个数组
2. 切片本质是数组片段的描述，包括数组的指针，片段的长度和容量，切片操作并不复制切片指向的元素，而是复用原来切片的底层数组
   - 长度是切片实际拥有的元素，使用 `len` 可得到切片长度
   - 容量是切片预分配的内存能够容纳的元素个数，使用cap 可得到切片容量
     - 当 append 之后的元素小于等于 cap，将会直接利用底层元素剩余的空间
     - 当 append 后的元素大于 cap，将会分配一块更大的区域来容纳新的底层数组，在容量较小的时候，通常是以 2 的倍数扩大
3. 可能存在只使用了一小段切片，但是底层数组仍被占用，得不到使用，推荐使用 `copy` 替代默认的 `re-slice`

#### for 和 range 的性能比较

总结：

只有当range遍历值(值不是指针且较大)的情况下，range的性能比for差，其他时候都差不多

#### Go Reflect 提高反射性能

总结：

reflect在取结构体的每个字段时，尽量用Field()而不是FieldByName()，因为是顺序存储的，用下标更快

#### Go 空结构体 struct{} 的使用

总结：

1.实现集合（Set）存储的是key => 空结构体 的结构

2.不发送数据的信道 发送空结构体到信道，仅仅作为通知，不消耗内存

3.仅包含方法的结构体 有些结构体只用于方法，不需要任何字段 这时候可以声明这个结构体为空结构体

#### Go struct 内存对齐

总结：

1.一个结构体实例所占据的空间等于各字段占据空间之和，再加上内存对齐的空间大小

2.合理的内存对齐可以提高内存读写的性能，并且便于实现变量操作的原子性

3.空结构体和空数组的内存为0，在对内存特别敏感的结构体的设计上，我们可以通过调整字段的顺序，减少内存的占用

4.结构体中的空结构体不要放在最后

### 并发编程

#### 读写锁和互斥锁的性能比较

总结：

1.互斥锁:就是两个代码片段互相排斥，只有当一个执行完后另一个才能执行

2.读写锁:读锁之间不会互斥，没有写所锁，可以多个协程同时获得读锁；写锁之间互斥；读锁与写锁之间互斥

3.当读的情况远大于写的情况下，读写锁的性能比互斥锁高

#### 如何退出协程 goroutine (超时场景)

总结：

当一个函数传入的参数时chan时：

1.传入的chan应该是有缓冲区的，以免超时goroutine阻塞，不能退出造成持续内存占用

2.或者设计函数时用select语句，goroutine执行这个函数时，执行不成功直接default：return，以免超时goroutine阻塞，不能退出造成持续内存占用

3.设计函数多段执行，执行一段代码就用select语句，执行不成功就直接default：return，以免超时goroutine阻塞，不能退出造成持续内存占用

#### 如何退出协程 goroutine (其他场景)

总结:

1.暴力退出：直接close(chan)，如果panic就defer中recover

2.礼貌退出：写一个结构体包含chan 和 sync.Once，确保这个chan只会被关闭一次，以免panic

3.优雅退出：1发送者 n接收者的情况，1发送者去关闭；n发送者 1接收者，这个接收者通过额外的goroutine通知发送者不要发送数据；n发送者 m接收者 任意一个goroutine都可以去通过一个额外的goroutine来通知每个goroutine

#### 控制协程(goroutine)的并发数量

总结：

如果不控制协程的并发数量，会导致内存被耗尽，所以我们可以通过带有缓冲区的chan来控制goroutine的并发数量，或者使用第三方的库，用pool技术。我们也可以修改系统参数，比如虚拟内存、修改同时打开文件数等

#### Go sync.Pool

总结：

使用sync.Pool可以减少内存分配，减少GC压力

使用sync.Pool：1.声明对象池 2.Get从对象池中拿 Put放回对象池

#### Go sync.Once

总结：

sync.Once是让函数全局只执行一次，多用于控制变量的初始化 变量须满足的三个条件：1.第一次访问的时候初始化；2.变量初始化过程中所有读都会阻塞直到初始化完成；3.变量只初始化一次，后常驻内存中

原理：就是一个标志有没有执行的flag以及互斥锁

一个结构体中要是函数经常被使用，可以把函数声明在第一个字段，运行更快

#### Go sync.Cond

