package sessdata

import (
	"fmt"
	"net"
	"os"
	"path"
	"testing"

	"github.com/bjdgyc/anylink/base"
	"github.com/bjdgyc/anylink/dbdata"
	"github.com/stretchr/testify/assert"
)

func preIpData() {
	base.Cfg.Ipv4Network = "192.168.3.0"
	base.Cfg.Ipv4Netmask = "255.255.255.0"
	base.Cfg.Ipv4Pool = []string{"192.168.3.1", "192.168.3.199"}
	tmpDb := path.Join(os.TempDir(), "anylink_test.db")
	base.Cfg.DbFile = tmpDb
	dbdata.Start()
}

func closeIpdata() {
	dbdata.Stop()
	tmpDb := path.Join(os.TempDir(), "anylink_test.db")
	os.Remove(tmpDb)
}

func TestIpPool(t *testing.T) {
	assert := assert.New(t)
	preIpData()
	defer closeIpdata()

	initIpPool()

	var ip net.IP

	for i := 1; i <= 100; i++ {
		ip = AcquireIp("user", fmt.Sprintf("mac-%d", i))
	}
	ip = AcquireIp("user", fmt.Sprintf("mac-new"))
	assert.True(net.IPv4(192, 168, 3, 101).Equal(ip))
	for i := 102; i <= 199; i++ {
		ip = AcquireIp("user", fmt.Sprintf("mac-%d", i))
	}
	assert.True(net.IPv4(192, 168, 3, 199).Equal(ip))
	ip = AcquireIp("user", fmt.Sprintf("mac-nil"))
	assert.Nil(ip)

	ReleaseIp(net.IPv4(192, 168, 3, 88), "mac-88")
	ReleaseIp(net.IPv4(192, 168, 3, 77), "mac-77")
	// 最早过期的ip
	ip = AcquireIp("user", "mac-release-new")
	assert.True(net.IPv4(192, 168, 3, 88).Equal(ip))
}
