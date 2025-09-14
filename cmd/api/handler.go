package api

import (
	"net/http"
	"strings"

	"coin-feed/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CryptoHandler struct {
	getCryptoCurrencyMap  *usecase.FetchCryptocurrencyMap
	getCryptoCurrencyById *usecase.GetLatestCryptoCurrencyDataById
}

func NewCryptoHandler(fetchCryptoCurrencyMap *usecase.FetchCryptocurrencyMap, getLatestCryptoCurrencyDataById *usecase.GetLatestCryptoCurrencyDataById) *CryptoHandler {
	return &CryptoHandler{
		getCryptoCurrencyMap:  fetchCryptoCurrencyMap,
		getCryptoCurrencyById: getLatestCryptoCurrencyDataById,
	}
}

func (h *CryptoHandler) FetchCryptoCurrencyMap(c *gin.Context) {
	cryptoCurrencyResponse, err := h.getCryptoCurrencyMap.Run(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, cryptoCurrencyResponse)
}

func (h *CryptoHandler) GetCryptoCurrencyById(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	cryptoCurrencyResponse, err := h.getCryptoCurrencyById.Run(c, id)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, cryptoCurrencyResponse)
}

func (h *CryptoHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/crypto-currency", h.FetchCryptoCurrencyMap)
	r.GET("/crypto-currency/:id", h.GetCryptoCurrencyById)
}
