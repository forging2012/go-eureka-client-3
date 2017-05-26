package eureka

import (
	"time"
	"encoding/json"
	"strings"
)

func Register(eurekaServerUrl,appName,port,securePort,instanceId,ipAddress string) {
	tpl := string(Template)
	tpl = strings.Replace(tpl, "${ipAddress}", ipAddress, -1)
	tpl = strings.Replace(tpl, "${port}", port, -1)
	tpl = strings.Replace(tpl, "${securePort}", securePort, -1)
	tpl = strings.Replace(tpl, "${instanceId}", instanceId, -1)
	tpl = strings.Replace(tpl, "${appName}", appName, -1)

	registerAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps/" + appName,
		Method:      "POST",
		ContentType: "application/json;charset=UTF-8",
		Body:        tpl,
	}
	var result bool
	for {
		result = doHttpRequest(registerAction)
		if result {
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

func Heartbeat(eurekaServerUrl,appName,ipAddress string) {
	heartbeatAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps/" + appName + "/" + ipAddress,
		Method:      "PUT",
		ContentType: "application/json;charset=UTF-8",
	}
	doHttpRequest(heartbeatAction)
}

func UnRegister(eurekaServerUrl,appName,ipAddress string) {
	unRegisterAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps/" + appName + "/" + ipAddress,
		ContentType: "application/json;charset=UTF-8",
		Method:      "DELETE",
	}
	doHttpRequest(unRegisterAction)
}

func GetServices(eurekaServerUrl string) ([]EurekaApplication, error) {
	var m EurekaApplicationsRootResponse
	queryAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps",
		Method:      "GET",
		Accept:      "application/json;charset=UTF-8",
		ContentType: "application/json;charset=UTF-8",
	}
	bytes, err := executeQuery(queryAction)
	if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			return nil, err
		}
		return m.Resp.Applications, nil
	}
}

func GetServiceInstances(eurekaServerUrl,appName string) ([]EurekaInstance, error) {
	var m EurekaServiceResponse
	queryAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps/" + appName,
		Method:      "GET",
		Accept:      "application/json;charset=UTF-8",
		ContentType: "application/json;charset=UTF-8",
	}
	bytes, err := executeQuery(queryAction)
	if err != nil {
		return nil, err
	} else {
		err := json.Unmarshal(bytes, &m)
		if err != nil {
			return nil, err
		}
		return m.Application.Instance, nil
	}
}