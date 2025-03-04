package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func TestAddNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory()
	assert.NotNil(f)

	f.Add("Switch", func() zebra.Resource { return new(network.Switch) })
	assert.NotNil(f.New("Switch"))
	assert.Nil(f.New("random"))
}

func TestNewResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(zebra.NewResourceList(nil))
}

func TestCopyResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resA := zebra.NewResourceList(nil)
	assert.NotNil(resA)

	resA.Resources = append(resA.Resources, new(network.IPAddressPool))

	resB := zebra.NewResourceList(nil)
	assert.NotNil(resB)
	assert.True(len(resB.Resources) == 0)

	zebra.CopyResourceList(resB, resA)
	assert.True(len(resB.Resources) == 1)
}

func TestListMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	resA := zebra.NewResourceList(funMap)
	assert.NotNil(resA)

	vlan := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "invalid",
			Labels: nil,
		},
		RangeStart: 0,
		RangeEnd:   10,
	}

	resA.Resources = append(resA.Resources, vlan)

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceList(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.NotNil(err)

	vlan.Type = "VLANPool"
	resA.Resources = []zebra.Resource{vlan}

	bytes, err = resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB = zebra.NewResourceList(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
	assert.True(len(resB.Resources) == 1)
}

func TestErrorMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	resList := zebra.NewResourceList(funMap)
	assert.NotNil(resList.UnmarshalJSON(nil))
	assert.NotNil(resList.UnmarshalJSON([]byte(`[{"id":"0100000001"}]`)))
	assert.NotNil(resList.UnmarshalJSON([]byte(`[{"id":"0100000001", "type":123}]`)))

	resMap := zebra.NewResourceMap(nil)
	assert.NotNil(resMap.UnmarshalJSON(nil))
	assert.NotNil(resMap.UnmarshalJSON([]byte(`{"VLANPool":[{"id":"0100000001", "type":123}]}`)))
}

func TestNewResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(zebra.NewResourceMap(nil))
}

func TestCopyResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resA := zebra.NewResourceMap(nil)
	assert.NotNil(resA)

	resA.Add(new(network.IPAddressPool), "IPAddressPool")

	resB := zebra.NewResourceMap(nil)
	assert.NotNil(resB)

	zebra.CopyResourceMap(resB, resA)
	assert.True(len(resB.Resources) == 1)
	assert.True(len(resB.Resources["IPAddressPool"].Resources) == 1)
}

func TestGetFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory()
	f.Add("Switch", func() zebra.Resource { return new(network.Switch) })

	resA := zebra.NewResourceMap(f)
	assert.NotNil(resA)

	assert.NotNil(resA.GetFactory())
}

func TestAdd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add("Switch", func() zebra.Resource { return new(network.Switch) })

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	switch1 := funMap.New("Switch")

	resA.Add(switch1, "Switch")
	assert.NotNil(len(resA.Resources["Switch"].Resources) == 1)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add("Switch", func() zebra.Resource { return new(network.Switch) })

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	switch1 := funMap.New("Switch")

	resA.Add(switch1, "Switch")
	assert.NotNil(len(resA.Resources["Switch"].Resources) == 1)

	resA.Delete(switch1, "Switch")
	assert.NotNil(len(resA.Resources["Switch"].Resources) == 0)
}

func TestMapMarshalUnMarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	vlan := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 0,
		RangeEnd:   10,
	}

	resA.Add(vlan, "VLANPool")

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceMap(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
}
