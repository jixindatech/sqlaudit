package monitor

import (
	"github.com/jixindatech/sqlaudit/audit/server"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"sync"
)

type Prometheus struct {
	addr string
	svr  *server.Server
	host string

	data sync.Map
}

//新建prometheus实例
func NewPrometheus(addr string, svr *server.Server) (*Prometheus, error) {
	prometheus := new(Prometheus)
	prometheus.addr = addr
	prometheus.svr = svr
	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	prometheus.host = hostName

	golog.Info("prometheus", zap.String("new", addr))

	return prometheus, nil
}

//启动prometheus的http监控
func (p *Prometheus) Run() {
	//workerDB := NewSqlAuditMetric("db")
	sqlaudit := NewSqlAuditMetric("sqlaudit")
	sqlaudit.hostname = p.host
	sqlaudit.svr = p.svr

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(sqlaudit)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	http.ListenAndServe(":9898", nil)
}

type SqlAuditMetric struct {
	Zone     string
	hostname string
	svr      *server.Server

	ClientCountDesc     *prometheus.Desc
	PacketPerClientDesc *prometheus.Desc
}

// Simulate prepare the data
func (c *SqlAuditMetric) GetMetricFromServer() (
	clientCount map[string]int, packetPerClient map[string]float64,
) {
	count, items := c.svr.GetMonitorData()
	// Just example fake data.
	clientCount = map[string]int{
		c.hostname: count,
	}
	packetPerClient = items
	return
}

// Describe simply sends the two Descs in the struct to the channel.
func (c *SqlAuditMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.ClientCountDesc
	ch <- c.PacketPerClientDesc
}

func (c *SqlAuditMetric) Collect(ch chan<- prometheus.Metric) {
	clientCount, packetPerClient := c.GetMetricFromServer()
	for host, oomCount := range clientCount {
		ch <- prometheus.MustNewConstMetric(
			c.ClientCountDesc,
			prometheus.GaugeValue,
			float64(oomCount),
			host,
		)
	}
	for host, ramUsage := range packetPerClient {
		ch <- prometheus.MustNewConstMetric(
			c.PacketPerClientDesc,
			prometheus.GaugeValue,
			ramUsage,
			host,
		)
	}
}

func NewSqlAuditMetric(zone string) *SqlAuditMetric {
	return &SqlAuditMetric{
		Zone: zone,
		ClientCountDesc: prometheus.NewDesc(
			"sqlaudit_client_number",
			"Number of sqlaudit clients.",
			[]string{"host"},
			prometheus.Labels{"zone": zone},
		),
		PacketPerClientDesc: prometheus.NewDesc(
			"sqlaudit_client_recevied_packets",
			"Number of per client received packets.",
			[]string{"client"},
			prometheus.Labels{"zone": zone},
		),
	}
}
