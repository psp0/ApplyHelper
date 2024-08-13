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
	"time"
	"sort"
	"context"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
)

const (
	externalSELAPTMainAPI   = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancDetail?page=1&perPage=100&cond%%5BSUBSCRPT_AREA_CODE_NM%%3A%%3AEQ%%5D=%%EC%%84%%9C%%EC%%9A%%B8&cond%%5BRCRIT_PBLANC_DE%%3A%%3AGTE%%5D=%s&cond%%5BRCRIT_PBLANC_DE%%3A%%3ALT%%5D=%s"
	externalKYGAPTMainAPI  = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancDetail?page=1&perPage=100&cond%%5BSUBSCRPT_AREA_CODE_NM%%3A%%3AEQ%%5D=%%EC%%9D%%B8%%EC%%B2%%9C&cond%%5BRCRIT_PBLANC_DE%%3A%%3AGTE%%5D=%s&cond%%5BRCRIT_PBLANC_DE%%3A%%3ALT%%5D=%s"
	externalINCAPTMainAPI  = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancDetail?page=1&perPage=100&cond%%5BSUBSCRPT_AREA_CODE_NM%%3A%%3AEQ%%5D=%%EA%%B2%%BD%%EA%%B8%%B0&cond%%5BRCRIT_PBLANC_DE%%3A%%3AGTE%%5D=%s&cond%%5BRCRIT_PBLANC_DE%%3A%%3ALT%%5D=%s"
	externalAPTDetailAPI = "https://api.odcloud.kr/api/ApplyhomeInfoDetailSvc/v1/getAPTLttotPblancMdl?page=1&perPage=20&cond%5BHOUSE_MANAGE_NO%3A%3AEQ%5D="
)

type ApiResponse struct {
	Data []MainData `json:"data"`
}

type MainData struct {
	HouseManageNo        string       `json:"HOUSE_MANAGE_NO"`
	RcritPblancDe        string       `json:"RCRIT_PBLANC_DE"`
	HouseDtlSecdNm       string       `json:"HOUSE_DTL_SECD_NM"`
	SubscrptAreaCodeNm   string       `json:"SUBSCRPT_AREA_CODE_NM"`	
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
	EtcKYGPoint       int    `json:"ETC_KYG_POINT"`
	EtcKYGRandZero    int    `json:"ETC_KYG_RAND_ZERO"`
	EtcKYGRandZeroOne int    `json:"ETC_KYG_RAND_ZERO_ONE"`
	EtcPoint         int    `json:"ETC_POINT"`
	EtcRandZero      int    `json:"ETC_RAND_ZERO"`
	EtcRandZeroOne   int    `json:"ETC_RAND_ZERO_ONE"`
}

func main() {
	// Load the .env.local file
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatal("Error loading .env.local file")
	}

	// Initialize the tracer and ensure it shuts down gracefully
	shutdown := initTracer()
	defer shutdown()

	// HTTP handlers with tracing
	http.HandleFunc("/APT", traceMiddleware(getAPTData))
	http.HandleFunc("/RemndrAPT", traceMiddleware(getRemndrAPTData))

	// CORS configuration
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Start the server with CORS and tracing enabled
	log.Fatal(http.ListenAndServe(":8080", corsHandler(http.DefaultServeMux)))
}
func traceMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tracer := otel.Tracer("go-backend")
		ctx, span := tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		// Pass the traced context to the next handler
		next(w, r.WithContext(ctx))
	}
}

func initTracer() func() {
    // Create the Jaeger exporter    
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger:14268/api/traces")))

    if err != nil {
        log.Fatal(err)
    }

    // Create the tracer provider with the exporter
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName("go-backend"),
        )),
    )

    otel.SetTracerProvider(tp)

    return func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Fatal(err)
        }
    }
}
func getRemndrAPTData(w http.ResponseWriter, r *http.Request) {
	
}
func getAPTData(w http.ResponseWriter, r *http.Request) {
	//처리시간 계산
	start := time.Now()
	defer func() {
		elapsed := time.Since(start)
		log.Printf("Total Processing time: %s", elapsed)
	}()
		
	serviceKey := os.Getenv("CHUNGYAK_INFO_API_KEY")
	pages := 1 // Number of pages you want to go
	startDate, endDate :=getDateRangeForPage(pages)			
	selURL := fmt.Sprintf(externalSELAPTMainAPI+"&serviceKey=%s", startDate, endDate,serviceKey)
	kygURL := fmt.Sprintf(externalKYGAPTMainAPI+"&serviceKey=%s", startDate, endDate,serviceKey)
	incURL := fmt.Sprintf(externalINCAPTMainAPI+"&serviceKey=%s", startDate, endDate,serviceKey)
	
	selResp, err := http.Get(selURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer selResp.Body.Close()

	kygResp, err := http.Get(kygURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer kygResp.Body.Close()

	incResp, err := http.Get(incURL)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}
	defer incResp.Body.Close()

	selBody, err := ioutil.ReadAll(selResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	kygBody, err := ioutil.ReadAll(kygResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	incBody, err := ioutil.ReadAll(incResp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}	

	var selApiResponse ApiResponse
	if err := json.Unmarshal(selBody, &selApiResponse); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	var kygApiResponse ApiResponse
	if err := json.Unmarshal(kygBody, &kygApiResponse); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	var incApiResponse ApiResponse
	if err := json.Unmarshal(incBody, &incApiResponse); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusInternalServerError)
		return
	}

	var apiResponse ApiResponse
	apiResponse.Data = append(apiResponse.Data, selApiResponse.Data...)
	apiResponse.Data = append(apiResponse.Data, kygApiResponse.Data...)
	apiResponse.Data = append(apiResponse.Data, incApiResponse.Data...)
	sort.Slice(apiResponse.Data, func(i, j int) bool {
		return apiResponse.Data[i].RcritPblancDe > apiResponse.Data[j].RcritPblancDe
	})
	log.Printf("First external Api processing time: %s", time.Since(start))
	checkTime := time.Now()
	for i, item := range apiResponse.Data {
		detailURL := fmt.Sprintf("%s%s&serviceKey=%s", externalAPTDetailAPI, item.HouseManageNo, serviceKey)
		var detailResponse struct {
			Data []DetailData `json:"data"`
		}
		if err := fetchDetailData(detailURL, &detailResponse); err != nil {
			http.Error(w, "Failed to fetch detail data", http.StatusInternalServerError)
			return
		}
		for _, detail := range detailResponse.Data {
			apiResponse.Data[i].DetailData = append(apiResponse.Data[i].DetailData, processDetailData(&apiResponse.Data[i], detail))
		}
	}
	log.Printf("Second external Api processing time: %s", time.Since(checkTime))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiResponse)
}

func fetchDetailData(url string, detailResponse *struct {
	Data []DetailData `json:"data"`
}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, detailResponse); err != nil {
		return err
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
					detail.EtcKYGPoint = 가점자물량
					detail.EtcKYGRandZero = 랜덤무주택물량
					detail.EtcKYGRandZeroOne = 랜덤무주택일주택물량
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
func getDateRangeForPage(page int) (string, string) {
	currentDate := time.Now()

	endDate := currentDate.AddDate(0, 0, -30*(page-1))
	startDate := endDate.AddDate(0, 0, -29)

	startDateStr := startDate.Format("2006-01-02")
	endDateStr := endDate.Format("2006-01-02")

	return startDateStr, endDateStr
}