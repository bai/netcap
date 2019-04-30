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

package types

import "strconv"

func (f TransportFlow) CSVHeader() []string {
	return filter([]string{
		"TimestampFirst",
		"TimestampLast",
		"Proto",
		"SrcPort",
		"DstPort",
		"TotalSize",
		"NumPackets",
		"UID",
		"Duration",
	})
}

func (f TransportFlow) CSVRecord() []string {
	return filter([]string{
		formatTimestamp(f.TimestampFirst),
		formatTimestamp(f.TimestampLast),
		f.Proto,
		formatInt32(f.SrcPort),
		formatInt32(f.DstPort),
		formatInt64(f.TotalSize),
		formatInt64(f.NumPackets),
		strconv.FormatUint(f.UID, 10),
		formatInt64(f.Duration),
	})
}

func (f TransportFlow) NetcapTimestamp() string {
	return f.TimestampFirst
}

func (u TransportFlow) JSON() (string, error) {
	return jsonMarshaler.MarshalToString(&u)
}
