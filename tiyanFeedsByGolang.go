package main

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type dtb_applyform struct {
	Id             int `xorm:"pk"`
	User_id        int
	Af_email       string
	Af_mobile      string
	Bb_name        string
	Bb_sex         string
	Bb_year        string
	Bb_month       string
	Bb_day         string
	Pa_name        string
	Pa_sex         string
	Af_province    string
	Af_city        string
	Af_county      string
	Af_add         string
	Af_add_type    string
	Af_post        string
	Af_tel         string
	Af_tel_type    string
	Af_source      string
	Af_source_type string
	Qs_id          int `xorm:"default null"`
	Qs_name        string
	Ql_id          int `xorm:"default null"`
	Ql_name        string
	Cus_id         string
	Hd_id          int `xorm:"default null"`
	Hd_name        string
	Hdl_id         int `xorm:"default null"`
	Hdl_name       string
	Ip             string
	Linkfrom       string
	Apply_type     int
	Ref_mobile     string
	Time           string
	Recom_id       string
	Apps           int
	Listcode       string
	Mobile_check   int
	Appid          int
	Bb_birthday    string
	Source         string
	Isuser         int
	Old_userid     int
	Datasource     int
	Update_on      string
	Isbewill       int
}

type dtb_mobile struct {
	Id     int `xorm:"pk"`
	Mobile string
}

type dtb_applyform_relate struct {
	Id       int `xorm:"pk"`
	Form_id  int
	Tiyanurl string
	Reserve1 string `xorm:"default null"`
	Reserve2 string `xorm:"default null"`
	Reserve3 string `xorm:"default null"`
	Reserve4 string `xorm:"default null"`
}

type dtb_log_visit struct {
	Id             int `xorm:"pk"`
	Supplier_id    int `xorm:"default null"`
	Ql_id          int `xorm:"default null"`
	Hd_id          int `xorm:"default null"`
	Hdl_id         int `xorm:"default null"`
	Source         string
	Recom_id       string
	Visitor_ip     string
	Visit_page_url string
	Visit_time     string
	Is_pc          int
	View_type      int
	Template_id    int
}

func addPvip(engine *xorm.Engine, id int, lv dtb_log_visit) dtb_log_visit {
	lv.Id = id
	lv.Visit_time = time.Now().Format("2006-01-02 15:04:05")
	if rand.Intn(8) > 2 { // 模版
		lv.View_type = 0
		lv.Template_id = rand.Intn(4) + 1
	} else { // ab版
		lv.View_type = rand.Intn(2) + 1
		lv.Template_id = 0
	}
	_, err := engine.Insert(&lv)
	if err != nil {
		log.Println(err.Error())
	}
	return lv
}

func coreHandler(hash, origin_af_mobile int, recomIds, pa_names, ips []string, hds, lms, sources, way []map[string][]byte) {
	engine, err := xorm.NewEngine("mysql", "root:123456@/laiyuan2?charset=utf8")
	if err != nil {
		log.Fatalln(err.Error())
	}

	if hash == 0 {
		hash += 10
	}

	//构建推荐人手机号
	var recomMobiles []int
	tmpRecomMobile := origin_af_mobile
	recomMobilesLen := InsertLen * 10
	for i := 0; i < recomMobilesLen; i++ {
		recomMobiles = append(recomMobiles, tmpRecomMobile)
		tmpRecomMobile++
	}

	recomIds_len, pa_names_len, ips_lenm, hds_len, lms_len, sources_len, way_len := len(recomIds)-1, len(pa_names)-1, len(ips)-1, len(hds)-1, len(lms)-1, len(sources)-1, len(way)-1

	for i := 0; i < InsertLen; i++ {

		applyform := dtb_applyform{
			Bb_year:     "2012",
			Bb_month:    "10",
			Bb_day:      "10",
			Bb_birthday: "2012-10-10",
			Af_province: "云南省",
			Af_city:     "保山市",
			Af_county:   "昌宁县",
			Af_add:      "啊啊啊啊啊啊啊啊啊啊啊啊",
			Af_post:     "123456",
			Apply_type:  3,
			Time:        time.Now().Format("2006-01-02 15:04:05"),
			Update_on:   time.Now().Format("2006-01-02 15:04:05"),
		}

		mobile := dtb_mobile{}

		apply_relate := dtb_applyform_relate{}

		log_visit := dtb_log_visit{
			Visit_page_url: tiyanurl,
		}

		applyform.Id = hash + i*10
		applyform.Mobile_check = rand.Intn(2)
		applyform.Af_mobile = strconv.Itoa(origin_af_mobile + applyform.Id)
		applyform.Pa_name = string(pa_names[rand.Intn(pa_names_len)])
		applyform.Af_city = string(pa_names[rand.Intn(pa_names_len)])
		applyform.Af_add = string(pa_names[rand.Intn(pa_names_len)])
		applyform.Ip = string(ips[rand.Intn(ips_lenm)])
		applyform.Apps = rand.Intn(2)
		applyform.Datasource = 2
		mobile.Id = applyform.Id
		mobile.Mobile = applyform.Af_mobile

		if applyform.Apps == 0 {
			log_visit.Is_pc = 1
		} else {
			log_visit.Is_pc = 0
		}
		log_visit.Visitor_ip = applyform.Ip
		//随机情况,插入名单表和手机表
		if rand.Intn(2) == 1 { //是否主动
			if rand.Intn(2) == 1 { // 友介的情况
				if rand.Intn(2) == 1 { //有source的情况
					applyform.Source = string(sources[rand.Intn(sources_len)]["code"])
				} else {
					applyform.Source = ""
				}
				applyform.Recom_id = recomIds[rand.Intn(recomIds_len)]
				applyform.Ref_mobile = applyform.Af_mobile
				applyform.Ref_mobile = strconv.Itoa(recomMobiles[rand.Intn(recomMobilesLen)])
				log_visit.Recom_id = applyform.Recom_id
				log_visit.Source = applyform.Source
			} else {
				var hq map[string][]byte
				if rand.Intn(2) == 1 { //活动
					hq = hds[rand.Intn(hds_len)]
					applyform.Hd_id, _ = strconv.Atoi(string(hq["qs_id"]))
					applyform.Hd_name = string(hq["huodong"])
					log_visit.Hd_id = applyform.Hd_id
					if rand.Intn(10) != 0 { //有来源
						applyform.Hdl_id, _ = strconv.Atoi(string(hq["id"]))
						applyform.Hdl_name = string(hq["name"])
						log_visit.Hdl_id = applyform.Hdl_id
					}
				} else { //联盟主
					hq = lms[rand.Intn(lms_len)]
					applyform.Qs_id, _ = strconv.Atoi(string(hq["qs_id"]))
					applyform.Qs_name = string(hq["huodong"])
					log_visit.Supplier_id = applyform.Qs_id
					if rand.Intn(10) != 0 { //有联盟
						applyform.Ql_id, _ = strconv.Atoi(string(hq["id"]))
						applyform.Ql_name = string(hq["name"])
						log_visit.Ql_id = applyform.Ql_id
					}
				}
			}
		} else { //主动情况

		}
		log_visit = addPvip(engine, applyform.Id, log_visit) //插入pvip
		_, err := engine.Insert(&applyform)                  //插入名单表
		if err != nil {
			log.Println(err.Error())
		}
		_, err = engine.Insert(&mobile)
		if err != nil {
			log.Println(err.Error())
		}

		apply_relate.Id = applyform.Id
		apply_relate.Form_id = apply_relate.Id
		apply_relate.Tiyanurl = tiyanurl
		apply_relate.Reserve1 = strconv.Itoa(rand.Intn(2) + 1)
		apply_relate.Reserve2 = string(way[rand.Intn(way_len)]["way"])
		if log_visit.View_type != 0 {
			apply_relate.Reserve3 = strconv.Itoa(log_visit.View_type)
		}
		if log_visit.Template_id != 0 {
			apply_relate.Reserve4 = strconv.Itoa(log_visit.Template_id)
		}

		_, err = engine.Insert(&apply_relate)
		if err != nil {
			log.Println(err.Error())
		}
	}

	wg.Done()
}

var wg sync.WaitGroup

const InsertLen int = 5000

var tiyanurl string = "http://m.tiyan.qiaohu.com"

func main() {
	engine, err := xorm.NewEngine("mysql", "root:123456@/laiyuan2?charset=utf8")
	if err != nil {
		log.Fatalln(err.Error())
	}

	pa_names := []string{"刘备", "关羽", "张飞", "赵云", "马超", "黄忠", "诸葛亮", "徐庶", "庞统", "法正", "关平", "关索", "张苞", "关兴", "马岱", "廖化", "魏延", "曹操", "荀彧", "荀攸", "郭嘉", "贾诩", "司马懿", "曹丕", "曹植", "曹冲", "夏侯惇", "夏侯渊", "许褚", "典韦", "徐晃", "张辽", "吕布", "庞德", "于禁", "乐进", "李典", "司马昭", "夏侯霸", "姜维", "钟会", "邓艾", "张颌", "颜良", "文丑", "华雄", "孙坚", "孙策", "孙权", "诸葛恪", "陆逊", "周瑜", "鲁肃", "诸葛瑾", "张昭", "吕蒙", "太史慈", "潘璋", "黄盖", "蒋干"}

	hds, _ := engine.Query("select * from dtb_lianmeng where lhtype = 1 and enableflag = 1")

	lms, _ := engine.Query("select * from dtb_lianmeng where lhtype = 0 and enableflag = 1")

	sources, _ := engine.Query("select * from dtb_youjie_source where enableflag = 1")

	way, _ := engine.Query("select * from dtb_additional_source_rule")

	var recomIds []string

	for i := 0; i < 100; i++ {
		recomIds = append(recomIds, strconv.Itoa(101000000+i))
	}

	var ips []string
	for i := 0; i < 50; i++ {
		ip1, ip2, ip3, ip4 := strconv.Itoa(1+rand.Intn(254)), strconv.Itoa(1+rand.Intn(254)), strconv.Itoa(1+rand.Intn(254)), strconv.Itoa(1+rand.Intn(254))
		ips = append(ips, ip1+"."+ip2+"."+ip3+"."+ip4)
	}

	// start_af_mobile, _ := engine.Query("select max(mobile) as mobile from dtb_mobile")
	// origin_af_mobile, _ := strconv.Atoi(string(start_af_mobile[0]["mobile"]))
	origin_af_mobile := 13000000000

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go coreHandler(i, origin_af_mobile, recomIds, pa_names, ips, hds, lms, sources, way)
	}

	wg.Wait()

}
