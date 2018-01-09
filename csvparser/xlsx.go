package csvparser

import (
	"github.com/tealeg/xlsx"
	"github.com/YAWAL/converter/config"
)

func XLSX(fileName string) (map[string][][]string, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return make(map[string][][]string), err
	}
	res := make(map[string][][]string)
	for _, sheet := range xlFile.Sheets {
		//if sheet not in the list - skip this sheet
		if _, ok := config.Get().Indices[sheet.Name]; !ok {
			continue
		}
		res[sheet.Name] = make([][]string, 0)
		for i, row := range sheet.Rows {
			//first empty row mean the end of the table
			if isRowEmpty(row) {
				break
			}
			res[sheet.Name] = append(res[sheet.Name], []string{})
			//add sources
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[config.Get().Indices[sheet.Name].Page].String())
			//add policy name
			res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[config.Get().Indices[sheet.Name].Name].String())
			//add role`s actions
			for c := config.Get().Indices[sheet.Name].RoleStart; c <= config.Get().Indices[sheet.Name].RoleEnd; c++ {
				res[sheet.Name][i] = append(res[sheet.Name][i], row.Cells[c].String())
			}
			//for _, cell := range row.Cells {
			//	res[sheet.Name][i] = append(res[sheet.Name][i], cell.String())
			//}
		}
	}
	return res, nil
}

func isRowEmpty(row *xlsx.Row) bool {
	if len(row.Cells) == 0 {
		return true
	}
	for _, r := range row.Cells {
		if r.String() != "" {
			return false
		}
	}
	return true
}
