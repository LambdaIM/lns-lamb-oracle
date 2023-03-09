package tasks

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	jsoniter "github.com/json-iterator/go"
	"io"
	"lns-lamb-oracle/config"
	"lns-lamb-oracle/contract"
	"lns-lamb-oracle/utils"
	"log"
	"math/big"
	"net/http"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Cmc struct {
	*config.Config
	value *big.Int
}

func NewCmcTask(cfg *config.Config) *Cmc {
	return &Cmc{Config: cfg}
}

func (c *Cmc) GetLambPriceFromMarket(secret []byte) error {
	rawKey, err := base64.StdEncoding.DecodeString(c.Market.ApiKey)
	if nil != err {
		return err
	}

	cmcKey, err := utils.Decrypt(rawKey, secret)
	if nil != err {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, c.Market.Url, nil)
	if nil != err {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-CMC_PRO_API_KEY", string(cmcKey))

	resp, err := client.Do(req)
	if nil != err {
		return err
	}

	defer func() {
		if err = resp.Body.Close(); nil != err {
			log.Println("close response body error:", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if nil != err {
		return err
	}

	price := json.Get(body, "data", "LAMB", "quote", "USD", "price").ToFloat64()

	value := toGWei(big.NewFloat(price))
	if 0 < value.Cmp(big.NewInt(0)) {
		return errors.New("invalid LAMB price")
	}
	c.value = toGWei(big.NewFloat(price))
	return nil
}

func toGWei(v *big.Float) *big.Int {
	wei := new(big.Int)
	new(big.Float).Mul(v, big.NewFloat(1e18)).Int(wei)
	return new(big.Int).Div(wei, big.NewInt(1e9))
}

func (c *Cmc) SetLambPriceToOracle(secret []byte) (*types.Transaction, error) {
	rawKey, err := base64.StdEncoding.DecodeString(c.Lambda.PrivateKey)
	if nil != err {
		return nil, err
	}

	lambdaKey, err := utils.Decrypt(rawKey, secret)
	if nil != err {
		return nil, err
	}

	client, err := ethclient.Dial(c.Lambda.RpcUrl)
	if nil != err {
		return nil, err
	}
	defer client.Close()

	chainId, err := client.ChainID(context.Background())
	if nil != err {
		return nil, err
	}

	priKey, err := crypto.HexToECDSA(string(lambdaKey))
	if nil != err {
		return nil, err
	}

	ops, err := bind.NewKeyedTransactorWithChainID(priKey, chainId)
	if nil != err {
		return nil, err
	}

	oracle, err := contract.NewContract(common.HexToAddress(c.Lambda.OracleAddr), client)
	if nil != err {
		return nil, err
	}

	tx, err := oracle.Set(ops, c.value)
	if nil != err {
		return nil, err
	}

	return tx, nil
}
