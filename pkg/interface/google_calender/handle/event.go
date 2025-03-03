package handle

import (
	"context"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
)

// Hàm tạo sự kiện trên Google Calendar
func createEvent(summary, location, description, startTime, endTime string) {
	client := GetClient()
	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Không thể kết nối Google Calendar API: %v", err)
	}

	event := &calendar.Event{
		Summary:     summary,
		Location:    location,
		Description: description,
		Start: &calendar.EventDateTime{
			DateTime: startTime, // Định dạng RFC3339: "2025-03-04T10:00:00+07:00"
			TimeZone: "Asia/Ho_Chi_Minh",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: "Asia/Ho_Chi_Minh",
		},
		Attendees: []*calendar.EventAttendee{
			{Email: "user@example.com"}, // Người tham gia (email)
		},
		Reminders: &calendar.EventReminders{
			UseDefault: false,
			Overrides: []*calendar.EventReminder{
				{Method: "email", Minutes: 30}, // Nhắc trước 30 phút qua email
				{Method: "popup", Minutes: 10}, // Nhắc trước 10 phút qua popup
			},
		},
	}

	event, err = srv.Events.Insert("primary", event).Do()
	if err != nil {
		log.Fatalf("Không thể tạo sự kiện: %v", err)
	}

	fmt.Printf("Sự kiện đã tạo: %s\n", event.HtmlLink)
}
