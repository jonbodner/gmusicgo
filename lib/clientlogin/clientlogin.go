package clientlogin

import (
	"bytes"
	"fmt"
	"github.com/aleics/gmusicgo/lib/gmusicjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//Clientlogin structure
type Clientlogin struct {
	Auth   string    `json:"auth"`
	Header [5]string `json:"header"`
}

func Init() *Clientlogin { //Initialize the Clientlogin structure with default values
	c := new(Clientlogin)
	return c
}
func (c Clientlogin) GetAuth() string { //Get auth value
	return c.Auth
}
func (c *Clientlogin) SetAuth(s string) { //Modify auth value
	c.Auth = s
}
func (c *Clientlogin) SetHeader(accountType, email, passwd, service, source string) { //Modify header values
	c.Header[0] = accountType
	c.Header[1] = email
	c.Header[2] = passwd
	c.Header[3] = service
	c.Header[4] = source
}
func (c Clientlogin) GetHeader() [5]string { //Get header
	return c.Header
}
func (c Clientlogin) AccountType() string { //Get accountType of header
	return c.Header[0]
}
func (c Clientlogin) User() string { //Get user of header
	return c.Header[1]
}
func (c Clientlogin) Passwd() string { //Get password of header
	return c.Header[2]
}
func (c Clientlogin) Service() string { //Get name of the service of the header
	return c.Header[3]
}
func (c Clientlogin) Source() string { //Get the source of the header
	return c.Header[4]
}

func (c *Clientlogin) MakeRequest(header [5]string) [2]string { //Make the ClientLogin Request. Return the response code and body

	//Host that will be done the request
	hostname := "https://developers.google.com"
	resource := "/accounts/docs/AuthForInstalledApps"

	//POST values
	data := url.Values{}
	data.Add("accountType", header[0])
	data.Add("Email", header[1])
	data.Add("Passwd", header[2])
	data.Add("service", header[3])
	data.Add("source", header[4])

	//Hostname in URL structure
	u, _ := url.ParseRequestURI(hostname)
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u) //Save the URL of the request in urlStr

	cl := &http.Client{}                                                          //Initialize Client request
	r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode())) //Insert the values on the request
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")             //Add Content-type header

	resp, err := cl.Do(r) //Make the request
	if err != nil {       //Error management
		return [2]string{"404 Not Found", ""}
	}

	defer resp.Body.Close()
	fmt.Println(r)
	b, err := ioutil.ReadAll(resp.Body) //Get the body of the response
	if err != nil {                     //Error management
		return [2]string{"404 Not Found", ""}
	}

	body := string(b[:])

	response := [2]string{resp.Status, body}
	fmt.Println(response)

	//Delete "Auth=" of the string and save it on auth variable of ClientLogin structure
	c.Auth = response[1][(strings.LastIndex(response[1], "Auth="))+5 : (len(response[1]) - 1)]

	return response //Return response array (status,body)
}

func (c Clientlogin) SaveInfo(path string) bool { //Save UserInfo in one JSON file. Return true if the info was properly saved

	p := []string{path, "userinfo.json"} //Path and name of the file
	jsonpath := strings.Join(p, "")

	_, err := gmusicjson.Export(c, jsonpath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
