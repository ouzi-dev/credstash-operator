package env

import (
	"fmt"
	"os"
)

const ServiceMonitorNamespaceEnvVar = "SERVICE_MONITOR_NAMESPACE"
const OperatorPodNamespaceEnvVar = "POD_NAMESPACE"

func GetServiceMonitorNamespace() (string, error) {
	ns, found := os.LookupEnv(ServiceMonitorNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", ServiceMonitorNamespaceEnvVar)
	}

	return ns, nil
}

func GetOperatorPodNamespace() (string, error) {
	ns, found := os.LookupEnv(OperatorPodNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", OperatorPodNamespaceEnvVar)
	}

	return ns, nil
}