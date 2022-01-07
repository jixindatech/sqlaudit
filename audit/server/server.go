package server

import (
	"bufio"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/reassembly"
	"github.com/jixindatech/sqlaudit/pkg/apps/mysql"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/jixindatech/sqlaudit/pkg/queue"
	"github.com/jixindatech/sqlaudit/pkg/storage"
	"github.com/jixindatech/sqlaudit/pkg/tcpreassembly"
	"go.uber.org/zap"
	"io"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Offline = iota
	Online
	Unknown
)

const timeout time.Duration = time.Minute * 6
const closeTimeout time.Duration = time.Minute * 6

type Server struct {
	cfg  *config.Config
	addr string

	statusIndex int32
	status      [2]int32
	/*
		logSqlIndex        int32
		logSql             [2]string
		slowLogTimeIndex   int32
		slowLogTime        [2]int
		blacklistSqlsIndex int32
	*/
	listener          net.Listener
	handle            *pcap.Handle
	capture           bool
	captureSqlType    int
	running           bool
	Queue             queue.Queue
	Storage           storage.Storage
	configUpdateMutex sync.RWMutex
	configVer         uint32

	clients      map[string]int
	monitorMutex sync.RWMutex
}

func (s *Server) Status() string {
	var status string
	switch s.status[s.statusIndex] {
	case Online:
		status = "online"
	case Offline:
		status = "offline"
	case Unknown:
		status = "unknown"
	default:
		status = "unknown"
	}
	return status
}

func NewServer(cfg *config.Config, capture bool, sqltype int, device string, queueInstance queue.Queue) (*Server, error) {
	s := new(Server)

	s.cfg = cfg
	s.addr = cfg.Addr
	s.Queue = queueInstance

	atomic.StoreInt32(&s.statusIndex, 0)
	s.status[s.statusIndex] = Online

	s.clients = make(map[string]int)

	var err error
	netProto := "tcp"

	s.Storage, err = storage.GetStorage(cfg.EsConfig)
	if err != nil {
		return nil, err
	}

	err = mysql.ParserSqlRules()
	if err != nil {
		return nil, err
	}

	if capture {
		s.capture = capture
		s.captureSqlType = sqltype
		const snapshot_len = 65536
		promiscuous := true
		s.handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, pcap.BlockForever)
		if err != nil {
			return nil, err
		}

		var filter string = "tcp and port 3306"
		err = s.handle.SetBPFFilter(filter)
		if err != nil {
			return nil, err
		}

	} else {
		s.listener, err = net.Listen(netProto, s.addr)

		if err != nil {
			return nil, err
		}
	}

	golog.Info("server", zap.String("new", fmt.Sprintf("%s:%s", netProto, s.addr)))
	return s, nil
}

func (s *Server) Run() {
	s.running = true

	if s.capture {
		s.onCapture()
	} else {
		for s.running {
			conn, err := s.listener.Accept()
			if err != nil {
				golog.Error("server", zap.String("listen", err.Error()))
				continue
			}

			s.monitorMutex.Lock()
			s.clients[conn.RemoteAddr().String()]++
			s.monitorMutex.Unlock()

			go s.onConn(conn)
		}
	}
}

func (s *Server) processPacket(sqlType int, packet gopacket.Packet, deFrag *ip4defrag.IPv4Defragmenter, assembler *reassembly.Assembler) {
	// Process packet here
	if true {
		ip4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ip4Layer == nil {
			return
		}
		ip4 := ip4Layer.(*layers.IPv4)
		l := ip4.Length
		newip4, err := deFrag.DefragIPv4(ip4)
		if err != nil {
			golog.Fatal("server", zap.String("defrag", err.Error()))
		} else if newip4 == nil {
			return
		}
		if newip4.Length != l {
			pb, ok := packet.(gopacket.PacketBuilder)
			if !ok {
				golog.Fatal("server", zap.String("defrag", "packetBuilder failed"))
			}
			nextDecoder := newip4.NextLayerType()
			nextDecoder.Decode(newip4.Payload, pb)
		}
	}

	tcp := packet.Layer(layers.LayerTypeTCP)
	if tcp != nil {
		tcp := tcp.(*layers.TCP)
		if false {
			err := tcp.SetNetworkLayerForChecksum(packet.NetworkLayer())
			if err != nil {
				golog.Error("tcp", zap.String("checksum", err.Error()))
			}
		}
		c := tcpreassembly.Context{
			SqlType:     sqlType,
			CaptureInfo: packet.Metadata().CaptureInfo,
		}
		assembler.AssembleWithContext(packet.NetworkLayer().NetworkFlow(), tcp, &c)
	}

	return
}
func (s *Server) onCapture() {
	nextFlush := time.Now().Add(timeout)

	deFrag := ip4defrag.NewIPv4Defragmenter()
	streamFactory := &tcpreassembly.TcpStreamFactory{
		Queue: s.Queue,
	}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	packetSource := gopacket.NewPacketSource(s.handle, s.handle.LinkType())
	for packet := range packetSource.Packets() {
		s.processPacket(s.captureSqlType, packet, deFrag, assembler)

		now := time.Now()
		if now.After(nextFlush) {
			_, _ = assembler.FlushWithOptions(reassembly.FlushOptions{T: now.Add(-timeout), TC: now.Add(-closeTimeout)})
			nextFlush = now.Add(timeout)
		}
	}
}

func getSqlType(data []byte) (int, error) {
	for i := 0; i < len(data)/4; i++ {
		tmp := data[i*4 : i*4+4]
		key := int(uint32(tmp[1]) | uint32(tmp[0])<<8)
		value := int(uint32(tmp[3]) | uint32(tmp[2])<<8)
		if key == config.SQL_CLASS && value == config.SQL_TYPE_MYSQL {
			return value, nil
		}
	}

	return 0, fmt.Errorf("%s", "sqltype not found")
}

func (s *Server) onConn(c net.Conn) {
	defer c.Close()
	defer func() {
		err := recover()
		if err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			golog.Error("server", zap.String("onConn", fmt.Sprintf("remote_addr:%s,stack:%s", c.RemoteAddr().String(), string(buf))))
		}
	}()

	nextFlush := time.Now().Add(timeout)

	deFrag := ip4defrag.NewIPv4Defragmenter()
	streamFactory := &tcpreassembly.TcpStreamFactory{
		Queue: s.Queue,
	}
	streamPool := reassembly.NewStreamPool(streamFactory)
	assembler := reassembly.NewAssembler(streamPool)

	buf := bufio.NewReaderSize(c, 8*1024)
	header := []byte{0, 0, 0, 0}

	if _, err := io.ReadFull(buf, header); err != nil {
		golog.Error("server", zap.String("io", err.Error()), zap.String("ip", c.RemoteAddr().String()))
		return
	}

	length := int(uint32(header[3]) | uint32(header[2])<<8 | uint32(header[1])<<16 | uint32(header[0])<<24)
	data := make([]byte, length)
	if _, err := io.ReadFull(buf, data); err != nil {
		golog.Error("server", zap.String("io", err.Error()), zap.String("ip", c.RemoteAddr().String()))
		return
	}

	sqltype, err := getSqlType(data)
	if err != nil {
		golog.Error("server", zap.String("err", err.Error()), zap.String("ip", c.RemoteAddr().String()))
		return
	}

	for {
		if _, err := io.ReadFull(buf, header); err != nil {
			golog.Error("server", zap.String("io", err.Error()), zap.String("ip", c.RemoteAddr().String()))
			break
		}

		fmt.Println(header[0], header[1], header[2], header[3])

		length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16 | uint32(header[3])<<24)
		data := make([]byte, length)
		if _, err := io.ReadFull(buf, data); err != nil {
			break
		}
		packet := gopacket.NewPacket(data[15:], layers.LinkTypeEthernet, gopacket.DecodeOptions{Lazy: true, NoCopy: true})
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			continue
		}

		_ = packet.Metadata().CaptureInfo.Timestamp.UnmarshalBinary(data[:15])

		s.monitorMutex.Lock()
		s.clients[c.RemoteAddr().String()]++
		s.monitorMutex.Unlock()

		s.processPacket(sqltype, packet, deFrag, assembler)
		now := time.Now()
		if now.After(nextFlush) {
			flushed, closed := assembler.FlushWithOptions(reassembly.FlushOptions{T: now.Add(-timeout), TC: now.Add(-closeTimeout)})
			nextFlush = now.Add(timeout)
			fmt.Printf("Forced flush: %d flushed, %d\n", flushed, closed)
		}
	}

	s.monitorMutex.Lock()
	delete(s.clients, c.RemoteAddr().String())
	s.monitorMutex.Unlock()

}

func (s *Server) GetMonitorData() (count int, res map[string]float64) {
	res = make(map[string]float64)

	s.monitorMutex.RLock()
	defer s.monitorMutex.RUnlock()

	for addr, data := range s.clients {
		if data != 0 {
			res[addr] = float64(data)
			s.clients[addr] = 0
			count++
		}
	}
	return count, res
}

func (s *Server) Close() {
	s.running = false
	if s.listener != nil {
		_ = s.listener.Close()
	}
}

func (s *Server) UpdateConfig(newCfg *config.Config) {

}
