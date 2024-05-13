package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ? Test that new works
func TomTomGeocoding(config Config, address, city, postalCode, state, country string) (float64, float64, error) {
	apiKey := config.TomTomApiKey

	// Create an instance of HTTPRequest with TomTom API base URL and headers
	tomtomBaseURL := "https://api.tomtom.com/search/2/geocode"
	headers := map[string]string{}
	tomtomRequest := NewHTTPRequest(tomtomBaseURL, headers)

	// Construct the query string
	query := fmt.Sprintf("%s, %s, %s, %s, %s", address, city, postalCode, state, country)

	// Make the GET request to TomTom API
	endpoint := fmt.Sprintf("%s.json?key=%s", query, apiKey)
	response, statusCode, err := tomtomRequest.Get(endpoint)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return 0, 0, err
	}

	switch statusCode {
	case http.StatusOK:
		// ! Parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return 0, 0, err
		}

		var lat, lng float64
		// Extract latitude and longitude from the response
		if results, ok := result["results"].([]interface{}); ok && len(results) > 0 {
			if position, ok := results[0].(map[string]interface{})["position"].(map[string]interface{}); ok {
				lat = position["lat"].(float64)
				lng = position["lon"].(float64)
				fmt.Printf("Latitude: %f\nLongitude: %f\n", lat, lng)
			} else {
				fmt.Println("Error extracting position from the response.")
			}
		} else {
			fmt.Println("No results found.")
		}
		return lat, lng, nil

	default:
		// Handle other status codes
		return 0, 0, errors.New("Unexpected status code: " + strconv.Itoa(statusCode))
	}

}

// Add reverse geocoding
func TomTomReverseGeocoding(config Config, latitude, longitude float64) (string, error) {
	apiKey := config.TomTomApiKey
	// latitude := 40.748817
	// longitude := -73.985428

	// Create an instance of HTTPRequest with TomTom API base URL and headers
	tomtomBaseURL := "https://api.tomtom.com/search/2/reverseGeocode"
	headers := map[string]string{}
	tomtomRequest := NewHTTPRequest(tomtomBaseURL, headers)

	// Construct the query string for reverse geocoding
	query := fmt.Sprintf("%f,%f.json?key=%s", latitude, longitude, apiKey)

	// Make the GET request to TomTom API for reverse geocoding
	response, statusCode, err := tomtomRequest.Get(query)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "", err
	}
	switch statusCode {
	case http.StatusOK:
		// ! Parse the JSON response
		var result map[string]interface{}
		err = json.Unmarshal(response, &result)
		if err != nil {
			fmt.Println("Error parsing JSON response:", err)
			return "", err
		}

		var addr string
		// Extract address information from the response
		if results, ok := result["addresses"].([]interface{}); ok && len(results) > 0 {
			if address, ok := results[0].(map[string]interface{})["address"].(map[string]interface{}); ok {
				addr = address["freeformAddress"].(string)
				fmt.Printf("Formatted Address: %s\n", addr)
			} else {
				fmt.Println("Error extracting address from the response.")
			}
		} else {
			fmt.Println("No results found.")
		}

		return addr, nil
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

// ! Old
/*
func TomTomGeocoding(config Config, address, city, postalCode, state, country string) (float64, float64, error) {
	apiKey := config.TomTomApiKey

	// Create an instance of HTTPRequest with TomTom API base URL and headers
	tomtomBaseURL := "https://api.tomtom.com/search/2/geocode"
	headers := map[string]string{}
	tomtomRequest := NewHTTPRequest(tomtomBaseURL, headers)

	// Construct the query string
	query := fmt.Sprintf("%s, %s, %s, %s, %s", address, city, postalCode, state, country)

	// Make the GET request to TomTom API
	endpoint := fmt.Sprintf("%s.json?key=%s", query, apiKey)
	response, err := tomtomRequest.Get(endpoint)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return 0, 0, err
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return 0, 0, err
	}

	var lat, lng float64
	// Extract latitude and longitude from the response
	if results, ok := result["results"].([]interface{}); ok && len(results) > 0 {
		if position, ok := results[0].(map[string]interface{})["position"].(map[string]interface{}); ok {
			lat = position["lat"].(float64)
			lng = position["lon"].(float64)
			fmt.Printf("Latitude: %f\nLongitude: %f\n", lat, lng)
		} else {
			fmt.Println("Error extracting position from the response.")
		}
	} else {
		fmt.Println("No results found.")
	}
	return lat, lng, nil
}

// Add reverse geocoding
func TomTomReverseGeocoding(config Config, latitude, longitude float64) (string, error) {
	apiKey := config.TomTomApiKey
	// latitude := 40.748817
	// longitude := -73.985428

	// Create an instance of HTTPRequest with TomTom API base URL and headers
	tomtomBaseURL := "https://api.tomtom.com/search/2/reverseGeocode"
	headers := map[string]string{}
	tomtomRequest := NewHTTPRequest(tomtomBaseURL, headers)

	// Construct the query string for reverse geocoding
	query := fmt.Sprintf("%f,%f.json?key=%s", latitude, longitude, apiKey)

	// Make the GET request to TomTom API for reverse geocoding
	response, err := tomtomRequest.Get(query)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return "", err
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(response, &result)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return "", err
	}

	var addr string
	// Extract address information from the response
	if results, ok := result["addresses"].([]interface{}); ok && len(results) > 0 {
		if address, ok := results[0].(map[string]interface{})["address"].(map[string]interface{}); ok {
			addr = address["freeformAddress"].(string)
			fmt.Printf("Formatted Address: %s\n", addr)
		} else {
			fmt.Println("Error extracting address from the response.")
		}
	} else {
		fmt.Println("No results found.")
	}

	return addr, nil
}
*/
