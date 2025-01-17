package credentials

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lfhy/alismsproxyapi/src/github.com/alibabacloud-go/tea/tea"

	"github.com/aliyun/credentials-go/credentials/request"
	"github.com/aliyun/credentials-go/credentials/utils"
)

var securityCredURL = "http://100.100.100.200/latest/meta-data/ram/security-credentials/"

// EcsRAMRoleCredential is a kind of credential
type EcsRAMRoleCredential struct {
	*credentialUpdater
	RoleName          string
	sessionCredential *sessionCredential
	runtime           *utils.Runtime
}

type ecsRAMRoleResponse struct {
	Code            string `json:"Code" xml:"Code"`
	AccessKeyId     string `json:"AccessKeyId" xml:"AccessKeyId"`
	AccessKeySecret string `json:"AccessKeySecret" xml:"AccessKeySecret"`
	SecurityToken   string `json:"SecurityToken" xml:"SecurityToken"`
	Expiration      string `json:"Expiration" xml:"Expiration"`
}

func newEcsRAMRoleCredential(roleName string, runtime *utils.Runtime) *EcsRAMRoleCredential {
	return &EcsRAMRoleCredential{
		RoleName:          roleName,
		credentialUpdater: new(credentialUpdater),
		runtime:           runtime,
	}
}

// GetAccessKeyId reutrns  EcsRAMRoleCredential's AccessKeyId
// if AccessKeyId is not exist or out of date, the function will update it.
func (e *EcsRAMRoleCredential) GetAccessKeyId() (*string, error) {
	if e.sessionCredential == nil || e.needUpdateCredential() {
		err := e.updateCredential()
		if err != nil {
			return tea.String(""), err
		}
	}
	return tea.String(e.sessionCredential.AccessKeyId), nil
}

// GetAccessSecret reutrns  EcsRAMRoleCredential's AccessKeySecret
// if AccessKeySecret is not exist or out of date, the function will update it.
func (e *EcsRAMRoleCredential) GetAccessKeySecret() (*string, error) {
	if e.sessionCredential == nil || e.needUpdateCredential() {
		err := e.updateCredential()
		if err != nil {
			return tea.String(""), err
		}
	}
	return tea.String(e.sessionCredential.AccessKeySecret), nil
}

// GetSecurityToken reutrns  EcsRAMRoleCredential's SecurityToken
// if SecurityToken is not exist or out of date, the function will update it.
func (e *EcsRAMRoleCredential) GetSecurityToken() (*string, error) {
	if e.sessionCredential == nil || e.needUpdateCredential() {
		err := e.updateCredential()
		if err != nil {
			return tea.String(""), err
		}
	}
	return tea.String(e.sessionCredential.SecurityToken), nil
}

// GetBearerToken is useless for EcsRAMRoleCredential
func (e *EcsRAMRoleCredential) GetBearerToken() *string {
	return tea.String("")
}

// GetType reutrns  EcsRAMRoleCredential's type
func (e *EcsRAMRoleCredential) GetType() *string {
	return tea.String("ecs_ram_role")
}

func getRoleName() (string, error) {
	runtime := utils.NewRuntime(1, 1, "", "")
	request := request.NewCommonRequest()
	request.URL = securityCredURL
	request.Method = "GET"
	content, err := doAction(request, runtime)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (e *EcsRAMRoleCredential) updateCredential() (err error) {
	if e.runtime == nil {
		e.runtime = new(utils.Runtime)
	}
	request := request.NewCommonRequest()
	if e.RoleName == "" {
		e.RoleName, err = getRoleName()
		if err != nil {
			return fmt.Errorf("refresh Ecs sts token err: %s", err.Error())
		}
	}
	request.URL = securityCredURL + e.RoleName
	request.Method = "GET"
	content, err := doAction(request, e.runtime)
	if err != nil {
		return fmt.Errorf("refresh Ecs sts token err: %s", err.Error())
	}
	var resp *ecsRAMRoleResponse
	err = json.Unmarshal(content, &resp)
	if err != nil {
		return fmt.Errorf("refresh Ecs sts token err: Json Unmarshal fail: %s", err.Error())
	}
	if resp.Code != "Success" {
		return fmt.Errorf("refresh Ecs sts token err: Code is not Success")
	}
	if resp.AccessKeyId == "" || resp.AccessKeySecret == "" || resp.SecurityToken == "" || resp.Expiration == "" {
		return fmt.Errorf("refresh Ecs sts token err: AccessKeyId: %s, AccessKeySecret: %s, SecurityToken: %s, Expiration: %s", resp.AccessKeyId, resp.AccessKeySecret, resp.SecurityToken, resp.Expiration)
	}

	expirationTime, err := time.Parse("2006-01-02T15:04:05Z", resp.Expiration)
	e.lastUpdateTimestamp = time.Now().Unix()
	e.credentialExpiration = int(expirationTime.Unix() - time.Now().Unix())
	e.sessionCredential = &sessionCredential{
		AccessKeyId:     resp.AccessKeyId,
		AccessKeySecret: resp.AccessKeySecret,
		SecurityToken:   resp.SecurityToken,
	}

	return
}
