package moderate

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func ModerateWithAI(content string) (bool, error) {
	apikey := os.Getenv("GEMINI_KEY")
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash-lite:generateContent?key=" + apikey

	payload := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": map[string]string{"text": "Strictly moderate: profanity, racism, nazism, threats, 18+, cyber-threats, bad content. If safe, output 'OK'. If unsafe, output 'BAD'. No other words."},
		},
		"contents": []interface{}{
			map[string]interface{}{
				"parts": []interface{}{
					map[string]string{"text": content},
				},
			},
		},
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	bodyBytes := new(bytes.Buffer)
	bodyBytes.ReadFrom(resp.Body)

	return !bytes.Contains(bodyBytes.Bytes(), []byte("BAD")), nil
}
