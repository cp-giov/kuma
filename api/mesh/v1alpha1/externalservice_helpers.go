package v1alpha1

import (
	"net"
	"strconv"

	util_proto "github.com/kumahq/kuma/pkg/util/proto"
)

func (m *ExternalService) UnmarshalJSON(data []byte) error {
	return util_proto.FromJSON(data, m)
}

func (m *ExternalService) MarshalJSON() ([]byte, error) {
	return util_proto.ToJSON(m)
}
func (t *ExternalService) DeepCopyInto(out *ExternalService) {
	util_proto.Merge(out, t)
}
func (t *ExternalService) DeepCopy() *ExternalService {
	if t == nil {
		return nil
	}
	out := new(ExternalService)
	t.DeepCopyInto(out)
	return out
}

// Matches is simply an alias for MatchTags to make source code more aesthetic.
func (es *ExternalService) Matches(selector TagSelector) bool {
	if es != nil {
		return es.MatchTags(selector)
	}
	return false
}

func (es *ExternalService) MatchTags(selector TagSelector) bool {
	return selector.Matches(es.Tags)
}

func (es *ExternalService) GetService() string {
	if es == nil {
		return ""
	}
	return es.Tags[ServiceTag]
}

func (es *ExternalService) GetProtocol() string {
	if es == nil {
		return ""
	}
	return es.Tags[ProtocolTag]
}

func (es *ExternalService) GetHost() string {
	if es == nil {
		return ""
	}
	host, _, err := net.SplitHostPort(es.Networking.Address)
	if err != nil {
		return ""
	}
	return host
}

func (es *ExternalService) GetPort() string {
	if es == nil {
		return ""
	}
	_, port, err := net.SplitHostPort(es.Networking.Address)
	if err != nil {
		return ""
	}
	return port
}

func (es *ExternalService) GetPortUInt32() uint32 {
	port := es.GetPort()
	iport, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}
	return uint32(iport)
}

func (es *ExternalService) TagSet() SingleValueTagSet {
	return es.Tags
}
