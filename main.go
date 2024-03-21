package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var checksumKey = "c548f46e2e5efa987944e750277c515532cfd13a8843a08a78c1f9cef556121c"

type Data struct {
	Code      string      `json:"code"`
	Desc      string      `json:"desc"`
	Data      DataPayload `json:"data"`
	Signature string      `json:"signature"`
}

type DataPayload struct {
	OrderCode              int    `json:"orderCode"`
	Amount                 int    `json:"amount"`
	Description            string `json:"description"`
	AccountNumber          string `json:"accountNumber"`
	Reference              string `json:"reference"`
	TransactionDateTime    string `json:"transactionDateTime"`
	Currency               string `json:"currency"`
	PaymentLinkId          string `json:"paymentLinkId"`
	Code                   string `json:"code"`
	Desc                   string `json:"desc"`
	CounterAccountBankId   string `json:"counterAccountBankId"`
	CounterAccountBankName string `json:"counterAccountBankName"`
	CounterAccountName     string `json:"counterAccountName"`
	CounterAccountNumber   string `json:"counterAccountNumber"`
	VirtualAccountName     string `json:"virtualAccountName"`
	VirtualAccountNumber   string `json:"virtualAccountNumber"`
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)
	fmt.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" || r.URL.Path != "/webhook" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	var webhookData Data
	if err := json.Unmarshal(body, &webhookData); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error unmarshalling JSON: %v", err)
		return
	}

	isValid := isValidData(webhookData)
	fmt.Println("Is valid:", isValid)

	// Your further processing logic here
	// For example, you can access webhookData.Data fields like webhookData.Data.OrderCode, webhookData.Data.Amount, etc.

	// Respond to the webhook request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}

func isValidData(data Data) bool {
	dataQueryStr := convertObjToQueryStr(data.Data)
	dataToSignature := hmac.New(sha256.New, []byte(checksumKey))
	dataToSignature.Write([]byte(dataQueryStr))
	expectedSignature := hex.EncodeToString(dataToSignature.Sum(nil))
	return expectedSignature == data.Signature
}

func convertObjToQueryStr(obj DataPayload) string {
	var queryStrings []string

	objMap := map[string]interface{}{
		"orderCode":              obj.OrderCode,
		"amount":                 obj.Amount,
		"description":            obj.Description,
		"accountNumber":          obj.AccountNumber,
		"reference":              obj.Reference,
		"transactionDateTime":    obj.TransactionDateTime,
		"currency":               obj.Currency,
		"paymentLinkId":          obj.PaymentLinkId,
		"code":                   obj.Code,
		"desc":                   obj.Desc,
		"counterAccountBankId":   obj.CounterAccountBankId,
		"counterAccountBankName": obj.CounterAccountBankName,
		"counterAccountName":     obj.CounterAccountName,
		"counterAccountNumber":   obj.CounterAccountNumber,
		"virtualAccountName":     obj.VirtualAccountName,
		"virtualAccountNumber":   obj.VirtualAccountNumber,
	}

	for key, value := range objMap {
		switch value := value.(type) {
		case int:
			queryStrings = append(queryStrings, fmt.Sprintf("%s=%d", key, value))
		case float64:
			queryStrings = append(queryStrings, fmt.Sprintf("%s=%f", key, value))
		case bool:
			queryStrings = append(queryStrings, fmt.Sprintf("%s=%t", key, value))
		case string:
			queryStrings = append(queryStrings, fmt.Sprintf("%s=%s", key, value))
		}
	}

	return strings.Join(queryStrings, "&")
}
