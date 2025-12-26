package mirg

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func GetPic(fileName string, outDir string, stockCodeNum string, picNum string) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err = os.MkdirAll(outDir, 0777)
		if err != nil {
			panic(err)
		}
	}

	index := 0
	for {
		sheet := f.GetSheetName(index)
		rows, err := f.GetRows(sheet)
		if err != nil {
			return
		}

		for rowIdx := range rows {
			rowNum := rowIdx + 1

			pics, err := f.GetPictures(sheet, fmt.Sprintf("%s%d", picNum, rowNum))
			if err != nil || len(pics) == 0 {
				continue
			}

			picPath := path.Join(outDir, sheet)
			if _, err := os.Stat(picPath); os.IsNotExist(err) {
				err = os.MkdirAll(picPath, 0777)
				if err != nil {
					fmt.Println("创建目录失败:", err)
				}
			}

			// B 列作为文件名
			name, _ := f.GetCellValue(sheet, fmt.Sprintf("%s%d", stockCodeNum, rowNum))
			if name == "" {
				continue
			}

			for index, pic := range pics {
				ext := pic.Extension
				if ext == "" {
					ext = ".jpg"
				}
				filename := sanitize(name)
				if len(pics) > 1 {
					filename = filename + "_" + strconv.Itoa(index+1)
				}
				filename = filename + ext
				err := os.WriteFile(filepath.Join(picPath, filename), pic.File, 0644)
				if err != nil {
					fmt.Println("写文件失败:", err)
				}
			}
		}
		index++
		fmt.Println("图片导出完成")
	}
}

func sanitize(name string) string {
	// 防止非法文件名
	r := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return r.Replace(name)
}
