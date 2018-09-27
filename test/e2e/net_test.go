package e2e

import (
	"fmt"
	"testing"
)

func TestNetTCPandHTTP(t *testing.T) {
	env := BuildEnv(t)
	logger := Logger{}
	kwt := Kwt{t, env.Namespace, Logger{}}
	kubectl := Kubectl{t, env.Namespace, Logger{}}
	kwtNet := NewKwtNet(kwt, t, Logger{})

	kwtNet.Start([]string{})
	defer kwtNet.End()

	guestbookAddrs := Guestbook{kwt, kubectl, t, logger}.Install()
	netProbe := NetworkProbe{t, logger}

	for _, url := range []string{
		fmt.Sprintf("http://%s", guestbookAddrs.FrontendSvcIP),
		fmt.Sprintf("http://%s", guestbookAddrs.FrontendSvcDomain),
	} {
		netProbe.HTTPGet(url, "Guestbook", "guestbook")
	}

	for i, addr := range []string{guestbookAddrs.RedisSvcIP, guestbookAddrs.RedisSvcDomain} {
		netProbe.RedisWriteRead(addr, fmt.Sprintf("value%d", i))
	}
}
