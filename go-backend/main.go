package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

const (
	externalMainAPI   = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancDetail?page=1&perPage=100"
	externalDetailAPI = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancMdl?page=1&perPage=100&cond%5BHOUSE_MANAGE_NO%3A%3AEQ%5D="
)

type ApiResponse struct {
	Data []MainData `json:"data"`
}

type MainData struct {
	HouseManageNo        string       `json:"HOUSE_MANAGE_NO"`
	RcritPblancDe        string       `json:"RCRIT_PBLANC_DE"`
	HouseDtlSecdNm       string       `json:"HOUSE_DTL_SECD_NM"`
	SubscrptAreaCodeNm   string       `json:"SUBSCRPT_AREA_CODE_NM"`
	RentSecdNm           string       `json:"RENT_SECD_NM"`
	HouseNm              string       `json:"HOUSE_NM"`
	RceptEndde           string       `json:"RCEPT_ENDDE"`
	SpecltRdnEarthAt     string       `json:"SPECLT_RDN_EARTH_AT"`
	LrsclBldlndAt        string       `json:"LRSCL_BLDLND_AT"`
	NplnPrvoprPublicHouseAt string     `json:"NPLN_PRVOPR_PUBLIC_HOUSE_AT"`
	DetailData           []DetailData `json:"DETAIL_DATA"`
}

type DetailData struct {
	HouseTy          string `json:"HOUSE_TY"`
	LttotTopAmount   string `json:"LTTOT_TOP_AMOUNT"`
	SupplyHshldco    int    `json:"SUPLY_HSHLDCO"`
	LocalPoint       int    `json:"LOCAL_POINT"`
	LocalRandZero    int    `json:"LOCAL_RAND_ZERO"`
	LocalRandZeroOne int    `json:"LOCAL_RAND_ZERO_ONE"`
	EtcGGPoint       int    `json:"ETC_GG_POINT"`
	EtcGGRandZero    int    `json:"ETC_GG_RAND_ZERO"`
	EtcGGRandZeroOne int    `json:"ETC_GG_RAND_ZERO_ONE"`
	EtcPoint         int    `json:"ETC_POINT"`
	EtcRandZero      int    `json:"ETC_RAND_ZERO"`
	EtcRandZeroOne   int    `json:"ETC_RAND_ZERO_ONE"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/data", getData)
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	log.Fatal(http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux)))
}

func getData(w http.ResponseWriter, r *http.Request) {
	serviceKey := os.Getenv("CHUNGYAK_INFO_API_KEY")
	mainURL := fmt.Sprintf("%s&serviceKey=%s", externalMainAPI, serviceKey)

	resp, err := http.Get(mainURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	serviceKey = os.Getenv("CHUNGYAK_INFO_API_KEY")
	for i, item := range apiResponse.Data {
		detailURL := fmt.Sprintf("%s%s&serviceKey=%s", externalDetailAPI, item.HouseManageNo, serviceKey)
		if err := fetchDetailData(detailURL, &apiResponse.Data[i]); err != nil {
			http.Error(w, "Failed to fetch detail data", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse)
}

func fetchDetailData(url string, mainData *MainData) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var detailResponse struct {
		Data []DetailData `json:"data"`
	}
	if err := json.Unmarshal(body, &detailResponse); err != nil {
		return err
	}

	for _, detail := range detailResponse.Data {
		mainData.DetailData = append(mainData.DetailData, processDetailData(mainData, detail))
	}
	return nil
}

func processDetailData(mainData *MainData, detail DetailData) DetailData {
	총물량 := detail.SupplyHshldco
	지역 := mainData.SubscrptAreaCodeNm
	투기과열지구 := mainData.SpecltRdnEarthAt
	대규모 := mainData.LrsclBldlndAt
	수도권공공주택지구 := mainData.NplnPrvoprPublicHouseAt
	타입 := detail.HouseTy
	detail.LttotTopAmount = convertAndTrim(detail.LttotTopAmount)
	if 대규모 == "N" {
		assignNonLargeScaleValues(수도권공공주택지구, 투기과열지구, 타입, 총물량, &detail)
	}else{
		if 지역 == "서울" || 지역 == "인천" {
			assignLargeScaleValues(수도권공공주택지구, 투기과열지구, 타입, 총물량, &detail,  []int{50, 50})
		} else if 지역 == "경기" {
			assignLargeScaleValues(수도권공공주택지구, 투기과열지구, 타입, 총물량, &detail, []int{30, 20, 50})
		}
	}		 
	return detail
}
func assignNonLargeScaleValues(publicHouse, earthAt, houseType string, supply int, detail *DetailData) {
	var pointsFunc func(string, int) int
	if publicHouse == "Y" || earthAt == "Y" {
		pointsFunc = pointsFunc478
	} else {
		pointsFunc = pointsFunc40
	}
	detail.LocalPoint = pointsFunc(houseType, supply)
	detail.LocalRandZero = int(math.Ceil(float64(supply-detail.LocalPoint) * 0.75))
	detail.LocalRandZeroOne = supply - detail.LocalRandZero - detail.LocalPoint
}

func assignLargeScaleValues(수도권공공주택지구, 투기과열지구, 타입 string, 총물량 int, detail *DetailData,percentages []int) {
	var pointsFunc func(string, int) int
	if 수도권공공주택지구 == "Y" || 투기과열지구 == "Y" {
		pointsFunc = pointsFunc478
	} else {
		pointsFunc = pointsFunc40
	}
	적정물량:= allocateN(총물량, percentages)	
		for i, 물량 := range 적정물량 {			
			가점자물량, 랜덤무주택물량, 랜덤무주택일주택물량 := calculateValues(물량,타입,pointsFunc)

			switch len(적정물량) {
			case 2:
				if i == 0 {
					detail.LocalPoint = 가점자물량
					detail.LocalRandZero = 랜덤무주택물량
					detail.LocalRandZeroOne = 랜덤무주택일주택물량
				} else if i == 1 {
					detail.EtcPoint = 가점자물량
					detail.EtcRandZero = 랜덤무주택물량
					detail.EtcRandZeroOne = 랜덤무주택일주택물량
				}
			case 3:
				if i == 0 {
					detail.LocalPoint = 가점자물량
					detail.LocalRandZero = 랜덤무주택물량
					detail.LocalRandZeroOne = 랜덤무주택일주택물량
				} else if i == 1 {
					detail.EtcGGPoint = 가점자물량
					detail.EtcGGRandZero = 랜덤무주택물량
					detail.EtcGGRandZeroOne = 랜덤무주택일주택물량
				} else if i == 2 {
					detail.EtcPoint = 가점자물량
					detail.EtcRandZero = 랜덤무주택물량
					detail.EtcRandZeroOne = 랜덤무주택일주택물량
				}
			}
		}			
}
func calculateValues( 물량 int, 타입 string, pointsFunc func(string, int) int) (int, int, int) {
		
	가점자물량 := pointsFunc(타입, 물량)
	랜덤물량 := 물량 - 가점자물량
	랜덤무주택물량 := int(math.Ceil(float64(랜덤물량) * 0.75))
	랜덤무주택일주택물량 := 랜덤물량 - 랜덤무주택물량

	return 가점자물량, 랜덤무주택물량, 랜덤무주택일주택물량
}

func allocateN(총물량 int, percentages []int) []int {
	var results []int

	for _, percent := range percentages {
		part := float64(총물량) * float64(percent) / 100
		roundedPart := int(math.Round(part))
		results = append(results, roundedPart)
	}

	totalAllocated := 0
	for _, value := range results {
		totalAllocated += value
	}

	if totalAllocated > 총물량 {
		surplus := totalAllocated - 총물량
		for i := len(results) - 1; i >= 0 && surplus > 0; i-- {
			if results[i] > 0 {
				reduction := min(surplus, results[i])
				results[i] -= reduction
				surplus -= reduction
			}
		}
	} else if totalAllocated < 총물량 {
		deficit := 총물량 - totalAllocated
		for i := 0; i < len(results) && deficit > 0; i++ {
			addition := min(deficit, results[i])
			results[i] += addition
			deficit -= addition
		}
	}

	return results
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func cleanNumberString(input string) string {
	re := regexp.MustCompile(`[^0-9.]`)
	return re.ReplaceAllString(input, "")
}

func pointsFunc478(houseType string, supply int) int {
	houseM2, err := strconv.ParseFloat(cleanNumberString(houseType), 64)
	if err != nil {
		return 0
	}

	switch {
	case houseM2 <= 60:
		return int(math.Ceil(float64(supply) * 0.40))
	case houseM2 <= 85:
		return int(math.Ceil(float64(supply) * 0.70))
	default:
		return int(math.Ceil(float64(supply) * 0.80))	
	}
}
func pointsFunc40(houseType string, supply int) int {
	houseM2, err := strconv.ParseFloat(cleanNumberString(houseType), 64)
	if err != nil {
		return 0
	}

	if houseM2 <= 85 {
		return int(math.Ceil(float64(supply) * 0.40))
	}
	return 0
}
func convertAndTrim(numberStr string) string {
    number, err := strconv.ParseFloat(numberStr, 64)
    if err != nil {
        fmt.Println("Error parsing string:", err)
        return ""
    }
    number = number / 10000
    number = math.Round(number*10) / 10
    return fmt.Sprintf("%.1f", number)
}