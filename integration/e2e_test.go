//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func waitForTCP(t *testing.T, addr string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(1 * time.Second)
	}
	t.Fatalf("TCP connection not established at %s within %s", addr, timeout)
}

func getAuthToken(t *testing.T, username, password string) string {
	authURL := "http://localhost:8080/api/auth"
	payload := map[string]string{
		"username": username,
		"password": password,
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)

	resp, err := http.Post(authURL, "application/json", bytes.NewReader(data))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "auth endpoint status code")
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)
	token, ok := result["token"].(string)
	require.True(t, ok, "token is missing in auth response")
	return token
}

func TestAuthIntegration(t *testing.T) {
	waitForTCP(t, "localhost:8080", 30*time.Second)

	username := "e2e_user"
	password := "testpassword"

	token := getAuthToken(t, username, password)
	require.NotEmpty(t, token, "expected a non-empty token")
	t.Logf("AuthIntegration: token for user %s: %s", username, token)
}

func TestInfoIntegration(t *testing.T) {
	waitForTCP(t, "localhost:8080", 30*time.Second)

	username := "e2e_user"
	password := "testpassword"

	token := getAuthToken(t, username, password)

	infoURL := "http://localhost:8080/api/info"
	req, err := http.NewRequest("GET", infoURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "info endpoint status code")
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("InfoIntegration: response for user %s: %s", username, string(body))
}

func TestBuyMerchIntegration(t *testing.T) {
	waitForTCP(t, "localhost:8080", 30*time.Second)

	username := "e2e_user"
	password := "testpassword"

	token := getAuthToken(t, username, password)

	item := "t-shirt"
	buyURL := fmt.Sprintf("http://localhost:8080/api/buy/%s", item)
	req, err := http.NewRequest("GET", buyURL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected status 200, got %d", resp.StatusCode))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("BuyMerchIntegration: response for user %s: %s", username, string(body))
}

func TestTransferCoinsIntegration(t *testing.T) {
	waitForTCP(t, "localhost:8080", 30*time.Second)

	senderUsername := "e2e_sender"
	receiverUsername := "e2e_receiver"
	password := "testpassword"

	senderToken := getAuthToken(t, senderUsername, password)
	_ = getAuthToken(t, receiverUsername, password)

	transferURL := "http://localhost:8080/api/sendCoin"
	payload := map[string]interface{}{
		"toUser": receiverUsername,
		"amount": 100,
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", transferURL, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+senderToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, fmt.Sprintf("expected status 200, got %d", resp.StatusCode))
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("TransferCoinsIntegration: response from %s to %s: %s", senderUsername, receiverUsername, string(body))

	// Выполняем обратный перевод: переводим 100 монет от получателя обратно к отправителю.
	receiverToken := getAuthToken(t, receiverUsername, password)
	returnPayload := map[string]interface{}{
		"toUser": senderUsername,
		"amount": 100,
	}
	returnData, err := json.Marshal(returnPayload)
	require.NoError(t, err)

	returnReq, err := http.NewRequest("POST", transferURL, bytes.NewReader(returnData))
	require.NoError(t, err)
	returnReq.Header.Set("Content-Type", "application/json")
	returnReq.Header.Set("Authorization", "Bearer "+receiverToken)

	returnResp, err := client.Do(returnReq)
	require.NoError(t, err)
	defer returnResp.Body.Close()

	require.Equal(t, http.StatusOK, returnResp.StatusCode, fmt.Sprintf("expected status 200 on return transfer, got %d", returnResp.StatusCode))
	returnBody, err := ioutil.ReadAll(returnResp.Body)
	require.NoError(t, err)
	t.Logf("TransferCoinsIntegration: return response from %s to %s: %s", receiverUsername, senderUsername, string(returnBody))
}
