package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chrismrivera/cmd"
	"github.com/mikec/msplapi/client"
	"github.com/mikec/msplapi/mq"
)

var cmdr *cmd.App = cmd.NewApp()
var cfg *Config
var mcTwistLogger *McTwistLogger = &McTwistLogger{}

func main() {

	_cfg, err := NewConfigFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	cfg = _cfg

	cmdr.AddCommand(runCmd)

	cmdr.Description = "[mctwist] build deployment manager for mspl"
	if err := cmdr.Run(os.Args); err != nil {
		if ue, ok := err.(*cmd.UsageErr); ok {
			ue.ShowUsage()
		} else {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		}
		os.Exit(1)
	}

}

func Run() error {
	msgChannel := mq.CreateChannel(cfg.RabbitMQUrl)
	defer msgChannel.CloseConnection()
	defer msgChannel.CloseChannel()

	msgs, err := msgChannel.ReceiveDeployQueueMessage()
	if err != nil {
		return err
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			HandleDeployMessage(d.Body)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
	return nil
}

var runCmd = cmd.NewCommand(
	"run", "Build", "Runs mctwist",
	func(cmd *cmd.Command) {},
	func(cmd *cmd.Command) error {
		return Run()
	},
)

func checkErrorResponse(logger Logger, cr *client.ClientResponse, err error) bool {
	if err != nil {
		logger.Error(err)
		return true
	}
	if cr.StatusCode != 200 {
		logger.ApiError(cr.Data)
		return true
	}
	return false
}
