package transfer

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"techtrain/token"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

type Block struct {
	Number string
}

func GachaTransfer(responseStatus *int, PRIVATE_KEY string, number uint32, court chan int) error {

	err := godotenv.Load(".env")
	if err != nil {
		*responseStatus = 401
		return err
	}

	INFURA_APIKEY := os.Getenv("INFURA_APIKEY")

	client, err := ethclient.Dial(INFURA_APIKEY)
	if err != nil {
		*responseStatus = 401
		return err
	}

	privateKey, err := crypto.HexToECDSA(PRIVATE_KEY)
	if err != nil {
		*responseStatus = 401
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		*responseStatus = 401
		return errors.New("No Such Public key...")
	}

	toAddress := common.HexToAddress("0x1Fa1520A45d5A28f2487D15915f8FF27FA538545")
	tokenAddress := common.HexToAddress("0x2813971A687011B1518731fB93D6C6a62cAeB2C4")
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		*responseStatus = 402
		return err
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, fromAddress)
	if err != nil {
		*responseStatus = 403
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		*responseStatus = 404
		return err
	}

	value := big.NewInt(0) // in wei (0 eth)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		*responseStatus = 404
		return err
	}

	transferFnSignature := []byte("transfer(address,uint256)")
	hash := crypto.NewKeccakState()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Printf("----------------Transcation infomation------------------------------------------------\n")
	fmt.Printf("Method ID: %s\n", hexutil.Encode(methodID))

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)
	fmt.Printf("To address: %s\n", hexutil.Encode(paddedAddress))

	amount := new(big.Int)
	amount.SetString(fmt.Sprint(number)+"000000000000000000", 10) // number tokens
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Printf("Token amount: %s\n", hexutil.Encode(paddedAmount))

	fmt.Printf("Account XY Balance: %d\n", bal)

	if bal.Cmp(amount) == -1 {
		*responseStatus = 403
		return errors.New("XY Balance not enough")
	}

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		*responseStatus = 404
		return err
	}
	fmt.Printf("Gas limit: %d\n", gasLimit*10)

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit*10, gasPrice.Mul(gasPrice, big.NewInt(10)), data)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		*responseStatus = 404
		return err
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		*responseStatus = 404
		return err
	}

	fmt.Printf("Tokens sent at TX: %s\n", signedTx.Hash().Hex())
	fmt.Printf("--------------------------------------------------------------------------------------\n")

	go Confirmation(client, signedTx, court)

	return nil
}

func Confirmation(client *ethclient.Client, signedTx *types.Transaction, court chan int) {
	fmt.Printf("Waiting for a blockchain transaction to be confirmed...\n")

	bind.WaitMined(context.Background(), client, signedTx)
	receipt, err := client.TransactionReceipt(context.Background(), signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status == 1 {
		fmt.Printf("Blockchain transaction success\n")
		court <- 1
	} else {
		fmt.Printf("Blockchain transaction false\n")
		court <- -1
	}
}

// func ConnInfura() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		fmt.Printf("cannot read .env: %v", err)
// 	}

// 	INFURA_APIKEY := os.Getenv("INFURA_APIKEY")
// 	fmt.Println(INFURA_APIKEY)

// 	client, err := rpc.Dial(INFURA_APIKEY)
// 	if err != nil {
// 		log.Fatalf("Could not connect to Infura: %v", err)
// 	}

// 	var lastBlock Block
// 	err = client.Call(&lastBlock, "eth_getBlockByNumber", "latest", true)
// 	if err != nil {
// 		fmt.Println("Cannot get the latest block:", err)
// 		return
// 	}

// 	fmt.Printf("Latest block: %v\n", lastBlock.Number)
// }
