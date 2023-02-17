package helper

import (
	"fmt"
	"testing"
	"time"
)

type UserInfo struct {
	NickName   string `xlsx:"昵称"`
	Account    string `xlsx:"账号"`
	Password   string `xlsx:"密码"`
	CreateTime string `xlsx:"创建时间"`
}

func GetUserInfoSlice() []UserInfo {
	return []UserInfo{
		{
			NickName:   "Lwj",
			Account:    "awp85589",
			Password:   "99999",
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			NickName:   "Lwj-1",
			Account:    "awp85588",
			Password:   "99999",
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			NickName:   "Lwj-2",
			Account:    "awp85587",
			Password:   "99999",
			CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

func TestCreateExcelFile(t *testing.T) {
	f, err := CreateExcelFile("sheet1", GetUserInfoSlice(), &SheetConfig{
		ProtectSheet: struct {
			AlgorithmName string
			Password      string
			IsProtect     bool
		}{AlgorithmName: "SHA-512", Password: DefaultPwd, IsProtect: true},
		PictureSheet: struct {
			Config []struct {
				ScaleX float64
				ScaleY float64
				Path   string
			}
			SpaceRow int
		}{Config: []struct {
			ScaleX float64
			ScaleY float64
			Path   string
		}{
			{
				ScaleX: 0.2,
				ScaleY: 0.2,
				Path:   "adapter.jpg",
			},
			{
				ScaleX: 0.3,
				ScaleY: 0.3,
				Path:   "adapter.jpg",
			},
		}, SpaceRow: 100},
	})
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	filePath := GenerateExcelName("")

	if err = f.SaveAs(filePath); err != nil {
		fmt.Println("SaveAs ERR:", err)
		return
	}
}
