package rds

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// ClearDedicatedHost invokes the rds.ClearDedicatedHost API synchronously
// api document: https://help.aliyun.com/api/rds/cleardedicatedhost.html
func (client *Client) ClearDedicatedHost(request *ClearDedicatedHostRequest) (response *ClearDedicatedHostResponse, err error) {
	response = CreateClearDedicatedHostResponse()
	err = client.DoAction(request, response)
	return
}

// ClearDedicatedHostWithChan invokes the rds.ClearDedicatedHost API asynchronously
// api document: https://help.aliyun.com/api/rds/cleardedicatedhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearDedicatedHostWithChan(request *ClearDedicatedHostRequest) (<-chan *ClearDedicatedHostResponse, <-chan error) {
	responseChan := make(chan *ClearDedicatedHostResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.ClearDedicatedHost(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// ClearDedicatedHostWithCallback invokes the rds.ClearDedicatedHost API asynchronously
// api document: https://help.aliyun.com/api/rds/cleardedicatedhost.html
// asynchronous document: https://help.aliyun.com/document_detail/66220.html
func (client *Client) ClearDedicatedHostWithCallback(request *ClearDedicatedHostRequest, callback func(response *ClearDedicatedHostResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *ClearDedicatedHostResponse
		var err error
		defer close(result)
		response, err = client.ClearDedicatedHost(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// ClearDedicatedHostRequest is the request struct for api ClearDedicatedHost
type ClearDedicatedHostRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DedicatedHostId      string           `position:"Query" name:"DedicatedHostId"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	FailoverMode         string           `position:"Query" name:"FailoverMode"`
}

// ClearDedicatedHostResponse is the response struct for api ClearDedicatedHost
type ClearDedicatedHostResponse struct {
	*responses.BaseResponse
	RequestId       string `json:"RequestId" xml:"RequestId"`
	TaskId          string `json:"TaskId" xml:"TaskId"`
	DedicatedHostId string `json:"DedicatedHostId" xml:"DedicatedHostId"`
}

// CreateClearDedicatedHostRequest creates a request to invoke ClearDedicatedHost API
func CreateClearDedicatedHostRequest() (request *ClearDedicatedHostRequest) {
	request = &ClearDedicatedHostRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("Rds", "2014-08-15", "ClearDedicatedHost", "rds", "openAPI")
	request.Method = requests.POST
	return
}

// CreateClearDedicatedHostResponse creates a response to parse from ClearDedicatedHost response
func CreateClearDedicatedHostResponse() (response *ClearDedicatedHostResponse) {
	response = &ClearDedicatedHostResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}
