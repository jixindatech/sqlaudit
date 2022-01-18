package tcpreassembly

import (
	"encoding/hex"
	"fmt"
	"github.com/go-basic/uuid"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/reassembly"
	"github.com/jixindatech/sqlaudit/pkg/apps/mysql"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/jixindatech/sqlaudit/pkg/queue"
	"github.com/jixindatech/sqlaudit/pkg/storage"
	"go.uber.org/zap"
	"sync"
)

var stats struct {
	ipdefrag            int
	missedBytes         int
	pkt                 int
	sz                  int
	totalsz             int
	rejectFsm           int
	rejectOpt           int
	rejectConnFsm       int
	reassembled         int
	outOfOrderBytes     int
	outOfOrderPackets   int
	biggestChunkBytes   int
	biggestChunkPackets int
	overlapBytes        int
	overlapPackets      int
}

/*
 * The assembler context
 */
type Context struct {
	SqlType     int
	CaptureInfo gopacket.CaptureInfo
	Storage     *storage.Storage
}

func (c *Context) GetCaptureInfo() gopacket.CaptureInfo {
	return c.CaptureInfo
}

/*
 * The TCP factory: returns a new Stream
 */
type TcpStreamFactory struct {
	Queue queue.Queue
}

func (factory *TcpStreamFactory) New(net, transport gopacket.Flow, tcp *layers.TCP, ac reassembly.AssemblerContext) reassembly.Stream {
	golog.Debug("tcp", zap.String("new", fmt.Sprintf("%s %s", net.String(), transport.String())))
	fsmOptions := reassembly.TCPSimpleFSMOptions{
		SupportMissingEstablishment: false,
	}
	context := ac.(*Context)
	stream := &TcpStream{
		net:        net,
		transport:  transport,
		sqlType:    context.SqlType,
		tcpstate:   reassembly.NewTCPSimpleFSM(fsmOptions),
		ident:      fmt.Sprintf("%s:%s", net, transport),
		optchecker: reassembly.NewTCPOptionCheck(),
	}

	if stream.sqlType == config.SQL_TYPE_MYSQL {
		info := new(mysql.MysqlInfo)
		info.Transaction = uuid.New()
		info.Src = net.Src().String()
		info.Dst = net.Dst().String()

		info.Queue = factory.Queue

		stream.session = info
	}

	return stream
}

/*
 * TCP stream
 */

/* It's a connection (bidirectional) */
type TcpStream struct {
	tcpstate       *reassembly.TCPSimpleFSM
	fsmerr         bool
	optchecker     reassembly.TCPOptionCheck
	net, transport gopacket.Flow
	sqlType        int
	reversed       bool
	session        interface{}
	urls           []string
	ident          string
	sync.Mutex
}

func (t *TcpStream) Accept(tcp *layers.TCP, ci gopacket.CaptureInfo, dir reassembly.TCPFlowDirection, nextSeq reassembly.Sequence, start *bool, ac reassembly.AssemblerContext) bool {
	// FSM
	if !t.tcpstate.CheckState(tcp, dir) {
		// golog.Error("tcp", zap.String("FSM", fmt.Sprintf("%s: Packet rejected by FSM (state:%s)", t.ident, t.tcpstate.String())))
		stats.rejectFsm++
		if !t.fsmerr {
			t.fsmerr = true
			stats.rejectConnFsm++
		}
	}
	// Options
	err := t.optchecker.Accept(tcp, ci, dir, nextSeq, start)
	if err != nil {
		golog.Error("tcp", zap.String("OptionChecker", fmt.Sprintf("%s: Packet rejected by OptionChecker: %s\n", t.ident, err)))
		stats.rejectOpt++
		if false {
			return false
		}
	}
	// Checksum
	accept := true
	if false {
		c, err := tcp.ComputeChecksum()
		if err != nil {
			golog.Error("tcp", zap.String("ChecksumCompute", fmt.Sprintf("%s: Got error computing checksum: %s", t.ident, err)))
			accept = false
		} else if c != 0x0 {
			golog.Error("tcp", zap.String("Checksum", fmt.Sprintf("%s: Invalid checksum: 0x%x", t.ident, c)))
			accept = false
		}
	}
	if !accept {
		stats.rejectOpt++
	}
	return accept
}

func (t *TcpStream) ReassembledSG(sg reassembly.ScatterGather, ac reassembly.AssemblerContext) {
	// dir, start, end, skip := sg.Info()
	dir, _, _, skip := sg.Info()
	length, saved := sg.Lengths()
	// update stats
	sgStats := sg.Stats()

	if skip > 0 {
		stats.missedBytes += skip
	}
	stats.sz += length - saved
	stats.pkt += sgStats.Packets
	if sgStats.Chunks > 1 {
		stats.reassembled++
	}
	stats.outOfOrderPackets += sgStats.QueuedPackets
	stats.outOfOrderBytes += sgStats.QueuedBytes
	if length > stats.biggestChunkBytes {
		stats.biggestChunkBytes = length
	}
	if sgStats.Packets > stats.biggestChunkPackets {
		stats.biggestChunkPackets = sgStats.Packets
	}
	if sgStats.OverlapBytes != 0 && sgStats.OverlapPackets == 0 {
		golog.Error("tcp", zap.String("sgstats", fmt.Sprintf("bytes:%d, pkts:%d", sgStats.OverlapBytes, sgStats.OverlapPackets)))
		golog.Fatal("tcp", zap.String("sgstats", "invalid overlap"))

	}
	stats.overlapBytes += sgStats.OverlapBytes
	stats.overlapPackets += sgStats.OverlapPackets
	/*
		var ident string
		if dir == reassembly.TCPDirClientToServer {
			ident = fmt.Sprintf("%v %v(%s): ", t.net, t.transport, dir)
		} else {
			ident = fmt.Sprintf("%v %v(%s): ", t.net.Reverse(), t.transport.Reverse(), dir)
		}
		fmt.Printf("%s: SG reassembled packet with %d bytes (start:%v,end:%v,skip:%d,saved:%d,nb:%d,%d,overlap:%d,%d)\n",
			ident, length, start, end, skip, saved, sgStats.Packets, sgStats.Chunks, sgStats.OverlapBytes, sgStats.OverlapPackets)
	*/
	if skip == -1 /*&& *allowmissinginit*/ {
		// this is allowed
	} else if skip != 0 {
		// Missing bytes in stream: do not even try to parse it
		return
	}
	data := sg.Fetch(length)
	if t.sqlType == config.SQL_TYPE_MYSQL {
		if length > 0 {
			if false {
				fmt.Printf("Feeding mysql with:\n%s", hex.Dump(data))
			}

			if dir == reassembly.TCPDirClientToServer && !t.reversed {
				golog.Debug("ProcessClient")
				err := mysql.ProcessClient(t.session, data)
				if err != nil {
					golog.Error("tcp", zap.String("mysql", fmt.Sprintf("mysql process client error:", err)))
				}
			} else {
				golog.Debug("ProcessServer")
				err := mysql.ProcessServer(t.session, data)
				if err != nil {
					golog.Error("tcp", zap.String("mysql", fmt.Sprintf("mysql process server error:", err)))
				}
			}
		}
	}
}

func (t *TcpStream) ReassemblyComplete(ac reassembly.AssemblerContext) bool {
	golog.Debug("tcp", zap.String("connection", fmt.Sprintf("%s: Connection closed", t.ident)))
	// do not remove the connection to allow last ACK
	return false
}
