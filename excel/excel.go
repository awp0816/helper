package excel

import (
	"fmt"
	"helper/ti"
	"log"
	"reflect"
	"strings"

	"github.com/xuri/excelize/v2"
)

var (
	//生成文件后缀
	Suffix = ".xlsx"
	//算法名称
	Algorithms = []string{"MD4", "MD5", "SHA-1", "SHA-256", "SHA-384", "SHA-512"}
	//默认密码
	DefaultPwd = "123456"
	//构造列名称
	Cols = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L",
		"M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
)

type SheetConfig struct {
	ProtectSheet struct {
		AlgorithmName string //MD4、MD5、SHA-1、SHA-256、SHA-384 或 SHA-512
		Password      string //密码
		IsProtect     bool   //是否保护当前sheet
	}
	PictureSheet struct {
		Config []struct {
			ScaleX float64 //图片缩放比例
			ScaleY float64 //图片缩放比例
			Path   string  //图片路径
		}
		SpaceRow int //图片直接空多少行。默认最小10
	}
}

func _CheckAlgorithm(algorithm string) string {
	var key string
	for _, v := range Algorithms {
		if algorithm == v {
			key = algorithm
		}
	}
	if key == "" {
		key = Algorithms[0]
		log.Println("未支持的加密方式,已改用SHA-512")
	}
	return key
}

func _CheckPassword(password string) string {
	if len(strings.TrimSpace(password)) == 0 {
		log.Println("已输入的密码为空,已改用123456")
		return DefaultPwd
	}
	return password
}

func CreateExcelFile(sheet string, records interface{}, config *SheetConfig) (xlsx *excelize.File, err error) {
	var (
		sheetIndex int
		style      int
		enable     = true
		imageCell  string
	)
	if sheet == "" {
		sheet = "sheet1"
	}
	xlsx = excelize.NewFile()
	if sheetIndex, err = xlsx.NewSheet(sheet); err != nil {
		return
	}
	xlsx.SetActiveSheet(sheetIndex)

	if style, err = xlsx.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	}); err != nil {
		return
	}

	t := reflect.TypeOf(records)
	if t.Kind() != reflect.Slice {
		return
	}
	s := reflect.ValueOf(records)
	for i := 0; i < s.Len(); i++ {
		elem := s.Index(i).Interface()
		elemType := reflect.TypeOf(elem)
		elemValue := reflect.ValueOf(elem)
		index := -1
		for j := 0; j < elemType.NumField(); j++ {
			field := elemType.Field(j)
			tag := field.Tag.Get("xlsx")
			name := tag
			if tag == "" || tag == "-" {
				continue
			}
			index++
			if index == -1 {
				continue
			}
			column := GetExcelColumn(index)
			//设置表头
			if i == 0 {
				headCell := fmt.Sprintf("%s%d", column, i+1)
				if err = xlsx.SetCellStyle(sheet, headCell, headCell, style); err != nil {
					return
				}
				if err = xlsx.SetCellValue(sheet, headCell, name); err != nil {
					return
				}
			}
			//设置内容
			bodyCell := fmt.Sprintf("%s%d", column, i+2)
			if err = xlsx.SetCellStyle(sheet, bodyCell, bodyCell, style); err != nil {
				return
			}
			if err = xlsx.SetCellValue(sheet, bodyCell, elemValue.Field(j).Interface()); err != nil {
				return
			}
		}
	}
	if config != nil {
		//工作表保护
		if config.ProtectSheet.IsProtect {
			if err = xlsx.ProtectSheet(sheet, &excelize.SheetProtectionOptions{
				AlgorithmName:       _CheckAlgorithm(config.ProtectSheet.AlgorithmName),
				Password:            _CheckPassword(config.ProtectSheet.Password),
				SelectLockedCells:   true,
				SelectUnlockedCells: true,
				EditScenarios:       true,
			}); err != nil {
				return
			}
		}
		//插入图片
		if config.PictureSheet.Config != nil {
			if _, err = xlsx.NewSheet("images"); err != nil {
				return
			}
			if config.PictureSheet.SpaceRow < 10 {
				config.PictureSheet.SpaceRow = 10
			}
			for i, v := range config.PictureSheet.Config {
				if i == 0 {
					imageCell = fmt.Sprintf("A1")
				} else {
					imageCell = fmt.Sprintf("A%d", config.PictureSheet.SpaceRow*i)
				}
				if err = xlsx.AddPicture("images", imageCell, v.Path, &excelize.GraphicOptions{
					ScaleX:      v.ScaleX,
					ScaleY:      v.ScaleY,
					PrintObject: &enable,
				}); err != nil {
					return
				}
			}
		}
	}
	return
}

func GetExcelColumn(num int) string {
	var column string
	v := num + 1
	for v > 0 {
		k := v % 26
		if k == 0 {
			k = 26
		}
		v = (v - k) / 26
		column = Cols[k] + column
	}
	return column
}

func GenerateExcelName(name string) string {
	if name == "" {
		return fmt.Sprintf("未命名-%s%s", ti.PrintTime(ti.ExcelTime), Suffix)
	} else {
		if strings.Contains(name, Suffix) {
			index := strings.Index(name, Suffix)
			if index == 0 {
				if name[5:] == "" {
					return fmt.Sprintf("未命名-%s%s", ti.PrintTime(ti.ExcelTime), Suffix)
				} else {
					return fmt.Sprintf("%s-%s%s", name[5:], ti.PrintTime(ti.ExcelTime), Suffix)
				}
			} else if len(name)-index == 5 {
				return fmt.Sprintf("%s-%s%s", name[:len(name)-5], ti.PrintTime(ti.ExcelTime), Suffix)
			} else {
				return fmt.Sprintf("%s-%s%s", name[:index], ti.PrintTime(ti.ExcelTime), Suffix)
			}
		} else {
			return fmt.Sprintf("%s-%s%s", name, ti.PrintTime(ti.ExcelTime), Suffix)
		}
	}
}
