package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"net/http"

)

var test string
var gtest string

type mysqlUp struct {
	Up *prometheus.Desc
}

func NewMysqlUp() *mysqlUp {
	return &mysqlUp{
		Up: prometheus.NewDesc("mysql_up", "Is MySQL up", nil, nil),

	}
}

func (m *mysqlUp) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.Up
}
func (m *mysqlUp) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(m.Up, prometheus.GaugeValue, 1)
}
type nginxUp struct {
	Up *prometheus.Desc
	Connections *prometheus.Desc
}

func NewNginxUp() *nginxUp {
	return &nginxUp{
		Up: prometheus.NewDesc("nginx_up", "Is Nginx up", []string{"host"}, nil),
		Connections: prometheus.NewDesc("nginx_connections", "Nginx connections", nil, nil),
		
	}
}

func (n *nginxUp) Describe(ch chan<- *prometheus.Desc) {
	ch <- n.Up
	ch <- n.Connections
}
func (n *nginxUp) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(n.Up, prometheus.GaugeValue, 1,"127.0.0.1")
	ch <- prometheus.MustNewConstMetric(n.Connections, prometheus.GaugeValue, 100) // Example value, replace with actual logic
}


var rootcmd = &cobra.Command{
	Short: "myapp is a simple app",
	Long:  "This is a longer description of myapp",
	Run: func(cmd *cobra.Command, args []string) {
		prometheus.MustRegister(NewMysqlUp())
		prometheus.MustRegister(NewNginxUp())
		fmt.Println("Here", test, gtest)
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println("Error starting server:", err)
		}
	},
}


func init() {
	rootcmd.AddCommand(vercmd)
	rootcmd.Flags().StringVarP(&test, "test", "t", "", "test flag")
	rootcmd.PersistentFlags().StringVarP(&gtest, "gtest", "g", "", "global test flag")
}

var vercmd = &cobra.Command{
	Use:   "version",
	Short: "get version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("myapp version 1.0.0", test, gtest)
	},
}


func Execute() {
	if err := rootcmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
