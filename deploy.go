package main

import (
	"encoding/json"
	"fmt"

	"github.com/mikec/msplapi/api/apidata"
	"github.com/mikec/msplapi/client"
	"github.com/mikec/msplapi/mq"
	"github.com/mikec/msplapi/provider"
)

func HandleDeployMessage(messageBody []byte) {

	msg := mq.DeployQueueMessage{}
	err := json.Unmarshal(messageBody, &msg)
	if err != nil {
		mcTwistLogger.Error(err)
		return
	}
	sp, err := provider.GetSourceProvider(msg.SourceProvider, nil)
	if err != nil {
		mcTwistLogger.Error(err)
		return
	}

	d := Deployer{
		Logger:         mcTwistLogger,
		MsplClient:     client.NewMsplClient(cfg.ApiUrl, ""),
		SourceProvider: sp,
		Token:          msg.SourceProviderAccessToken,
		BuildHandle:    msg.BuildHandle,
		ProjectHandle:  msg.ProjectHandle,
	}

	d.StartDeploy()
}

type Deployer struct {
	Logger         Logger
	MsplClient     client.Client
	SourceProvider provider.SourceProvider
	Token          string
	BuildHandle    string
	ProjectHandle  string
}

func (d *Deployer) StartDeploy() {
	spKey := d.SourceProvider.Key()
	buildData, err := d.SourceProvider.GetBuildData(d.Token, d.BuildHandle, d.ProjectHandle)
	if err != nil {
		d.Logger.Error(err)
		return
	}

	// TODO: make sure this call is authenticated in some way...
	cr, _, err := d.MsplClient.SaveBuild(&apidata.BuildReq{
		Handle:        d.BuildHandle,
		ProjectHandle: d.ProjectHandle,
		Provider:      spKey,
		ProviderUrl:   buildData.Url,
	})
	if checkErrorResponse(d.Logger, cr, err) {
		return
	}

	d.Logger.Log(fmt.Sprintf("Saved %s build `%s` for project `%s`", spKey, d.BuildHandle, d.ProjectHandle))

	// NEXT:

	// 1. get source code from provider .. d.SourceProvider.GetBuildSource()

	// 2. save the source code to AWS ... (no provider implementation needed)

	// 3. call to msplapi to set build status to `completed`

}
