package device

import (
	"fmt"
	"strings"
)

type Generic struct {
	DeviceId    uint32
	DeviceIp    string
	DeviceType  Type
	DeviceModel string
	Name        string
	support     map[string]bool
}

func NewGeneric(id uint32, model string, name string, support string) *Generic {
	g := &Generic{
		DeviceId:    id,
		DeviceType:  CheckDevice(model),
		DeviceModel: model,
		Name:        name,
	}
	g.Support(support)
	return g
}

func (g *Generic) Type() Type {
	return g.DeviceType
}

func (g *Generic) Id() uint32 {
	return g.DeviceId
}

func (g *Generic) Start(ip string, update chan Device) error {
	g.DeviceIp = ip
	return fmt.Errorf("pseudo device can't start")
}
func (g *Generic) Close() error {
	return fmt.Errorf("pseudo device can't stop")
}

func (g *Generic) Cmd(cmd string, v1 string, v2 string, v3 string) ([]string, error) {
	return nil, fmt.Errorf("can't process commands")
}

func (g *Generic) Support(v string) {
	g.support = make(map[string]bool)
	sp := strings.Split(v, " ")
	for _, t := range sp {
		g.support[t] = true
	}
}

func (g *Generic) String() string {
	return fmt.Sprintf(`{"id":"%s","ip":"%s","model":"%s","name":"%s"}`,
		g.DeviceId, g.DeviceIp, g.DeviceModel, g.Name)
}
