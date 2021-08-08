package main

import (
	"context"
	"log"
	"time"

	"github.com/hansmi/prometheus-lvm-exporter/lvmreport"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

type collector struct {
	timeout time.Duration
	load    func(context.Context) (*lvmreport.ReportData, error)
	upDesc  *prometheus.Desc
	gc      []*groupCollector
}

func newEmptyCollector() *collector {
	return &collector{
		timeout: time.Minute,
		upDesc:  prometheus.NewDesc("up", "Whether scrape was successful", []string{"status"}, nil),
	}
}

func newCollector() *collector {
	c := newEmptyCollector()

	for _, i := range allGroups {
		c.gc = append(c.gc, newGroupCollector(i))
	}

	return c
}

func newCommandCollector(args []string) *collector {
	cmd := lvmreport.NewCommand(args)

	var sfg singleflight.Group

	c := newCollector()
	c.load = func(ctx context.Context) (*lvmreport.ReportData, error) {
		// Avoid concurrent invocations
		data, err, _ := sfg.Do("", func() (interface{}, error) {
			return cmd.Run(ctx)
		})

		return data.(*lvmreport.ReportData), err
	}

	log.Printf("LVM command: %s", cmd.String())

	return c
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.upDesc

	for _, i := range c.gc {
		i.describe(ch)
	}
}

func (c *collector) collect(ctx context.Context, ch chan<- prometheus.Metric) error {
	data, err := c.load(ctx)
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)

	for _, i := range c.gc {
		i := i
		g.Go(func() error {
			return i.collect(ch, data)
		})
	}

	return g.Wait()
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	status := float64(1)
	statusMsg := ""

	if err := c.collect(ctx, ch); err != nil {
		log.Printf("Scrape failed: %v", err)
		status = 0
		statusMsg = err.Error()
	}

	ch <- prometheus.MustNewConstMetric(c.upDesc, prometheus.GaugeValue, status, statusMsg)
}
