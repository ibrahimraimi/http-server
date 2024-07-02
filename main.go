package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type IPInfo struct {
	IP       string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

func getIPInfo(r *http.Request) (*IPInfo, error) {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	resp, err := http.Get("https://ipapi.co/json/")
	if err != nil {
		return nil, fmt.Errorf("error getting IP info: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("error decoding IP info: %v", err)
	}

	city, _ := data["city"].(string)

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENWEATHER_API_KEY not set")
	}

	weatherResp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey))
	if err != nil {
		return nil, fmt.Errorf("error getting weather info: %v", err)
	}
	defer weatherResp.Body.Close()

	var weatherData map[string]interface{}
	if err := json.NewDecoder(weatherResp.Body).Decode(&weatherData); err != nil {
		return nil, fmt.Errorf("error decoding weather info: %v", err)
	}

	temp, _ := weatherData["main"].(map[string]interface{})["temp"].(float64)

	visitorName := r.URL.Query().Get("visitor_name")
	if visitorName == "" {
		visitorName = os.Getenv("NAME")
		if visitorName == "" {
			visitorName = "six-shot"
		}
	}

	visitorName = strings.Trim(visitorName, "\"")

	greeting := fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", visitorName, temp, city)

	return &IPInfo{
		IP:       clientIP,
		Location: city,
		Greeting: greeting,
	}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	ipInfo, err := getIPInfo(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ipInfo)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default port
	}

	http.HandleFunc("/api/hello", handler)
	fmt.Printf("Server listening on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
