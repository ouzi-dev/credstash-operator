package env

import (
	"fmt"
	"os"
)

const ServiceMonitorNamespaceEnvVar = "SERVICE_MONITOR_NAMESPACE"

func GetServiceMonitorNamespace() (string, error) {
	ns, found := os.LookupEnv(ServiceMonitorNamespaceEnvVar)
	if !found {
		return "", fmt.Errorf("%s must be set", ServiceMonitorNamespaceEnvVar)
	}

	return ns, nil
}
