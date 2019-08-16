package crawler

import (
	"github.com/tealeg/xlsx"
	"fmt"
)

var (
	Header []string
)

func pushHeader() {
	Header = append(Header, "CNVD-ID")
	Header = append(Header, "公开日期")
	Header = append(Header, "危害级别")
	Header = append(Header, "影响产品")
	Header = append(Header, "CVE ID")
	Header = append(Header, "漏洞描述")
	Header = append(Header, "漏洞类型")
	Header = append(Header, "参考链接")
	Header = append(Header, "漏洞解决方案")
	Header = append(Header, "厂商补丁")
	Header = append(Header, "验证信息")
	Header = append(Header, "报送时间")
	Header = append(Header, "收录时间")
	Header = append(Header, "更新时间")
}
func CreateHeader(sheet *xlsx.Sheet, header []string) {

	row := sheet.AddRow()
	row.SetHeightCM(1)
	for i := 0; i < len(header); i++ {
		cell := row.AddCell()
		cell.Value = header[i]
	}
}
func WriteExecl() {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		fmt.Printf(err.Error())
	}
	pushHeader()

	CreateHeader(sheet, Header)

	for i:=0;i<len(HoleList);i++{
		row := sheet.AddRow()
		row.SetHeightCM(1)
		cell := row.AddCell()
		cell.Value = HoleList[i].CnvdID
		cell = row.AddCell()
		cell.Value = HoleList[i].OpenDate
		cell = row.AddCell()
		cell.Value = HoleList[i].Level
		cell = row.AddCell()
		cell.Value = HoleList[i].AffectProduct
		cell = row.AddCell()
		cell.Value = HoleList[i].CVEID
		cell = row.AddCell()
		cell.Value = HoleList[i].Description
		cell = row.AddCell()
		cell.Value = HoleList[i].HoleType
		cell = row.AddCell()
		cell.Value = HoleList[i].RefLinking
		cell = row.AddCell()
		cell.Value = HoleList[i].Solution
		cell = row.AddCell()
		cell.Value = HoleList[i].Patch
		cell = row.AddCell()
		cell.Value = HoleList[i].AuthInfo
		cell = row.AddCell()
		cell.Value = HoleList[i].SubmitTime
		cell = row.AddCell()
		cell.Value = HoleList[i].CollectionTime
		cell = row.AddCell()
		cell.Value = HoleList[i].UpdateTime
	}
	err = file.Save("C:/Users/HP/Desktop/test_write.xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

}
