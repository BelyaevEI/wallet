package service

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func (service *Service) GetBalance(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	// Parsing wallet id for return balance
	parts := strings.Split(request.URL.Path, "/")
	walletID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		service.log.Log.Error("reading wallet id from request is failed: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check exists wallet
	ok, err := service.repository.CheckExists(ctx, uint32(walletID))
	if err != nil {
		service.log.Log.Error("check exists is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Wallet not exists
	if !ok {
		service.log.Log.Info("wallet not exists")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Getting balance by wallet id
	balance, err := service.repository.GetBalanceByID(ctx, uint32(walletID))
	if err != nil {
		service.log.Log.Error("geting balance is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf("%d", balance)))
}
