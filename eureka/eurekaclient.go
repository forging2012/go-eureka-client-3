package eureka

import (
	"encoding/json"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func handleSigterm(eurekaServerUrl, appName, ipAddress string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		UnRegister(eurekaServerUrl, appName, ipAddress)
		os.Exit(1)
	}()
}

func startHeartbeat(eurekaServerUrl, appName, ipAddress string) {
	for {
		time.Sleep(time.Second * 30)
		heartbeat(eurekaServerUrl, appName, ipAddress)
	}
}

func heartbeat(eurekaServerUrl, appName, ipAddress string) {
	heartbeatAction := HttpAction{
		Url:         eurekaServerUrl + "/eureka/apps/" + appName + "/" + ipAddress,
		Method:      "PUT",
		ContentType: "application/json;charset=UTF-8",
	}
	doHttpRequest(heartbeatAction)
}

func Register(eurekaServerUrl, appName, port, securePort, instanceId, ipAddress string) {
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
			handleSigterm(eurekaServerUrl, appName, ipAddress)
			go startHeartbeat(eurekaServerUrl, appName, ipAddress)
			break
		} else {
			time.Sleep(time.Second * 5)
		}
	}
}

func UnRegister(eurekaServerUrl, appName, ipAddress string) {
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

func GetServiceInstances(eurekaServerUrl, appName string) ([]EurekaInstance, error) {
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
