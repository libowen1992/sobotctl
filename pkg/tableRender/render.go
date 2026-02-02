package tableRender

import (
	"github.com/olekukonko/tablewriter"
	"os"
)

func Render(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetRowLine(true)
	var colors []tablewriter.Colors
	for i := 0; i < len(headers); i++ {
		colors = append(colors, tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor})
	}
	table.AppendBulk(data)
	table.Render()
}
