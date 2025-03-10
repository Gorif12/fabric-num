1. 智能合约结构定义
   // SmartContract 定义房地产交易智能合约
type SmartContract struct {
    contractapi.Contract  // 继承自Fabric链码基础合约
}
// 文档类型常量（用于构造复合键）
const (
    REAL_ESTATE = "RE"  // 房产键前缀
    TRANSACTION = "TX"  // 交易键前缀
)
// 房产状态枚举
type RealEstateStatus string
const (
    NORMAL         RealEstateStatus = "NORMAL"         // 正常状态
    IN_TRANSACTION RealEstateStatus = "IN_TRANSACTION" // 交易中状态
)
// 交易状态枚举
type TransactionStatus string
const (
    PENDING   TransactionStatus = "PENDING"   // 待付款
    COMPLETED TransactionStatus = "COMPLETED" // 已完成
)
// RealEstate 定义房产数据结构
type RealEstate struct {
    ID              string           // 唯一标识
    PropertyAddress string           // 地址
    Area            float64          // 面积
    CurrentOwner    string           // 当前所有者
    Status          RealEstateStatus // 状态
    CreateTime      time.Time        // 创建时间
    UpdateTime      time.Time        // 更新时间
}
// Transaction 定义交易数据结构
type Transaction struct {
    ID           string            // 交易ID
    RealEstateID string            // 关联的房产ID
    Seller       string            // 卖家
    Buyer        string            // 买家
    Price        float64           // 价格
    Status       TransactionStatus // 状态
    CreateTime   time.Time         // 创建时间
    UpdateTime   time.Time         // 更新时间
}
// QueryResult 分页查询结果结构
type QueryResult struct {
    Records             []interface{} // 数据记录列表
    RecordsCount        int32         // 返回记录数
    Bookmark            string        // 分页书签
    FetchedRecordsCount int32         // 总记录数
}
//逻辑说明：
   复合键策略：通过 REAL_ESTATE 和 TRANSACTION 前缀结合状态（如 NORMAL）生成唯一键，便于分类查询。
   状态枚举：明确限定房产和交易的合法状态，避免非法状态输入

   
2. 权限控制与通用方法
// 获取客户端所属组织的MSP ID
func (s *SmartContract) getClientIdentityMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
    clientID, err := cid.New(ctx.GetStub())  // 解析客户端身份
    if err != nil {
        return "", err
    }
    return clientID.GetMSPID()  // 返回组织标识（如Org1MSP）
}
// 创建复合键（格式：前缀_属性1_属性2...）
func (s *SmartContract) getCompositeKey(ctx contractapi.TransactionContextInterface, objectType string, attributes []string) (string, error) {
    return ctx.GetStub().CreateCompositeKey(objectType, attributes)
}
// 通用方法：反序列化读取状态
func (s *SmartContract) getState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
    data, err := ctx.GetStub().GetState(key)  // 从账本读取
    if err != nil || data == nil {
        return err
    }
    return json.Unmarshal(data, value)  // JSON反序列化
}

// 通用方法：序列化保存状态
func (s *SmartContract) putState(ctx contractapi.TransactionContextInterface, key string, value interface{}) error {
    data, err := json.Marshal(value)  // JSON序列化
    if err != nil {
        return err
    }
    return ctx.GetStub().PutState(key, data)  // 写入账本
}
  逻辑说明：
    MSP ID验证：确保只有授权组织（如不动产机构）可执行敏感操作。
    复合键管理：通过 CreateCompositeKey 实现数据分类存储（如按状态分组）。

    
3. 核心业务方法
 3.1 创建房产（仅限不动产机构）
  func (s *SmartContract) CreateRealEstate(ctx contractapi.TransactionContextInterface, id, address string, area float64, owner string, createTime time.Time) error {
    // 权限校验
    if mspID, _ := s.getClientIdentityMSPID(ctx); mspID != REALTY_ORG_MSPID {
        return errors.New("无权操作")
    }
    // 参数校验
    if area <= 0 || owner == "" { /*...*/ }
    // 检查房产唯一性（所有状态）
    for _, status := range []RealEstateStatus{NORMAL, IN_TRANSACTION} {
        key, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(status), id})
        if exists, _ := ctx.GetStub().GetState(key); exists != nil {
            return errors.New("房产已存在")
        }
    }
    // 保存新房产（状态：NORMAL）
    key, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), id})
    return s.putState(ctx, key, &RealEstate{
        ID:              id,
        PropertyAddress: address,
        Area:            area,
        CurrentOwner:    owner,
        Status:          NORMAL,
        CreateTime:      createTime,
        UpdateTime:      createTime,
    })
}

3.2 创建交易（仅限交易平台）
func (s *SmartContract) CreateTransaction(ctx contractapi.TransactionContextInterface, txID, realEstateID, seller, buyer string, price float64, createTime time.Time) error {
    // 权限校验...
    
    // 验证房产归属
    realEstateKey, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), realEstateID})
    realEstate, _ := s.getState(ctx, realEstateKey, &RealEstate{})
    if realEstate.CurrentOwner != seller {
        return errors.New("卖家无权出售")
    }
    // 更新房产状态为交易中
    ctx.GetStub().DelState(realEstateKey)  // 删除旧状态记录
    newRealEstateKey, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(IN_TRANSACTION), realEstateID})
    s.putState(ctx, newRealEstateKey, realEstate)
    // 保存新交易（状态：PENDING）
    txKey, _ := s.getCompositeKey(ctx, TRANSACTION, []string{string(PENDING), txID})
    return s.putState(ctx, txKey, &Transaction{
        ID:           txID,
        RealEstateID: realEstateID,
        Seller:       seller,
        Buyer:        buyer,
        Price:        price,
        Status:       PENDING,
        CreateTime:   createTime,
        UpdateTime:   createTime,
    })
}

3.3 完成交易（仅限银行）
func (s *SmartContract) CompleteTransaction(ctx contractapi.TransactionContextInterface, txID string, updateTime time.Time) error {
    // 权限校验...
    // 查询待完成交易
    txKey, _ := s.getCompositeKey(ctx, TRANSACTION, []string{string(PENDING), txID})
    tx, _ := s.getState(ctx, txKey, &Transaction{})
    // 更新房产所有权
    realEstateKey, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(IN_TRANSACTION), tx.RealEstateID})
    realEstate, _ := s.getState(ctx, realEstateKey, &RealEstate{})
    realEstate.CurrentOwner = tx.Buyer
    realEstate.Status = NORMAL
    // 删除旧记录并保存新状态
    ctx.GetStub().DelState(realEstateKey)
    ctx.GetStub().DelState(txKey)
    newRealEstateKey, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(NORMAL), tx.RealEstateID})
    newTxKey, _ := s.getCompositeKey(ctx, TRANSACTION, []string{string(COMPLETED), txID})
    s.putState(ctx, newRealEstateKey, realEstate)
    s.putState(ctx, newTxKey, tx)
    return nil
}
业务逻辑总结：
   创建房产：不动产机构验证身份后，生成初始状态（NORMAL）的房产记录。
   发起交易：交易平台验证卖家所有权后，将房产标记为交易中（IN_TRANSACTION），记录交易为待处理（PENDING）。
   完成交易：银行确认支付后，转移所有权，房产恢复正常状态，交易标记为已完成（COMPLETED）。

   
4. 查询方法
// 按ID查询房产（自动检查所有状态）
func (s *SmartContract) QueryRealEstate(ctx contractapi.TransactionContextInterface, id string) (*RealEstate, error) {
    for _, status := range []RealEstateStatus{NORMAL, IN_TRANSACTION} {
        key, _ := s.getCompositeKey(ctx, REAL_ESTATE, []string{string(status), id})
        if data, _ := ctx.GetStub().GetState(key); data != nil {
            var re RealEstate
            json.Unmarshal(data, &re)
            return &re, nil
        }
    }
    return nil, errors.New("未找到")
}
// 分页查询房产列表
func (s *SmartContract) QueryRealEstateList(ctx contractapi.TransactionContextInterface, pageSize int32, bookmark, status string) (*QueryResult, error) {
    // 根据状态构造查询条件
    iterator, metadata, _ := ctx.GetStub().GetStateByPartialCompositeKeyWithPagination(
        REAL_ESTATE,
        []string{status},
        pageSize,
        bookmark,
    )
    defer iterator.Close()
    // 遍历结果集
    var records []interface{}
    for iterator.HasNext() {
        item, _ := iterator.Next()
        var re RealEstate
        json.Unmarshal(item.Value, &re)
        records = append(records, re)
    }
    return &QueryResult{
        Records:             records,
        Bookmark:            metadata.Bookmark,
        FetchedRecordsCount: metadata.FetchedRecordsCount,
    }, nil
}
查询逻辑：
   按ID查询：遍历所有可能状态，确保找到目标房产。
   分页查询：利用 GetStateByPartialCompositeKeyWithPagination 实现高效分页。

   
6. 链码入口
func main() {
    // 初始化链码实例
    chaincode, _ := contractapi.NewChaincode(&SmartContract{})
    // 启动服务
    if err := chaincode.Start(); err != nil {
        log.Panicf("链码启动失败: %v", err)
    }
}
