// awair-exporter - a Prometheus exporter for the Awair Local API (https://support.getawair.com/hc/en-us/articles/360049221014-Awair-Element-Local-API-Feature)

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type airData struct {
	Hostname                         string
	Score                            float64 `json:"score"`
	DewPoint                         float64 `json:"dew_point"`
	Temperature                      float64 `json:"temp"`
	RelativeHumidity                 float64 `json:"humid"`
	AbsoluteHumidity                 float64 `json:"abs_humid"`
	CarbonDioxide                    float64 `json:"co2"`
	CarbonDioxideEstimate            float64 `json:"co2_est"`
	CarbonDioxideEstimateBaseline    float64 `json:"co2_est_baseline"`
	VolatileOrganicCompounds         float64 `json:"voc"`
	VolatileOrganicCompoundsBaseline float64 `json:"voc_baseline"`
	VolatileOrganicCompoundsHydrogen float64 `json:"voc_h2_raw"`
	VolatileOrganicCompoundsEthanol  float64 `json:"voc_ethanol_raw"`
	ParticulateMatter25              float64 `json:"pm25"`
	ParticulateMatter10              float64 `json:"pm10_est"`
}

type awairExporter struct {
	URL string
}

var (
	awairScore = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "awair_score"), "Awair Score.", []string{
			"instance",
		}, nil)
	dewPoint = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "dew_point"), "Dew Point. The temperature the air needs to be cooled to (at constant pressure) in order to achieve a relative humidity of 100%.", []string{
			"instance",
		}, nil)
	temperature = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "temperature"), "Temperature.", []string{
			"instance",
		}, nil)
	relativeHumidity = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "relative_humidity"), "Relative Humidity.", []string{
			"instance",
		}, nil)
	absoluteHumidity = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "absolute_humidity"), "Absolute Humidity.", []string{
			"instance",
		}, nil)
	carbonDioxide = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "co2"), "Carbon Dioxide (CO2) levels", []string{
			"instance",
		}, nil)
	carbonDioxideEstimate = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "co2_estimate"), "Carbon Dioxide (CO2) estimated levels", []string{
			"instance",
		}, nil)
	carbonDioxideEstimateBaseline = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "co2_estimate_baseline"), "Carbon Dioxide (CO2) estimated baseline levels", []string{
			"instance",
		}, nil)
	volatileOrganicCompounds = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "voc"), "Volatile Organic Compounds (VOC) levels", []string{
			"instance",
		}, nil)
	volatileOrganicCompoundsBaseline = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "voc_baseline"), "Volatile Organic Compounds (VOC) baseline levels", []string{
			"instance",
		}, nil)
	volatileOrganicCompoundsHydrogen = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "voc_h2_raw"), "Volatile Organic Compounds (VOC) Molecular Hydrogen raw", []string{
			"instance",
		}, nil)
	volatileOrganicCompoundsEthanol = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "voc_ethanol_raw"), "Volatile Organic Compounds (VOC) Ethanol raw", []string{
			"instance",
		}, nil)
	particulateMatter = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "pm25"), "Particulate Matter 2.5 micrometers or smaller", []string{
			"instance",
		}, nil)
	particulateMatter10 = prometheus.NewDesc(
		prometheus.BuildFQName(
			"awair", "", "pm10_estimate"), "Particulate Matter 10 micrometers or smaller", []string{
			"instance",
		}, nil)
)

func newAwairExporter(url string) *awairExporter {
	return &awairExporter{
		URL: url,
	}
}

// Describe provides the superset of descriptors to the provided channel
func (e *awairExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- awairScore
	ch <- dewPoint
	ch <- temperature
	ch <- relativeHumidity
	ch <- absoluteHumidity
	ch <- carbonDioxide
	ch <- carbonDioxideEstimate
	ch <- carbonDioxideEstimateBaseline
	ch <- volatileOrganicCompounds
	ch <- volatileOrganicCompoundsBaseline
	ch <- volatileOrganicCompoundsHydrogen
	ch <- volatileOrganicCompoundsEthanol
	ch <- particulateMatter
	ch <- particulateMatter10
}

func (e *awairExporter) Collect(ch chan<- prometheus.Metric) {
	endpoint := url.URL{Scheme: "http", Host: e.URL, Path: "air-data/latest"}
	req, err := http.NewRequest(http.MethodGet, endpoint.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &http.Client{}
	req.Header.Set("User-Agent", "github.com/Ichabond/awair-exporter")
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	air := airData{Hostname: e.URL}
	err = json.Unmarshal(data, &air)
	if err != nil {
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(
		awairScore, prometheus.GaugeValue, air.Score, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		dewPoint, prometheus.GaugeValue, air.Score, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		temperature, prometheus.GaugeValue, air.Temperature, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		relativeHumidity, prometheus.GaugeValue, air.RelativeHumidity, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		absoluteHumidity, prometheus.GaugeValue, air.AbsoluteHumidity, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		carbonDioxide, prometheus.GaugeValue, air.CarbonDioxide, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		carbonDioxideEstimate, prometheus.GaugeValue, air.CarbonDioxideEstimate, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		carbonDioxideEstimateBaseline, prometheus.GaugeValue, air.CarbonDioxideEstimateBaseline, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		volatileOrganicCompounds, prometheus.GaugeValue, air.VolatileOrganicCompounds, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		volatileOrganicCompoundsBaseline, prometheus.GaugeValue, air.VolatileOrganicCompoundsBaseline, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		volatileOrganicCompoundsHydrogen, prometheus.GaugeValue, air.VolatileOrganicCompoundsHydrogen, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		volatileOrganicCompoundsEthanol, prometheus.GaugeValue, air.VolatileOrganicCompoundsEthanol, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		particulateMatter, prometheus.GaugeValue, air.ParticulateMatter25, air.Hostname,
	)
	ch <- prometheus.MustNewConstMetric(
		particulateMatter10, prometheus.GaugeValue, air.ParticulateMatter10, air.Hostname,
	)
}
func main() {
	listenAddress := flag.String("l", ":2112", "Listen Address")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage: %s [FLAGS...] HOSTNAME_TO_QUERY\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if len(flag.Args()) == 0 || len(flag.Args()) > 1 {
		log.Fatal("Incorrect arguments passed, see usage.")
	}
	host := flag.Args()[0]
	exporter := newAwairExporter(host)
	prometheus.MustRegister(exporter)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(*listenAddress, nil)
	if err != http.ErrServerClosed {
		log.Fatal(err)
		os.Exit(1)
	}
}
