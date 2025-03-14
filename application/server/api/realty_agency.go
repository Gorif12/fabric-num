//这段代码定义了一个 HTTP API 处理器，用于处理与房产信息相关的请求，包括创建房产信息、查询房产信息、分页查询房产列表和分页查询区块列表。
//它使用了 Gin 框架来处理 HTTP 请求，并且与服务层交互来执行具体的业务逻辑。
package api

import (
	"application/service"
	"application/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RealtyAgencyHandler struct {
	realtyService *service.RealtyAgencyService
}

func NewRealtyAgencyHandler() *RealtyAgencyHandler {
	return &RealtyAgencyHandler{
		realtyService: &service.RealtyAgencyService{},
	}
}

// CreateRealEstate 创建房产信息（仅不动产登记机构组织可以调用）
func (h *RealtyAgencyHandler) CreateRealEstate(c *gin.Context) {
	var req struct {
		ID      string  `json:"id"`
		Address string  `json:"address"`
		Area    float64 `json:"area"`
		Owner   string  `json:"owner"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "格式错误")
		return
	}

	err := h.realtyService.CreateRealEstate(req.ID, req.Address, req.Area, req.Owner)
	if err != nil {
		utils.ServerError(c, "创建信息失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "创建成功", nil)
}

// QueryRealEstate 查询房产信息
func (h *RealtyAgencyHandler) QueryRealEstate(c *gin.Context) {
	id := c.Param("id")
	realEstate, err := h.realtyService.QueryRealEstate(id)
	if err != nil {
		utils.ServerError(c, "查询信息失败："+err.Error())
		return
	}

	utils.Success(c, realEstate)
}

// QueryRealEstateList 分页查询房产列表
func (h *RealtyAgencyHandler) QueryRealEstateList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")
	status := c.DefaultQuery("status", "")

	result, err := h.realtyService.QueryRealEstateList(int32(pageSize), bookmark, status)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryBlockList 分页查询区块列表
func (h *RealtyAgencyHandler) QueryBlockList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	result, err := h.realtyService.QueryBlockList(pageSize, pageNum)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}
