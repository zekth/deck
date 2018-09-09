package kong

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpstreamsService(T *testing.T) {
	assert := assert.New(T)

	client, err := NewClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	upstream := &Upstream{
		Name: String("virtual-host1"),
	}

	createdUpstream, err := client.Upstreams.Create(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(createdUpstream)

	upstream, err = client.Upstreams.Get(defaultCtx, createdUpstream.ID)
	assert.Nil(err)
	assert.NotNil(upstream)

	upstream.Name = String("virtual-host2")
	upstream, err = client.Upstreams.Update(defaultCtx, upstream)
	assert.Nil(err)
	assert.NotNil(upstream)
	assert.Equal("virtual-host2", *upstream.Name)

	err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	assert.Nil(err)

	// PUT request is not yet supported
	// TODO uncomment this upstream entity is migrated over to new DAO

	// ID can be specified
	// id := uuid.NewV4().String()
	// upstream = &Upstream{
	// 	Name: String("key-auth"),
	// 	ID:   String(id),
	// }

	// createdUpstream, err = client.Upstreams.Create(defaultCtx, upstream)
	// assert.Nil(err)
	// assert.NotNil(createdUpstream)
	// assert.Equal(id, *createdUpstream.ID)

	// err = client.Upstreams.Delete(defaultCtx, createdUpstream.ID)
	// assert.Nil(err)
}

func TestUpstreamListEndpoint(T *testing.T) {
	assert := assert.New(T)

	client, err := NewClient(nil, nil)
	assert.Nil(err)
	assert.NotNil(client)

	// fixtures
	upstreams := []*Upstream{
		&Upstream{
			Name: String("vhost1.com"),
		},
		&Upstream{
			Name: String("vhost2.com"),
		},
		&Upstream{
			Name: String("vhost3.com"),
		},
	}

	// create fixturs
	for i := 0; i < len(upstreams); i++ {
		upstream, err := client.Upstreams.Create(defaultCtx, upstreams[i])
		assert.Nil(err)
		assert.NotNil(upstream)
		upstreams[i] = upstream
	}

	upstreamsFromKong, next, err := client.Upstreams.List(defaultCtx, nil)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(upstreamsFromKong)
	assert.Equal(3, len(upstreamsFromKong))

	// check if we see all upstreams
	assert.True(compareUpstreams(upstreams, upstreamsFromKong))

	// Test pagination
	upstreamsFromKong = []*Upstream{}

	// first page
	page1, next, err := client.Upstreams.List(defaultCtx, &ListOpt{Size: 1})
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page1)
	assert.Equal(1, len(page1))
	upstreamsFromKong = append(upstreamsFromKong, page1...)

	// second page
	page2, next, err := client.Upstreams.List(defaultCtx, next)
	assert.Nil(err)
	assert.NotNil(next)
	assert.NotNil(page2)
	assert.Equal(1, len(page2))
	upstreamsFromKong = append(upstreamsFromKong, page2...)

	// last page
	page3, next, err := client.Upstreams.List(defaultCtx, next)
	assert.Nil(err)
	assert.Nil(next)
	assert.NotNil(page3)
	assert.Equal(1, len(page3))
	upstreamsFromKong = append(upstreamsFromKong, page3...)

	assert.True(compareUpstreams(upstreams, upstreamsFromKong))

	for i := 0; i < len(upstreams); i++ {
		assert.Nil(client.Upstreams.Delete(defaultCtx, upstreams[i].ID))
	}
}

func compareUpstreams(expected, actual []*Upstream) bool {
	var expectedNames, actualNames []string
	for _, upstream := range expected {
		expectedNames = append(expectedNames, *upstream.Name)
	}

	for _, upstream := range actual {
		actualNames = append(actualNames, *upstream.Name)
	}

	return (compareSlices(expectedNames, actualNames))
}
