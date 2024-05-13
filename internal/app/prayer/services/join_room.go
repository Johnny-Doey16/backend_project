package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/utils"
)

func JoinRoomInSignallingServer(roomId, roomName string) (string, error) {
	// Create an instance of HTTPRequest with TomTom API base URL and headers
	headers := map[string]string{}
	request := utils.NewHTTPRequest(constants.WebRtcSignallingServerUrl, headers)

	// Make the GET request to the signalling server
	endpoint := fmt.Sprintf("room/%s", roomId)
	response, statusCode, err := request.Get(endpoint)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "", err
	}

	// Handle different status codes
	switch statusCode {
	case http.StatusOK:
		// Parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return "", err
		}

		roomUrl, ok := result["RoomWebsocketAddr"].(string)
		if !ok {
			return "", errors.New("error: roomUrl is not a string")
		}
		return roomUrl, nil
	case http.StatusBadRequest:
		// Handle Bad Request (400) error
		return "", errors.New("Bad request: " + string(response))
	case http.StatusNotFound:
		// TODO: Revisit
		log.Println("Inside 404!!")
		_, err = CreateRoomInSignallingServer(roomId, roomName)
		if err != nil {
			return "", err
		}

		// TODO: Fix
		// ws := "ws"
		// if os.Getenv("ENVIRONMENT") == "PRODUCTION" {
		// 	ws = "wss"
		// }
		ws := "ws"
		roomUrl := fmt.Sprintf("%s://%s/room/%s/websocket", ws, constants.WebRtcSignallingServerHost, roomId)

		return roomUrl, nil
	default:
		// Handle other status codes
		return "", errors.New("Unexpected status code: " + strconv.Itoa(statusCode))
	}
}

// func JoinRoomInSignallingServer(roomId, roomName string) (string, error) {
// 	// Create an instance of HTTPRequest with TomTom API base URL and headers
// 	headers := map[string]string{}
// 	request := utils.NewHTTPRequest(constants.WebRtcSignallingServerUrl, headers)

// 	// Make the GET request to the signalling server
// 	endpoint := fmt.Sprintf("room/%s", roomId)
// 	response, err := request.Get(endpoint)
// 	if err != nil {
// 		fmt.Println("Error making the request:", err)
// 		return "", err
// 	}

// 	// Parse the JSON response
// 	var result map[string]interface{}
// 	err = json.Unmarshal(response, &result)
// 	if err != nil {
// 		fmt.Println("Error parsing JSON response:", err)
// 		return "", err
// 	}

// 	// "RoomWebsocketAddr":   fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Hostname(), uuid),
// 	// 	"RoomName":            w.Rooms[uuid].Name,

// 	roomUrl, ok := result["RoomWebsocketAddr"].(string)
// 	if !ok {
// 		return "", errors.New("error: roomUrl is not a string")
// 	}

// 	return roomUrl, nil
// }
