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

	retFile := retrievalFile{}
	err := c.ShouldBindJSON(&retFile)

	url := retFile.Url
	sellerID := retFile.Seller_id

	resp, err := http.Get(url)
	file, err := io.ReadAll(resp.Body)

	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if sellerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	productList, err := goods.ParseXlsx(file, sellerID)

	//save or update goods in the database
	stats, err := goods.SaveGoods(productList)

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func HandleGetGoods(c *gin.Context) {
	getGoodsreq := getGoodsrequest{}
	err := c.ShouldBindJSON(&getGoodsreq)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	offer_id, _ := strconv.Atoi(getGoodsreq.Offer_id)
	seller_id, _ := strconv.Atoi(getGoodsreq.Seller_id)

	goods_list := goods.GetGoods(offer_id, seller_id, getGoodsreq.GoodSubstring)

	c.JSON(http.StatusOK, gin.H{"goods": goods_list})
}
