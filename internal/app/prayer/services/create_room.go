package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jaevor/go-nanoid"
	"github.com/steve-mir/diivix_backend/constants"
	"github.com/steve-mir/diivix_backend/utils"
)

func GenerateRoomId() (string, error) {
	// Generate a 11-character string.
	decenaryID, err := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyz", 11)
	if err != nil {
		return "", err
	}

	id := decenaryID()
	// Insert hyphens to match the pattern "abcs-def-hij".
	formattedID := fmt.Sprintf("%s-%s-%s", id[0:4], id[4:7], id[7:11])
	return formattedID, nil
}

func CreateRoomInSignallingServer(roomId, room_name string) (string, error) {
	// Create an instance of HTTPRequest with TomTom API base URL and headers
	headers := map[string]string{}
	request := utils.NewHTTPRequest(constants.WebRtcSignallingServerUrl, headers)

	var endpoint string
	if roomId == "" {
		endpoint = fmt.Sprintf("room/create?name=%s", room_name)
	} else {
		log.Println("Room id passed")
		endpoint = fmt.Sprintf("room/create_with_id?uuid=%s&name=%s", roomId, room_name)
	}

	// Make the GET request to the signalling server
	response, statusCode, err := request.Get(endpoint)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "", err
	}

	switch statusCode {
	case http.StatusOK:
		//! Parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return "", err
		}

		// fmt.Printf("Response %+v", result)
		roomId, ok := result["roomId"].(string)
		if !ok {
			return "", errors.New("error: roomId is not a string")
		}

		return roomId, nil
	case http.StatusBadRequest:
		// Handle Bad Request (400) error
		return "", errors.New("Bad request: " + string(response))
	case http.StatusNotFound:
		// Handle Not Found (404) error
		return "", errors.New("Not found: " + string(response))
	default:
		// Handle other status codes
		return "", errors.New("Unexpected status code: " + strconv.Itoa(statusCode))
	}
}
