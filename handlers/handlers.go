package handlers

import (
	"github.com/gin-gonic/gin"
	"io"
	"merchant-experience/goods"
	"net/http"
	"strconv"
)

type retrievalFile struct {
	Url       string `json:"url"`
	Seller_id string `json:"seller_id"`
}

type getGoodsrequest struct {
	Seller_id     string `json:"seller_id"`
	Offer_id      string `json:"offer_id"`
	GoodSubstring string `json:"goodSubstring"`
}

func HandleXlsxProcessing(c *gin.Context) {
	var retFile retrievalFile
	if err := c.ShouldBindJSON(&retFile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if retFile.Url == "" || retFile.Seller_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL and seller ID are required"})
		return
	}

	resp, err := http.Get(retFile.Url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	file, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productList, err := goods.ParseXlsx(file, retFile.Seller_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats, err := goods.SaveGoods(productList)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func HandleGetGoods(c *gin.Context) {
	var getGoodsreq getGoodsrequest
	if err := c.ShouldBindJSON(&getGoodsreq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	offerID, err := strconv.ParseFloat(getGoodsreq.Offer_id, 64)
	if err != nil && getGoodsreq.Offer_id != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offer ID"})
		return
	}

	sellerID, err := strconv.ParseFloat(getGoodsreq.Seller_id, 64)
	if err != nil && getGoodsreq.Seller_id != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid seller ID"})
		return
	}

	goodsList, err := goods.GetGoods(int(offerID), int(sellerID), getGoodsreq.GoodSubstring)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"goods": goodsList})
}
