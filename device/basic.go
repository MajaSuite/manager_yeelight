package device

import "fmt"

type Basic struct {
	DeviceType Type
	DeviceIp   string
	DeviceId   uint32
}

func NewBasic(id uint32, ip string) Device {
	return &Basic{
		DeviceType: NO_TYPE,
		DeviceId:   id,
		DeviceIp:   ip,
	}
}

func (b *Basic) Type() Type {
	return b.DeviceType
}

func (b *Basic) ID() uint32 {
	return b.DeviceId
}

func (b *Basic) IP() string {
	return b.DeviceIp
}

func (b *Basic) Token() []byte {
	return nil
}

func (b *Basic) Model() string {
	return ""
}

func (b *Basic) Name() string {
	return "Basic device"
}

func (b *Basic) Close() error {
	return ErrInvalidCommand
}

func (b *Basic) String() string {
	return fmt.Sprintf("{type:%s, id:%x, ip:%s}", b.DeviceType.String(), b.DeviceId, b.DeviceIp)
}

func Close() error {
	return nil
}

func (b *Basic) SetIP(ip string) error {
	b.DeviceIp = ip

	return nil
}

func (b *Basic) SetName(name string) error {
	return ErrInvalidCommand
}
