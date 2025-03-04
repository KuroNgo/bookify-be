package google_calender

import (
	"bookify/pkg/interface/google_calender/handle"
	"context"
	"fmt"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
)

func SetupCalendar() {
	client := handle.GetClient()

	srv, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	events, err := srv.Events.List("primary").MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve events: %v", err)
	}

	fmt.Printf("Found %d events\n", len(events.Items))
	if len(events.Items) > 0 {
		fmt.Println("Found events:")
	} else {
		for _, item := range events.Items {
			fmt.Printf("- %s\n", item.Summary)
		}
	}

}
