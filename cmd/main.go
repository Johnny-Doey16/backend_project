package main

// "github.com/steve-mir/diivix_backend/utils"

// GenerateSecureToken generates a cryptographically secure random string.
// The length of the resulting string will be approximately 4/3 of the size due
// to base64 encoding.
// func GenerateSecureRandomNumber(max int64) (int64, error) {
// 	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return nBig.Int64(), nil
// }

func main() {
	/*
		maxNumber := int64(1000000) // for a 6-digit number, set max to 1 million
		randomNumber, err := utils.GenerateSecureRandomNumber(maxNumber)
		if err != nil {
			fmt.Println("Error generating secure random number:", err)
			return
		}
		print(randomNumber)
		pwdResetCodeStr := fmt.Sprintf("%06d", randomNumber)
		pwdResetCodeStr2 := strconv.FormatInt(randomNumber, 10)
		fmt.Printf("Secure random number: %06d\n", randomNumber) // prints the number with leading zeros
		println("Str is ", pwdResetCodeStr)
		println("Str2 is ", pwdResetCodeStr2)
	*/

	/*
		// Assuming the reading plan started on November 1st
		startDate := time.Date(2023, time.November, 1, 0, 0, 0, 0, time.UTC)

		// Current date is February 17th
		currentDate := time.Date(2024, time.February, 17, 0, 0, 0, 0, time.UTC)

		// Reading plan duration is 255 days
		durationInDays := 255

		// Path to the XML file containing the Bible data
		biblePath := "../bible-translations/eng-dra.osis.xml"

		// Get the current reading for the specified day
		currentReading, err := services.GetCurrentReading(biblePath, startDate, durationInDays, currentDate)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Current Reading:")
		fmt.Println(currentReading)
	*/
}
