package show

import (
	"io"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var (
	NodeColumn = [...]string{" ", " ", " ", " ", " "}

	NodeRow = [...]string{
		"5", "10", "15", "20", "25", "30", "35", "40", "45", "50",
		"55", "60", "65", "70", "75", "80", "85", "90", "95", "100",
	}
)

func genHeatMapData(nodesRem []int64) []opts.HeatMapData {
	items := make([]opts.HeatMapData, 0)
	n := len(nodesRem)
	columnLen := len(NodeColumn)
	sort.Slice(nodesRem, func(i, j int) bool {
		return nodesRem[i] > nodesRem[j]
	})

	for i := 0; i < n; i++ {
		column := i % columnLen
		row := i / columnLen
		if nodesRem[i] == 0 {
			items = append(items, opts.HeatMapData{Value: [3]interface{}{row, column, "-"}})
		} else {
			items = append(items, opts.HeatMapData{Value: [3]interface{}{row, column, nodesRem[i]}})
		}
	}
	return items
}

func heatMapBase(Title string, nodesUse []int64) *charts.HeatMap {
	hm := charts.NewHeatMap()
	hm.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: Title,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type:      "category",
			Data:      NodeRow,
			SplitArea: &opts.SplitArea{Show: true},
			//AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value}"},
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type:      "category",
			Data:      NodeColumn,
			SplitArea: &opts.SplitArea{Show: true},
			//AxisLabel: &opts.AxisLabel{Show: true, Formatter: "{value}"},
		}),

		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Min:        200,
			Max:        550,
			InRange: &opts.VisualMapInRange{
				Color: []string{"#eac736", "#ff0000"},
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1080px",
			Height: "360px",
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
	)

	hm.SetXAxis(NodeRow).AddSeries("heatmap", genHeatMapData(nodesUse))
	return hm
}

type HeatmapExamples struct{}

func (HeatmapExamples) Experiments(StrategyHeat map[string][]int64) {
	page := components.NewPage()
	var Charters []components.Charter
	for title, nodesUse := range StrategyHeat {
		Charters = append(Charters, heatMapBase(title, nodesUse))
	}
	page.AddCharts(Charters...)

	f, err := os.Create("show/html/heatmap.html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
}
