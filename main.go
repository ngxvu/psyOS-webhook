package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Định nghĩa endpoint để nhận webhook
	http.HandleFunc("/webhook", handleWebhook)

	// Khởi động máy chủ và lắng nghe cổng 3000
	fmt.Println("Server is running on http://localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// Đảm bảo rằng request được gửi đến /webhook và là phương thức POST
	if r.URL.Path != "https://api-merchant.payos.vn/confirm-webhook" || r.Method != "POST" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Đọc dữ liệu từ webhook
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error reading request body: %v", err)
		return
	}

	// Xử lý dữ liệu từ webhook (trong ví dụ này, chúng ta chỉ in ra nó)
	fmt.Println("Received webhook data:", string(body))

	// Phản hồi để xác nhận nhận được webhook thành công
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received successfully"))
}
