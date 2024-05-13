package fbservice

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"

	// "github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

func getDecodedFireBaseKey() ([]byte, error) {

	fireBaseAuthKey := os.Getenv("FIREBASE_AUTH_KEY")

	decodedKey, err := base64.StdEncoding.DecodeString(fireBaseAuthKey)
	if err != nil {
		return nil, err
	}

	return decodedKey, nil
}

func SendPushNotificationSingle(deviceToken, title, body, imgUrl string, data map[string]string) error {

	opts := option.WithCredentialsFile("diivix1-firebase-adminsdk-vu9c3-a25d50f8c8.json")
	app, err := firebase.NewApp(context.Background(), nil, opts)

	if err != nil {
		log.Printf("Error in initializing firebase app" + err.Error())
		return err
	}

	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return err
	}

	response, err := fcmClient.Send(context.Background(), &messaging.Message{

		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: imgUrl,
		},
		Token: deviceToken, // it's a single device token
		Data:  data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				Title:       title,
				Body:        body,
				ImageURL:    imgUrl,
				ChannelID:   "diivix_notification_id",
			},
		},
	})

	if err != nil {
		return err
	}

	log.Println(response)

	return nil
}

func SendPushNotificationMulti(deviceTokens []string, title, body, imgUrl string, data map[string]string) error {
	fmt.Println("Tokens", deviceTokens)
	opts := option.WithCredentialsFile("diivix1-firebase-adminsdk-vu9c3-a25d50f8c8.json")
	app, err := firebase.NewApp(context.Background(), nil, opts)

	if err != nil {
		log.Printf("Error in initializing firebase app" + err.Error())
		return err
	}

	fcmClient, err := app.Messaging(context.Background())

	if err != nil {
		return err
	}

	response, err := fcmClient.SendMulticast(context.Background(), &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title:    title, //"Congratulations!!",
			Body:     body,  //"You have just implement push notification",
			ImageURL: imgUrl,
		},
		Tokens: deviceTokens, // it's an array of device tokens
		Data:   data,
		Android: &messaging.AndroidConfig{
			Notification: &messaging.AndroidNotification{
				ClickAction: "FLUTTER_NOTIFICATION_CLICK",
				Title:       title,
				Body:        body,
				ImageURL:    imgUrl,
				ChannelID:   "diivix_notification_id",
			},
		},
	})

	if err != nil {
		return err
	}

	log.Println("Response success count : ", response.SuccessCount)
	log.Println("Response failure count : ", response.FailureCount)
	return nil
}
