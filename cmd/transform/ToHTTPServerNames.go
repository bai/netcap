package transform

import (
	"github.com/dreadl0ck/netcap/maltego"
	"github.com/dreadl0ck/netcap/types"
)

func toHTTPServerNames() {
	maltego.HTTPTransform(
		nil,
		func(lt maltego.LocalTransform, trx *maltego.Transform, http *types.HTTP, min, max uint64, profilesFile string, ipaddr string) {
			if http.SrcIP == ipaddr {
				if http.ServerName != "" {
					trx.AddEntity("netcap.ServerName", http.ServerName)
				}
			}
		},
		false,
	)
}
