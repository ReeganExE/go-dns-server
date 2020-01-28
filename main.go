package main

func main() {
	server := NewDNSServer(53, "127.0.0.1", []string{"abc.com"})
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
