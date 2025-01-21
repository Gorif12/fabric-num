package fabric

import (
	"context"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	bolt "go.etcd.io/bbolt"
)

const (
	BlocksBucket = "blocks"        // 存储区块数据
	LatestBucket = "latest_blocks" // 存储最新区块信息
)

// BlockData 区块数据结构
type BlockData struct {
	BlockNum  uint64    `json:"block_num"`
	BlockHash string    `json:"block_hash"`
	DataHash  string    `json:"data_hash"`
	PrevHash  string    `json:"prev_hash"`
	TxCount   int       `json:"tx_count"`
	SaveTime  time.Time `json:"save_time"`
}

// LatestBlock 最新区块信息
type LatestBlock struct {
	BlockNum uint64    `json:"block_num"`
	SaveTime time.Time `json:"save_time"`
}

// BlockEventListener 区块事件监听器
type BlockEventListener struct {
	sync.RWMutex
	networks map[string]*client.Network
	ctx      context.Context
	cancel   context.CancelFunc
	dataDir  string
	db       *bolt.DB
}

var (
	listener     *BlockEventListener
	listenerOnce sync.Once
)

// InitBlockListener 初始化区块监听器
func InitBlockListener(dataDir string) error {
	var initErr error
	listenerOnce.Do(func() {
		// 创建数据目录
		if err := os.MkdirAll(dataDir, 0755); err != nil {
			initErr = fmt.Errorf("创建数据目录失败：%w", err)
			return
		}

		// 打开BBolt数据库
		dbPath := filepath.Join(dataDir, "blocks.db")
		db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 10 * time.Second})
		if err != nil {
			initErr = fmt.Errorf("打开数据库失败：%w", err)
			return
		}

		// 创建Buckets
		if err := db.Update(func(tx *bolt.Tx) error {
			if _, err := tx.CreateBucketIfNotExists([]byte(BlocksBucket)); err != nil {
				return fmt.Errorf("创建blocks bucket失败: %w", err)
			}
			if _, err := tx.CreateBucketIfNotExists([]byte(LatestBucket)); err != nil {
				return fmt.Errorf("创建latest_blocks bucket失败: %w", err)
			}
			return nil
		}); err != nil {
			db.Close()
			initErr = fmt.Errorf("初始化数据库失败：%w", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		listener = &BlockEventListener{
			networks: make(map[string]*client.Network),
			ctx:      ctx,
			cancel:   cancel,
			dataDir:  dataDir,
			db:       db,
		}
	})

	return initErr
}

// AddNetwork 添加网络
func AddNetwork(orgName string, network *client.Network) error {
	if listener == nil {
		return fmt.Errorf("区块监听器未初始化")
	}

	listener.Lock()
	defer listener.Unlock()

	listener.networks[orgName] = network
	go listener.startBlockListener(orgName)

	return nil
}

// startBlockListener 启动区块监听
func (l *BlockEventListener) startBlockListener(orgName string) {
	network := l.networks[orgName]
	if network == nil {
		fmt.Printf("组织[%s]的网络未找到\n", orgName)
		return
	}

	lastBlockNum := l.getLastBlockNum(orgName)
	startBlock := lastBlockNum + 1

	events, err := network.BlockEvents(l.ctx, client.WithStartBlock(startBlock))
	if err != nil {
		fmt.Printf("创建区块事件请求失败：%v\n", err)
		return
	}

	for {
		select {
		case <-l.ctx.Done():
			return
		case block := <-events:
			l.saveBlock(orgName, block)
		}
	}
}

// getLastBlockNum 获取最后保存的区块号
func (l *BlockEventListener) getLastBlockNum(orgName string) uint64 {
	var lastBlock LatestBlock

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(LatestBucket))
		data := b.Get([]byte(orgName))
		if data == nil {
			return nil
		}
		return json.Unmarshal(data, &lastBlock)
	})

	if err != nil {
		fmt.Printf("获取最后区块号失败：%v\n", err)
		return 0
	}

	return lastBlock.BlockNum
}

// saveBlock 保存区块数据
func (l *BlockEventListener) saveBlock(orgName string, block *common.Block) {
	blockNum := block.GetHeader().GetNumber()

	// 计算区块哈希
	blockHeader := struct {
		Number       *big.Int
		PreviousHash []byte
		DataHash     []byte
	}{
		PreviousHash: block.GetHeader().GetPreviousHash(),
		DataHash:     block.GetHeader().GetDataHash(),
		Number:       new(big.Int).SetUint64(blockNum),
	}
	headerBytes, err := asn1.Marshal(blockHeader)
	if err != nil {
		fmt.Printf("序列化区块头失败：%v\n", err)
		return
	}
	blockHash := sha256.Sum256(headerBytes)

	// 准备区块数据
	blockData := BlockData{
		BlockNum:  blockNum,
		BlockHash: fmt.Sprintf("%x", blockHash[:]),
		DataHash:  fmt.Sprintf("%x", block.GetHeader().GetDataHash()),
		PrevHash:  fmt.Sprintf("%x", block.GetHeader().GetPreviousHash()),
		TxCount:   len(block.GetData().GetData()),
		SaveTime:  time.Now(),
	}

	// 使用事务保存数据
	err = l.db.Update(func(tx *bolt.Tx) error {
		// 保存区块数据
		blocksBucket := tx.Bucket([]byte(BlocksBucket))
		blockKey := fmt.Sprintf("%s_%d", orgName, blockNum)
		blockJSON, err := json.Marshal(blockData)
		if err != nil {
			return fmt.Errorf("序列化区块数据失败：%v", err)
		}
		if err := blocksBucket.Put([]byte(blockKey), blockJSON); err != nil {
			return fmt.Errorf("保存区块数据失败：%v", err)
		}

		// 更新最新区块信息
		latestBucket := tx.Bucket([]byte(LatestBucket))
		latestBlock := LatestBlock{
			BlockNum: blockNum,
			SaveTime: time.Now(),
		}
		latestJSON, err := json.Marshal(latestBlock)
		if err != nil {
			return fmt.Errorf("序列化最新区块信息失败：%v", err)
		}
		if err := latestBucket.Put([]byte(orgName), latestJSON); err != nil {
			return fmt.Errorf("保存最新区块信息失败：%v", err)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("保存区块失败：%v\n", err)
		return
	}

	fmt.Printf("已保存组织[%s]的区块[%d]\n", orgName, blockNum)
}

// Close 关闭监听器
func (l *BlockEventListener) Close() error {
	if l.cancel != nil {
		l.cancel()
	}
	if l.db != nil {
		return l.db.Close()
	}
	return nil
}

// GetBlockByNumber 根据组织名和区块号查询区块
func (l *BlockEventListener) GetBlockByNumber(orgName string, blockNum uint64) (*BlockData, error) {
	var blockData BlockData

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			return fmt.Errorf("blocks bucket不存在")
		}

		blockKey := fmt.Sprintf("%s_%d", orgName, blockNum)
		data := b.Get([]byte(blockKey))
		if data == nil {
			return fmt.Errorf("区块不存在")
		}

		return json.Unmarshal(data, &blockData)
	})

	if err != nil {
		return nil, err
	}

	return &blockData, nil
}

// BlockQueryResult 区块查询结果
type BlockQueryResult struct {
	Blocks   []*BlockData `json:"blocks"`    // 区块数据列表
	Total    int          `json:"total"`     // 总记录数
	PageSize int          `json:"page_size"` // 每页大小
	PageNum  int          `json:"page_num"`  // 当前页码
	HasMore  bool         `json:"has_more"`  // 是否还有更多数据
}

// GetBlocksByOrg 分页查询组织的区块列表（按区块号降序）
func (l *BlockEventListener) GetBlocksByOrg(orgName string, pageSize, pageNum int) (*BlockQueryResult, error) {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNum <= 0 {
		pageNum = 1
	}

	var result BlockQueryResult
	result.PageSize = pageSize
	result.PageNum = pageNum

	err := l.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BlocksBucket))
		if b == nil {
			return fmt.Errorf("blocks bucket不存在")
		}

		// 获取组织的最新区块号
		latestBucket := tx.Bucket([]byte(LatestBucket))
		if latestBucket == nil {
			return fmt.Errorf("latest_blocks bucket不存在")
		}

		var latestBlock LatestBlock
		latestData := latestBucket.Get([]byte(orgName))
		if latestData == nil {
			return fmt.Errorf("组织数据不存在")
		}
		if err := json.Unmarshal(latestData, &latestBlock); err != nil {
			return err
		}

		// 计算总记录数和分页参数
		result.Total = int(latestBlock.BlockNum) + 1 // 从0开始，所以要加1

		// 计算起始和结束区块号
		startIdx := result.Total - (pageNum * pageSize)
		endIdx := startIdx + pageSize
		if startIdx < 0 {
			startIdx = 0
		}
		if endIdx > result.Total {
			endIdx = result.Total
		}

		// 判断是否还有更多数据
		result.HasMore = startIdx > 0

		// 收集区块数据
		blocks := make([]*BlockData, 0, pageSize)
		for i := endIdx - 1; i >= startIdx; i-- {
			blockKey := fmt.Sprintf("%s_%d", orgName, i)
			data := b.Get([]byte(blockKey))
			if data != nil {
				var block BlockData
				if err := json.Unmarshal(data, &block); err != nil {
					return err
				}
				blocks = append(blocks, &block)
			}
		}
		result.Blocks = blocks

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}
