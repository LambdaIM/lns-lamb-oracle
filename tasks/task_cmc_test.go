package tasks

import (
	"lns-lamb-oracle/config"
	"os"
	"path/filepath"
	"testing"
)

func TestCmc_GetLambPriceFromMarket(t *testing.T) {
	pwd, err := os.Getwd()
	if nil != err {
		t.Fatal(err)
	}
	cfg, err := config.NewConfig(filepath.Dir(pwd))
	task := NewCmcTask(cfg)

	err = task.GetLambPriceFromMarket([]byte("0123456789abcdef"))
	if nil != err {
		t.Error(err)
	}
}

func TestCmc_SetLambPriceToOracle(t *testing.T) {
	pwd, err := os.Getwd()
	if nil != err {
		t.Fatal(err)
	}
	cfg, err := config.NewConfig(filepath.Dir(pwd))
	task := NewCmcTask(cfg)

	err = task.GetLambPriceFromMarket([]byte("0123456789abcdef"))
	if nil != err {
		t.Error(err)
	}

	tx, err := task.SetLambPriceToOracle([]byte("0123456789abcdef"))
	if nil != err {
		t.Error(err)
	}
	t.Log("tx:", tx.Hash())
}
