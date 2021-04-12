// +build ignore
package main

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	fruits = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

// generate random data for line chart
func generateLineItems(rList []float64) []opts.LineData {
	// r is debt ratio
	// p is probability of sending
	pList := make([]opts.LineData, 0)
	for _, r := range rList {
		p := 1 - (1 / (1 + math.Exp(6-3*r)))
		pList = append(pList, opts.LineData{Value: p})
	}
	return pList
}
func generateRList(itemCnt int) ([]float64, []string) {
	left := float64(0)
	step := float64(1)
	rList := make([]float64, 0)
	rSList := make([]string, 0)
	for i := 0; i < itemCnt; i++ {
		r := left + step*float64(i)
		rList = append(rList, r)
		rSList = append(rSList, fmt.Sprintf("%.2f", r))
	}
	return rList, rSList
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:   "smooth style",
			SubLink: "子标题",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "P(send|r)",
			SplitLine: &opts.SplitLine{
				Show: false,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Debt Ratio",
		}),
	)

	// Put data into instance
	rList, rSList := generateRList(5)
	line.SetXAxis(rSList).AddSeries("Category A", generateLineItems(rList)).
		SetSeriesOptions(charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: false,
			}),
		)
	line.Render(w)
}

func main() {
	//	log.Printf("exp(1) = %+v = 2.718.",math.Exp(1))

	http.HandleFunc("/", httpserver)
	log.Println("running server at http://localhost:8848")
	log.Fatal(http.ListenAndServe(":8848", nil))
}
