//这段代码定义了一个 TradingPlatformHandler 结构体，用于处理与交易平台相关的 HTTP 请求。它包含了创建交易、查询房产信息、查询交易信息、分页查询交易列表和分页查询区块列表的方法。
//每个方法都使用了 Gin 框架来处理 HTTP 请求，并调用了 service 包中的服务来完成具体的业务逻辑。
package api

import (
	"application/service"
	"application/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TradingPlatformHandler struct {
	tradingService *service.TradingPlatformService
}

func NewTradingPlatformHandler() *TradingPlatformHandler {
	return &TradingPlatformHandler{
		tradingService: &service.TradingPlatformService{},
	}
}

// CreateTransaction 生成交易（仅交易平台组织可以调用）
func (h *TradingPlatformHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		TxID         string  `json:"txId"`
		RealEstateID string  `json:"realEstateId"`
		Seller       string  `json:"seller"`
		Buyer        string  `json:"buyer"`
		Price        float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "交易信息格式错误")
		return
	}

	err := h.tradingService.CreateTransaction(req.TxID, req.RealEstateID, req.Seller, req.Buyer, req.Price)
	if err != nil {
		utils.ServerError(c, "生成交易失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "交易创建成功", nil)
}

// QueryRealEstate 查询房产信息
func (h *TradingPlatformHandler) QueryRealEstate(c *gin.Context) {
	id := c.Param("id")
	realEstate, err := h.tradingService.QueryRealEstate(id)
	if err != nil {
		utils.ServerError(c, "查询信息失败："+err.Error())
		return
	}

	utils.Success(c, realEstate)
}

// QueryTransaction 查询交易信息
func (h *TradingPlatformHandler) QueryTransaction(c *gin.Context) {
	txID := c.Param("txId")
	transaction, err := h.tradingService.QueryTransaction(txID)
	if err != nil {
		utils.ServerError(c, "查询交易信息失败："+err.Error())
		return
	}

	utils.Success(c, transaction)
}

// QueryTransactionList 分页查询交易列表
func (h *TradingPlatformHandler) QueryTransactionList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")
	status := c.DefaultQuery("status", "")

	result, err := h.tradingService.QueryTransactionList(int32(pageSize), bookmark, status)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryBlockList 分页查询区块列表
func (h *TradingPlatformHandler) QueryBlockList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	result, err := h.tradingService.QueryBlockList(pageSize, pageNum)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}
