/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017-2020 Philipp Mieden <dreadl0ck [at] protonicmp [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package maltego

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogo/protobuf/proto"

	"github.com/dreadl0ck/netcap/defaults"
	"github.com/dreadl0ck/netcap/types"
)

// ICMPv4CountFunc is a function that counts something over multiple ICMPv4 audit records.
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type ICMPv4CountFunc func()

// ICMPv4TransformationFunc is a transformation over ICMPv4 audit records.
//goland:noinspection GoUnnecessarilyExportedIdentifiers
type ICMPv4TransformationFunc = func(lt LocalTransform, trx *Transform, icmp *types.ICMPv4, min, max uint64, path string, ip string)

// ICMPv4Transform applies a maltego transformation over ICMPv4 audit records.
func ICMPv4Transform(count ICMPv4CountFunc, transform ICMPv4TransformationFunc) {
	var (
		lt               = ParseLocalArguments(os.Args[1:])
		path             = lt.Values["path"]
		ipaddr           = lt.Values["ipaddr"]
		dir              = filepath.Dir(path)
		icmpAuditRecords = filepath.Join(dir, "ICMPv4.ncap.gz")
		trx              = Transform{}
	)

	f, path := openFile(icmpAuditRecords)

	// check if its an audit record file
	if !strings.HasSuffix(f.Name(), defaults.FileExtensionCompressed) && !strings.HasSuffix(f.Name(), defaults.FileExtension) {
		die(errUnexpectedFileType, f.Name())
	}

	r := openNetcapArchive(path)

	// read netcap header
	header, errFileHeader := r.ReadHeader()
	if errFileHeader != nil {
		die("failed to read file header", errFileHeader.Error())
	}

	if header.Type != types.Type_NC_ICMPv4 {
		die("file does not contain ICMPv4 records", header.Type.String())
	}

	var (
		icmp = new(types.ICMPv4)
		pm   proto.Message
		ok   bool
	)
	pm = icmp

	if _, ok = pm.(types.AuditRecord); !ok {
		panic("type does not implement types.AuditRecord interface")
	}

	var (
		min uint64 = 10000000
		max uint64 = 0
		err error
	)

	if count != nil {
		for {
			err = r.Next(icmp)
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			} else if err != nil {
				die(err.Error(), errUnexpectedReadFailure)
			}

			count()
		}

		err = r.Close()
		if err != nil {
			log.Println("failed to close audit record file: ", err)
		}
	}

	r = openNetcapArchive(path)

	// read netcap header - ignore err as it has been checked before
	_, _ = r.ReadHeader()

	for {
		err = r.Next(icmp)
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		} else if err != nil {
			panic(err)
		}

		transform(lt, &trx, icmp, min, max, path, ipaddr)
	}

	err = r.Close()
	if err != nil {
		log.Println("failed to close audit record file: ", err)
	}

	trx.AddUIMessage("completed!", UIMessageInform)
	fmt.Println(trx.ReturnOutput())
}
