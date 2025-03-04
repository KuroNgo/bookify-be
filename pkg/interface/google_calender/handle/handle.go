package handle

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"log"
	"net/http"
	"os"
)

func GetClient() *http.Client {
	// Đọc token.json
	b, err := os.ReadFile("pkg/interface/google_calender/config/token.json")
	if err != nil {
		log.Fatalf("Không thể đọc file credentials: %v", err)
	}

	// Load OAuth2 Config
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Không thể parse credentials: %v", err)
	}

	// Kiểm tra token đã có chưa
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

// Lấy token từ web nếu chưa có
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Truy cập link này để xác thực: %v\n", authURL)

	var authCode string
	fmt.Print("Nhập mã xác thực: ")
	fmt.Scan(&authCode)

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Không thể lấy token: %v", err)
	}
	return tok
}

// Lưu token vào file
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Lưu token vào %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Không thể tạo file token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Đọc token từ file
func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
