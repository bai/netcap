/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package export

import (
	"github.com/namsral/flag"
	"os"
)

func Flags() (flags []string) {
	fs.VisitAll(func(f *flag.Flag) {
		flags = append(flags, f.Name)
	})
	return
}

var (
	fs         = flag.NewFlagSetWithEnvPrefix(os.Args[0], "NC", flag.ExitOnError)
	flagMetricsAddress = fs.String("address", "127.0.0.1:7777", "set address for exposing metrics")
	flagDumpJSON       = fs.Bool("dumpJson", false, "dump as JSON")
	flagReplay         = fs.Bool("replay", false, "replay traffic (only works when exporting audit records directly!)")
	flagDir            = fs.String("dir", "", "path to directory with netcap audit records")
	flagInput          = fs.String("read", "", "read specified file, can either be a pcap or netcap audit record file")
	flagInterface      = fs.String("iface", "", "attach to network interface and capture in live mode")
	flagWorkers        = fs.Int("workers", 1000, "number of workers")
	flagPacketBuffer   = fs.Int("pbuf", 100, "set packet buffer size, for channels that feed data to workers")
	flagIngoreUnknown  = fs.Bool("ignore-unknown", false, "disable writing unknown packets into a pcap file")
	flagPromiscMode    = fs.Bool("promisc", true, "toggle promiscous mode for live capture")
	flagSnapLen        = fs.Int("snaplen", 1514, "configure snaplen for live capture from interface")
	flagVersion        = fs.Bool("version", false, "print netcap package version and exit")
	flagBaseLayer      = fs.String("base", "ethernet", "select base layer")
	flagDecodeOptions  = fs.String("opts", "lazy", "select decoding options")
	flagPayload        = fs.Bool("payload", false, "capture payload for supported layers")
	flagCompress       = fs.Bool("comp", true, "compress output with gzip")
	flagBuffer         = fs.Bool("buf", true, "buffer data in memory before writing to disk")
	flagOutDir         = fs.String("out", "", "specify output directory, will be created if it does not exist")
	flagBPF            = fs.String("bpf", "", "supply a BPF filter to use prior to processing packets with netcap")
	flagInclude        = fs.String("include", "", "include specific encoders")
	flagExclude        = fs.String("exclude", "LinkFlow,TransportFlow,NetworkFlow", "exclude specific encoders")
	flagMemProfile     = fs.Bool("memprof", false, "create memory profile")
	flagCSV            = fs.Bool("csv", false, "print output data as csv with header line")
	flagContext        = fs.Bool("context", true, "add packet flow context to selected audit records")
	flagMemBufferSize  = fs.Int("membuf-size", 1024*1024*10, "set size for membuf")
	flagListInterfaces = fs.Bool("interfaces", false, "list all visible network interfaces")
	flagReverseDNS     = fs.Bool("reverse-dns", false, "resolve ips to domains via the operating systems default dns resolver")
	flagLocalDNS       = fs.Bool("local-dns", false, "resolve DNS locally via hosts file in the database dir")
	flagMACDB          = fs.Bool("macDB", false, "use mac to vendor database for device profiling")
	flagJa3DB          = fs.Bool("ja3DB", false, "use ja3 database for device profiling")
	flagServiceDB      = fs.Bool("serviceDB", false, "use serviceDB for device profiling")
	flagGeolocationDB  = fs.Bool("geoDB", false, "use geolocation for device profiling")
	flagDPI            = fs.Bool("dpi", false, "use DPI for device profiling")
)
