package names

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Separator is the separator used in resource names.
const Separator = ":"

func formatPort(port uint32) string {
	return strconv.FormatUint(uint64(port), 10)
}

// Join uses Separator to join the given parts into a resource name.
func Join(parts ...string) string {
	return strings.Join(parts, Separator)
}

func GetLocalClusterName(port uint32) string {
	return Join("localhost", formatPort(port))
}

func GetSplitClusterName(service string, idx int) string {
	return fmt.Sprintf("%s-_%d_", service, idx)
}

func GetPortForLocalClusterName(cluster string) (uint32, error) {
	parts := strings.Split(cluster, Separator)
	if len(parts) != 2 {
		return 0, errors.Errorf("failed to  parse local cluster name: %s", cluster)
	}
	port, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(port), nil
}

func GetInboundListenerName(address string, port uint32) string {
	return Join("inbound",
		net.JoinHostPort(address, formatPort(port)))
}

func GetOutboundListenerName(address string, port uint32) string {
	return Join("outbound",
		net.JoinHostPort(address, formatPort(port)))
}

func GetInboundRouteName(service string) string {
	return Join("inbound", service)
}

func GetOutboundRouteName(service string) string {
	return Join("outbound", service)
}

func GetEnvoyAdminClusterName() string {
	return Join("kuma", "envoy", "admin")
}

func GetMetricsHijackerClusterName() string {
	return Join("kuma", "metrics", "hijacker")
}

func GetPrometheusListenerName() string {
	return Join("kuma", "metrics", "prometheus")
}

func GetAdminListenerName() string {
	return Join("kuma", "envoy", "admin")
}

func GetTracingClusterName(backendName string) string {
	return Join("tracing", backendName)
}

func GetDNSListenerName() string {
	return Join("kuma", "dns")
}

func GetGatewayListenerName(gatewayName string, protoName string, port uint32) string {
	return Join(gatewayName, protoName, formatPort(port))
}

// GetSecretName constructs a secret name that has a good chance of being
// unique across subsystems that are unaware of each other.
//
// category should be used to indicate the type of the secret resource. For
// example, is this a TLS certificate, or a ValidationContext, or something else.
//
// scope is a qualifier within which identifier can be considered to be unique.
// For example, the name of a Kuma file DataSource is unique across file
// DataSources, but may collide with the name of a secret DataSource.
//
// identifier is a name that should be unique within a category and scope.
func GetSecretName(category string, scope string, identifier string) string {
	return Join(category, scope, identifier)
}
