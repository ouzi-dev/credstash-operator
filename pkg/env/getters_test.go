package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServiceMonitorNamespace(t *testing.T) {
	err := os.Unsetenv(ServiceMonitorNamespaceEnvVar)
	assert.NoError(t, err)

	ns, err := GetServiceMonitorNamespace()
	assert.Error(t, err)
	assert.Equal(t, "", ns)

	err = os.Setenv(ServiceMonitorNamespaceEnvVar, "namespace")
	assert.NoError(t, err)

	ns, err = GetServiceMonitorNamespace()
	assert.NoError(t, err)
	assert.Equal(t, "namespace", ns)
}
