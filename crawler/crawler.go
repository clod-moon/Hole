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
	CVEID          string `gorm:"type:varchar(30)"`
	Description    string `gorm:"type:text"`
	HoleType       string `gorm:"type:varchar(30)"`
	RefLinking     string `gorm:"type:text"`
	Solution       string `gorm:"type:text"`
	Patch          string `gorm:"type:text"`
	AuthInfo       string `gorm:"type:text"`
	SubmitTime     string `gorm:"type:varchar(30)"`
	CollectionTime string `gorm:"type:varchar(30)"`
	UpdateTime     string `gorm:"type:varchar(30)"`
}

var (
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

	WriteExecl()
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
	//fmt.Println("len:", len(HoleList))

	DBHd.Find(&HoleList,`cnvd_id = ""`)

	fmt.Println("len:",len(HoleList))

	for i := 0; i < len(HoleList); i++ {
		ParseHole(&HoleList[i])
		time.Sleep(time.Second*3)
	}
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

	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate, br")
	reqest.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	reqest.Header.Add("Cache-Control", "max-age=0")
	reqest.Header.Add("Connection", "Keep-Alive")
	reqest.Header.Add("Cookie", "JSESSIONID=05DCDCA2AE084CEB0F45FDE197CE8A9E; __jsl_clearance=1565972211.667|0|YqiRYqKoW8ardzabAGKB3q2UKK8%3D; __jsluid_s=1562270a633103338cd4b3fffd9e1ef8; __jsluid_h=5eeca75fab78d68d624506f6a6760378; puk=077df5d2afa8609862b895e3ed271fa410faba86a4fc40a1a8a6c04929f030cea6649567ec415a23e5026db436036cad4da1809494523bd49af54d388bc7060fa16949beba6e8fa7747837a03541794131f76b74596b843193d7c3f26b8354c5b815fe8c8a2ed15298834005267c87937b0955568a34aaec58eaf8baf4d1c1eb")
	reqest.Header.Add("Host", "www.cnvd.org.cn")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18362")
	reqest.Header.Add("Referer", h.Url)
	reqest.Header.Add("Upgrade", "1")
	//reqest.Header.Add("Sec-Fetch-Mode","navigate")
	//reqest.Header.Add("Sec-Fetch-Site","none")
	//reqest.Header.Add("Sec-Fetch-User","?1")

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

func ParseHole(h *Hole) {

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
	fmt.Println(h.CnvdID)
	DBHd.Save(h)
}
