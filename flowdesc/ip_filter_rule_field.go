package flowdesc

type IPFilterRuleFieldList []IPFilterRuleField

type IPFilterRuleField interface {
	Set(*IPFilterRule) error
}

type IPFilterAction struct {
	Action Action
}

func (i *IPFilterAction) Set(r *IPFilterRule) error {
	return r.SetAction(i.Action)
}

type IPFilterDirection struct {
	Direction Direction
}

func (i *IPFilterDirection) Set(r *IPFilterRule) error {
	return r.SetDirection(i.Direction)
}

type IPFilterProto struct {
	Proto uint8
}

func (i *IPFilterProto) Set(r *IPFilterRule) error {
	return r.SetProtocol(i.Proto)
}

type IPFilterSourceIP struct {
	Src string
}

func (i *IPFilterSourceIP) Set(r *IPFilterRule) error {
	return r.SetSourceIP(i.Src)
}

type IPFilterSourcePorts struct {
	Ports string
}

func (i *IPFilterSourcePorts) Set(r *IPFilterRule) error {
	return r.SetSourcePorts(i.Ports)
}

type IPFilterDestinationIP struct {
	Src string
}

func (i *IPFilterDestinationIP) Set(r *IPFilterRule) error {
	return r.SetDestinationIP(i.Src)
}

type IPFilterDestinationPorts struct {
	Ports string
}

func (i *IPFilterDestinationPorts) Set(r *IPFilterRule) error {
	return r.SetDestinationPorts(i.Ports)
}

func BuildIPFilterRuleFromField(cl IPFilterRuleFieldList) (*IPFilterRule, error) {
	rule := NewIPFilterRule()

	var err error
	for _, config := range cl {
		err = config.Set(rule)
		if err != nil {
			return nil, flowDescErrorf("build ip filter rule failed by %s", err)
		}
	}

	return rule, nil
}
