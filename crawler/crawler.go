package crawler

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"bytes"
	"regexp"
	"strings"
	"time"
)

type Hole struct {
	ID             int    `gorm:"primary_key;type:int(11);AUTO_INCREMENT`
	Title          string `gorm:"type:text`
	Url            string `gorm:"type:varchar(255);not null`
	CnvdID         string `gorm:"type:text`
	OpenDate       string `gorm:"type:varchar(30)"`
	Level          string `gorm:"type:varchar(30)"`
	AffectProduct  string `gorm:"type:text"`
	CVEID          string `gorm:"type:text"`
	Description    string `gorm:"type:text"`
	HoleType       string `gorm:"type:text"`
	RefLinking     string `gorm:"type:text"`
	Solution       string `gorm:"type:text"`
	Patch          string `gorm:"type:text"`
	AuthInfo       string `gorm:"type:text"`
	SubmitTime     string `gorm:"type:text"`
	CollectionTime string `gorm:"type:text"`
	UpdateTime     string `gorm:"type:text"`
}

var (
	Cookie     string
	Accept     string
	Coding     string
	Language   string
	Control    string
	Connection string
	CNDVHost       string
	Agent      string
	HoleList   []Hole
	rgx        = regexp.MustCompile(`<span class="([a-z]+) showInfo"></span>`)
	rgxProduct = regexp.MustCompile(`([^<]*)<br/>`)
)

const RED = "red"
const YELLOW = "yellow"
const GREEN = "green"

func transLevel(color string) string {
	if color == RED {
		return "高"
	} else if color == YELLOW {
		return "中"
	} else if color == GREEN {
		return "低"
	} else {
		return "低"
	}
	return "低"
}

func StartCrawler() {

	Init()

	GetPages()

	//WriteExecl()
}

func GetPages() {

	//for i := 0; i < 3; i++ {
	//	url := fmt.Sprintf("https://ics.cnvd.org.cn/?max=1000&offset=%d", i*1000)
	//	doc, err := goquery.NewDocument(url)
//
	//	if err != nil {
	//		return
	//	}
	//	ParsePages(doc)
	//	time.Sleep(time.Second * 5)
	//}
//
	fmt.Println("len:", len(HoleList))

	DBHd.Find(&HoleList,`ref_linking = ""`)

	fmt.Println("len:",len(HoleList))

	for i := 0; i < 7; i++ {
		ParseHole(&HoleList[i])
		fmt.Println(HoleList[i].ID)
		time.Sleep(time.Second*5)
	}

	//var h Hole
	//h.ID = 812
	//h.Url = `https://www.cnvd.org.cn/flaw/show/CNVD-2019-22236`
	//ParseHole(&h)
	//fmt.Println(h)
}

func ParsePages(doc *goquery.Document) {
	//HoleList = append(HoleList, Hole{Page: 1, Url: ""})
	doc.Find("#tr td a ").Each(func(i int, s *goquery.Selection) {

		var h Hole

		h.Url, _ = s.Attr("href")

		h.Title, _ = s.Attr("title")

		var th Hole
		DBHd.Find(&th, "url = ?", h.Url)
		if th.ID == 0{
			DBHd.Create(h)
			HoleList = append(HoleList, h)
		}else{
			if len(th.CnvdID) == 0{
				HoleList = append(HoleList, th)
			}
		}
	})
	return
}

func ParseHolePage(h *Hole) (*goquery.Document, error) {
	client := &http.Client{}
	//生成要访问的url

	//提交请求
	reqest, err := http.NewRequest("GET", h.Url, nil)
	//增加header选项

	reqest.Header.Add("Accept", Accept)
	reqest.Header.Add("Accept-Encoding", Coding)
	reqest.Header.Add("Accept-Language", Language)
	reqest.Header.Add("Cache-Control", Control)
	reqest.Header.Add("Connection", Connection)
	reqest.Header.Add("Cookie", Cookie)
	reqest.Header.Add("Host", CNDVHost)
	reqest.Header.Add("User-Agent", Agent)
	reqest.Header.Add("Referer", h.Url)
	//reqest.Header.Add("Upgrade", "1")

	if err != nil {
		panic(err)
	}
	//处理返回结果
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	ret, err := ParseGzip(body, false)
	r := bytes.NewReader(ret)

	return goquery.NewDocumentFromReader(r)

}

func ParseHole2(h *Hole) {

	doc, err := ParseHolePage(h)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	doc.Find("table.gg_detail  tbody tr td").Each(func(i int, s *goquery.Selection) {
		switch i {
		case 1:
			h.CnvdID = strings.TrimSpace(s.Text())
			break
		case 3:
			h.OpenDate = strings.TrimSpace(s.Text())
			break
		case 5:
			str, _ := s.Html()
			if len(rgx.FindStringSubmatch(str)) > 1 {
				h.Level = transLevel(rgx.FindStringSubmatch(str)[1])
			} else {
				h.Level = "低"
			}
			break
		case 7:
			str, _ := s.Html()
			strs := rgxProduct.FindAllStringSubmatch(str, -1)
			for i := 0; i < len(strs); i++ {
				if i != 0 {
					h.AffectProduct += "\n\r"
				}
				h.AffectProduct += strings.TrimSpace(strs[i][1])
			}
			break
		case 9:
			h.CVEID = strings.TrimSpace(s.Text())
			break
		case 11:
			str, _ := s.Html()
			strs := rgxProduct.FindAllStringSubmatch(str, -1)
			for i := 0; i < len(strs); i++ {
				h.Description += strings.TrimSpace(strs[i][1])
			}
			break
		case 13:
			str, _ := s.Html()
			h.HoleType = strings.TrimSpace(str)
			break
		case 15:
			h.RefLinking, _ = s.Find("a").Attr("href")
			break
		case 17:
			str, _ := s.Html()
			strs := rgxProduct.FindAllStringSubmatch(str, -1)
			for i := 0; i < len(strs); i++ {
				h.Solution += strings.TrimSpace(strs[i][1])
			}
			break
		case 19:
			h.Patch = strings.TrimSpace(s.Text())
			break
		case 21:
			h.AuthInfo = strings.TrimSpace(s.Text())
			break
		case 23:
			h.SubmitTime = strings.TrimSpace(s.Text())
			break
		case 25:
			h.CollectionTime = strings.TrimSpace(s.Text())
			break
		case 27:
			h.UpdateTime = strings.TrimSpace(s.Text())
			break
		default:
			break
		}
	})
	fmt.Println(h.AffectProduct)
	fmt.Println(h.Description)
	//fmt.Println(h.CnvdID)
	//DBHd.Save(h)
}

func ParseHole(h *Hole) {

	doc, err := ParseHolePage(h)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	doc.Find("table.gg_detail  tbody tr ").Each(func(i int, s *goquery.Selection) {
		//fmt.Println(s.Text())
		align := s.Find("td.alignRight").Text()
		s.Find("td").Each(func(t int, td *goquery.Selection) {
			//fmt.Println(td.Next().Text())
			if align == "CNVD-ID"{
				h.CnvdID = strings.TrimSpace(td.Text())
			}else if align == "公开日期"{
				h.OpenDate = strings.TrimSpace(td.Text())
			} else if align == "危害级别"{
				str, _ := td.Next().Html()
				if len(rgx.FindStringSubmatch(str)) > 1 {
					h.Level = transLevel(rgx.FindStringSubmatch(str)[1])
				} else {
					h.Level = "低"
				}
			} else if align == "影响产品"{
				str, _ := td.Next().Html()
				strs := rgxProduct.FindAllStringSubmatch(str, -1)
				for i := 0; i < len(strs); i++ {
					if i != 0 {
						h.AffectProduct += "\n\r"
					}
					h.AffectProduct += strings.TrimSpace(strs[i][1])
				}
			} else if align == "CVE ID"{
				h.CVEID = strings.TrimSpace(td.Text())
			} else if align == "漏洞描述"{
				str, _ := td.Next().Html()
				strs := rgxProduct.FindAllStringSubmatch(str, -1)
				for i := 0; i < len(strs); i++ {
					h.Description += strings.TrimSpace(strs[i][1])
				}
			} else if align == "漏洞类型"{
				str, _ := td.Html()
				h.HoleType = strings.TrimSpace(str)
			} else if align == "参考链接"{
				h.RefLinking, _ = td.Find("a").Attr("href")
				if len(h.RefLinking) < 1{
					h.RefLinking = h.Url
				}
			} else if align == "漏洞解决方案"{
				str, _ := td.Next().Html()
				strs := rgxProduct.FindAllStringSubmatch(str, -1)
				for i := 0; i < len(strs); i++ {
					h.Solution += strings.TrimSpace(strs[i][1])
				}
			} else if align == "厂商补丁"{
				h.Patch = strings.TrimSpace(td.Text())
			} else if align == "验证信息"{
				h.AuthInfo = strings.TrimSpace(td.Text())
			} else if align == "报送时间"{
				h.SubmitTime = strings.TrimSpace(td.Text())
			} else if align == "收录时间"{
				h.CollectionTime = strings.TrimSpace(td.Text())
			} else if align == "更新时间"{
				h.UpdateTime = strings.TrimSpace(td.Text())
			}
		})
	})

	DBHd.Save(h)
}