package main

import (
	"fmt"
	"github.com/PlagueCat-Miao/goipfs-test/merkle/mtest"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"log"
	"net/http"
)

var (
	labels = []string{"A+B+C组","A+B组","A+C组","B+C组","A组","B组","C组"}
	Timeee = map[string][2]float32{
		"A+B+C组":{110.23,46.97},
		"A+B组":{67.33,28.51},
		"A+C组":{74.11,31.93},
		"B+C组":{65.28,27.85},
		"A组":{36.57,15.70},
		"B组":{32.18,13.47},
		"C组":{39.85,17.28},
	}
	treeKB []int64
	dagKB []int64

	treeMs []float32
	dagMs []float32
)

// generate random data for line chart
func generateBarItems(Value []int64) []opts.BarData {
	items := make([]opts.BarData, 0)
	itemCnt := len(labels)
	for i := 0; i < itemCnt; i++ {
		items = append(items, opts.BarData{Value: Value[i]})
	}
	return items
}

func generateKeyPerformanceIndicator() error {
	for _, label := range labels {
		T_KB, T_Time, D_KB, D_Time, _, err := mtest.TestSpeedAndSpace(label)
		if err != nil {
			return err
		}
		treeKB = append(treeKB, T_KB)
		dagKB = append(dagKB, D_KB)
		treeMs = append(treeMs,float32(T_Time)/1000)
		dagMs = append(dagMs,float32(D_Time)/1000)
		//treeMs = append(treeMs,Timeee[label][0])
		//dagMs = append(dagMs,Timeee[label][1])
	}
	return nil
}

func httpserver(w http.ResponseWriter, _ *http.Request) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "空间占用对比实验",
			Right: "40%",
		}),
		charts.WithToolboxOpts(opts.Toolbox{
			Show:  true,
			Right: "20%",
			Feature: &opts.ToolBoxFeature{
				SaveAsImage: &opts.ToolBoxFeatureSaveAsImage{
					Show:  true,
					Type:  "png",
					Title: "下载图片",
				},
				DataView: &opts.ToolBoxFeatureDataView{
					Show:  true,
					Title: "数据展示",
					// set the language
					// Chinese version: ["数据视图", "关闭", "刷新"]
					Lang: []string{"data view", "turn off", "refresh"},
				},
			}},
		),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "视频数据组名",
			SplitLine: &opts.SplitLine{
				Show: true,
			},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "占用空间",
			AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} KB"},
			//AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value} s"},
		}),

	)
	// Put data into instance
	bar.SetXAxis(labels).
		AddSeries("merkle Tree", generateBarItems(treeKB)).
		AddSeries("merkle Dag", generateBarItems(dagKB))
	bar.Render(w)
}

func main() {
	err := generateKeyPerformanceIndicator()
	if err != nil {
		fmt.Printf("Err: %v",err)
		return
	}
	//treeKB = append(treeKB,5)
	//dagKB = append(dagKB,5)
	http.HandleFunc("/", httpserver)
	log.Println("running server at http://localhost:8848")
	log.Fatal(http.ListenAndServe(":8848", nil))

}
