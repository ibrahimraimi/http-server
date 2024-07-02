# IP Address and Greeting Service (Go)

This Go program provides a basic HTTP server that does the following:

1. **Retrieves the client's IP address.**
2. **Returns a JSON response** containing:
   * The client's IP address
   * Get the current client's IP address
   * A personalized greeting with the client name and temperature

### How to Use

1. **Prerequisites:**
   * Go installed on your system ([https://golang.org/](https://golang.org/))

2. **Run the Server:**
   * Open your terminal.
   * Navigate to the directory where this code is saved.
   * Execute:  `go run main.go`

3. **Access the Service:**
   * Open your web browser or use a tool like `curl`.
   * Visit: `http://localhost:8000`
   * You'll see a JSON response similar to this:

   ```json
   {"ip": ${client's ip},"location":${client's location},"greeting":"Hello, ${client's name}, The temperature is 1${client's temperature} degrees Celsius in ${client's location}"}

### How to Build and Run with Docker

* **Build:** Open a terminal in that directory and run:

```sh
docker build -t ip-greeter-image .

# This creates a Docker image named ip-greeter-image.
```

- **Run**

```sh
docker run -p 8000:8000 ip-greeter-image

# This starts a container from the image, publishing port 8000 inside the container to port 8000 on your local machine.
```

```go
package main

import (
 "encoding/json"
 "fmt"
 "io"
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

func getIPInfo(visitorName string) (*IPInfo, error) {
 resp, err := http.Get("https://ipapi.co/json/")
 if err != nil {
  return nil, err
 }
 defer resp.Body.Close()

 body, err := io.ReadAll(resp.Body)
 if err != nil {
  return nil, err
 }

 var data map[string]interface{}
 if err := json.Unmarshal(body, &data); err != nil {
  return nil, err
 }

 city, ok := data["city"].(string)
 if !ok || city == "" {
  return nil, fmt.Errorf("city not found in IP API response")
 }

 apiKey := os.Getenv("OPENWEATHER_API_KEY")
 weatherResp, err := http.Get(fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey))
 if err != nil {
  return nil, err
 }
 defer weatherResp.Body.Close()

 weatherBody, err := io.ReadAll(weatherResp.Body)
 if err != nil {
  return nil, err
 }
 var weatherData map[string]interface{}
 if err := json.Unmarshal(weatherBody, &weatherData); err != nil {
  return nil, err
 }

 temperature, ok := weatherData["main"].(map[string]interface{})["temp"].(float64)
 if !ok {
  return nil, fmt.Errorf("temperature not found in weather API response")
 }

 name := os.Getenv("NAME")
 if name == "" {
  name = "six-shot"
 }

 if visitorName != "" {
  visitorName = strings.Trim(visitorName, "\"")
  name = visitorName
 }

 greeting := fmt.Sprintf("Hello, %s! The temperature is %.1f degrees Celsius in %s", name, temperature, city)

 info := &IPInfo{
  IP:       data["ip"].(string),
  Location: city,
  Greeting: greeting,
 }

 return info, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
 visitorName := r.URL.Query().Get("visitor_name")
 ipInfo, err := getIPInfo(visitorName)
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
  fmt.Println("Error loading .env file")
 }

 http.HandleFunc("/api/hello", handler)
 fmt.Println("Server listening on port 8000")
 if err := http.ListenAndServe(":8000", nil); err != nil {
  panic(err)
 }
}
```