package main

import (
	"log"
	"regexp"
	"strings"

	fbservice "github.com/steve-mir/diivix_backend/fb-service"
	// "github.com/steve-mir/diivix_backend/utils"
)

func main() {

	err := fbservice.SendPushNotificationSingle("token",
		"Title", "This is the body", "", map[string]string{"activity": "PostActivity"})
	if err != nil {
		log.Println("Error sending push notification", err.Error())
	}

	// str := detectMentions("content string @john doe will see @Peter tomoror for call in @Andrew place to talk to @Peter @JoHn")
	// log.Printf("Result %+v", str)

	// // Geocoding
	// // Use Viper for configuration management
	// config, err := utils.LoadConfig("../../")
	// if err != nil {
	// 	log.Fatal("cannot load config " + err.Error())
	// }

	// lat, lng, err := utils.TomTomGeocoding(config, "13B New Nkisi Rd", "Onitsha", "434106", "Anambra state", "Nigeria")
	// if err != nil {
	// 	log.Println("Error geocoding", err.Error())
	// }
	// log.Printf("Latitude is %f and longitude is %f", lat, lng)

	// addrs, err := utils.TomTomReverseGeocoding(config, 6.123730, 6.794020)
	// if err != nil {
	// 	log.Println("Error Reverse geocoding", err.Error())
	// }
	// log.Println("Address is", addrs)

}

func detectMentions(content string) []string {
	mentionRegex := regexp.MustCompile(`@(\w+)`)
	matches := mentionRegex.FindAllStringSubmatch(content, -1)

	usernameSet := make(map[string]struct{})
	for _, match := range matches {
		username := strings.ToLower(match[1])
		usernameSet[username] = struct{}{}
	}

	usernames := make([]string, 0, len(usernameSet))
	for username := range usernameSet {
		usernames = append(usernames, username)
	}

	if len(usernames) == 0 {
		return []string{}
	}
	return usernames
}
