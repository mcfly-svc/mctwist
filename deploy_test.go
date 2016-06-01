package main_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/mikec/mctwist"
	"github.com/mcfly-svc/mcfly/api/apidata"
	"github.com/mcfly-svc/mcfly/client"
	"github.com/mcfly-svc/mcfly/client/mockclient"
	"github.com/mcfly-svc/mcfly/provider"
	"github.com/mcfly-svc/mcfly/provider/mockprovider"
	"github.com/stretchr/testify/assert"
)

var apiErrOutput = "mcflyapi responded with an error:"

type MockClientReturn struct {
	CR   *client.ClientResponse
	Resp *http.Response
	Err  error
}

type MockBuildDataReturn struct {
	BuildData *provider.BuildData
	Err       error
}

type DeployerTest struct {
	Message                   string
	Token                     string
	BuildHandle               string
	ProjectHandle             string
	MockGetBuildDataReturn    *MockBuildDataReturn
	MockSaveBuildReturn       *MockClientReturn
	ExpectSaveBuildCalledWith *apidata.BuildReq
	ExpectErrorOutput         string
}

func TestStartDeploy(t *testing.T) {
	validProvider := "jabroni.com"
	validToken := "abc123"
	validBuildHandle := "mock_valid_build_handle"
	validProjectHandle := "mock_valid_project_handle"
	validBuildUrl := strPtr("mockprojecturl://")
	validBuildDataReturn := &MockBuildDataReturn{
		BuildData: &provider.BuildData{
			Handle: validBuildHandle,
			Url:    validBuildUrl,
			Config: provider.NewDefaultBuildConfig(),
		},
		Err: nil,
	}
	validSaveBuildCalledWith := &apidata.BuildReq{
		Handle:        validBuildHandle,
		ProjectHandle: validProjectHandle,
		Provider:      validProvider,
		ProviderUrl:   validBuildUrl,
	}
	validSaveBuildReturn := &MockClientReturn{
		CR:   &client.ClientResponse{StatusCode: 200, Data: "mock"},
		Resp: nil,
		Err:  nil,
	}

	successCaseTest := DeployerTest{
		Message:                   "When all data is valid",
		Token:                     validToken,
		BuildHandle:               validBuildHandle,
		ProjectHandle:             validProjectHandle,
		MockGetBuildDataReturn:    validBuildDataReturn,
		ExpectSaveBuildCalledWith: validSaveBuildCalledWith,
		MockSaveBuildReturn:       validSaveBuildReturn,
		ExpectErrorOutput:         "",
	}

	saveBuildErrTest := successCaseTest
	saveBuildErrTest.Message = "When SaveBuild returns an error"
	saveBuildErrTest.MockSaveBuildReturn = &MockClientReturn{
		CR:   nil,
		Resp: nil,
		Err:  errors.New("mock error"),
	}
	saveBuildErrTest.ExpectErrorOutput = "mock error"

	saveBuildErrResponseTest := successCaseTest
	saveBuildErrResponseTest.Message = "When SaveBuild responds with an error"
	saveBuildErrResponseTest.MockSaveBuildReturn = &MockClientReturn{
		CR:   &client.ClientResponse{StatusCode: 400, Data: "mock error"},
		Resp: nil,
		Err:  nil,
	}
	saveBuildErrResponseTest.ExpectErrorOutput = fmt.Sprintf("%s mock error\n", apiErrOutput)

	tests := []DeployerTest{
		successCaseTest,
		saveBuildErrTest,
		saveBuildErrResponseTest,
	}

	for _, test := range tests {
		logger := &MockLogger{}

		mc := new(mockclient.MockClient)
		ret := test.MockSaveBuildReturn
		mc.On("SaveBuild", test.ExpectSaveBuildCalledWith).Return(ret.CR, ret.Resp, ret.Err)

		mp := new(mockprovider.MockProvider)
		mp.On("Key").Return("jabroni.com")
		bdRet := test.MockGetBuildDataReturn
		mp.On("GetBuildData", test.Token, test.BuildHandle, test.ProjectHandle).Return(bdRet.BuildData, bdRet.Err)

		d := &main.Deployer{
			Logger:         logger,
			McflyClient:     mc,
			SourceProvider: mp,
			Token:          test.Token,
			BuildHandle:    test.BuildHandle,
			ProjectHandle:  test.ProjectHandle,
		}
		d.StartDeploy()
		assert.Equal(t, test.ExpectErrorOutput, logger.ErrorOutput, fmt.Sprintf("StartDeploy(): %s", test.Message))
	}
}

func strPtr(s string) *string {
	return &s
}
