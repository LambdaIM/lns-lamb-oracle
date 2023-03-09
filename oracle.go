package main

import (
	"errors"
	"fmt"
	"golang.org/x/term"
	"lns-lamb-oracle/config"
	"lns-lamb-oracle/tasks"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiKeySecret := secret("Your secret for apiKey: ")
	priKeySecret := secret("Your secret for priKey: ")

	log.Println(string(apiKeySecret), string(priKeySecret))
	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt)

	pwd, err := os.Getwd()
	fatal(err)

	cfg, err := config.NewConfig(pwd)
	fatal(err)

	task := tasks.NewCmcTask(cfg)

	ticker := time.NewTicker(time.Duration(cfg.Lambda.DurationH) * time.Hour)
	defer ticker.Stop()
	log.Println("lamb oracle process running...")
	for {
		select {
		case <-ticker.C:
			err = task.GetLambPriceFromMarket(apiKeySecret)
			if nil != err {
				log.Println("execute task error:", err)
				continue
			}
			tx, err := task.SetLambPriceToOracle(priKeySecret)
			if nil != err {
				log.Println("update LAMB price to oracle error:", err)
				continue
			}
			log.Println("update LAMB price to oracle success.", tx.Hash())
		case <-interrupt:
			log.Println("system INTERRUPT signal received. Exiting...")
			return
		}
	}
}

func secret(message string) []byte {
	_, err := fmt.Fprint(os.Stdout, message)
	fatal(err)
	sec, err := term.ReadPassword(int(syscall.Stdin))
	fatal(err)
	length := len(sec)
	if 0 == length || (16 != length && 24 != length && 32 != length) {
		fmt.Println()
		fatal(errors.New("wrong secret entered, length of secret must be 16 or 24 or 32 bytes"))
	}
	fmt.Println()
	return sec
}

func fatal(err error) {
	if nil != err {
		log.Fatal(err)
	}
}
