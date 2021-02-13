package main

/*
var client = &http.Client{
	Timeout: time.Second * 10,
}
*/

/*
func sendReq(){
	req, err := http.NewRequest("GET", "http://192.168.1.6:8080/admin/api.php", nil)

	if err != nil {
		log.Fatalf("Error generating get repo request: %+v", err)
	}

	response, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %+v", err.Error())
	}
	defer response.Body.Close()
	parsedBody, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(parsedBody))
}
*/
