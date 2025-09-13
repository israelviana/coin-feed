package api

import (
	"net/http"

	"coin-feed/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CryptoHandler struct {
	fetchCryptoCurrencyMap *usecase.FetchCryptocurrencyMap
}

func NewCryptoHandler(fetchCryptoCurrencyMap *usecase.FetchCryptocurrencyMap) *CryptoHandler {
	return &CryptoHandler{
		fetchCryptoCurrencyMap: fetchCryptoCurrencyMap,
	}
}

func (h *CryptoHandler) FetchCryptoCurrencyMap(c *gin.Context) {
	cryptoCurrencyResponse, err := h.fetchCryptoCurrencyMap.Fetch(c)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, cryptoCurrencyResponse)
}

func (h *CryptoHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/crypto-currency", h.FetchCryptoCurrencyMap)
}
