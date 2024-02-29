package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/BelyaevEI/wallet/internal/models"
)

func (service *Service) TransferFunds(writer http.ResponseWriter, request *http.Request) {

	var (
		buf      bytes.Buffer
		walletTo models.Wallet
	)

	ctx := request.Context()

	// Parsing wallet id for return balance
	parts := strings.Split(request.URL.Path, "/")
	walletID, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		service.log.Log.Error("reading wallet id from request is failed: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// We will get a wallet where you need to transfer funds
	_, err = buf.ReadFrom(request.Body)
	if err != nil {
		service.log.Log.Error("reading body from request is failed: ", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal data
	if err := json.Unmarshal(buf.Bytes(), &walletTo); err != nil {
		service.log.Log.Error("unmarshal data is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check exists wallet FROM
	ok, err := service.repository.CheckExists(ctx, uint32(walletID))
	if err != nil {
		service.log.Log.Error("check exists is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Wallet not exists FROM
	if !ok {
		service.log.Log.Info("wallet not exists")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Check exists wallet TO
	ok, err = service.repository.CheckExists(ctx, walletTo.ID)
	if err != nil {
		service.log.Log.Error("check exists is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Wallet not exists TO
	if !ok {
		service.log.Log.Info("wallet not exists")
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	// Let's check that there are enough funds in the wallet
	ok, err = service.repository.CheckFundsByID(ctx, uint32(walletID), walletTo.Amount)
	if err != nil {
		service.log.Log.Error("check funds is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insufficient funds
	if !ok {
		service.log.Log.Info("insufficient funds")
		writer.WriteHeader(http.StatusNotAcceptable)
		return
	}

	mes := models.Transfer{
		WalletIDFrom: uint32(walletID),
		WalletIDTo:   walletTo.ID,
		Amount:       walletTo.Amount,
	}

	// Send message to workers
	err = service.repository.SendMessageToWorker(ctx, mes)
	if err != nil {
		service.log.Log.Error("send message to broker is failed: ", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success transfer funds to wallet
	writer.WriteHeader(http.StatusOK)
}
