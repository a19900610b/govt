/*
govt is a VirusTotal API v2 client written for Google Go.

Written by Willi Ballenthin while at Mandiant.
June, 2013.
*/
package govt

import "log"
import "net/url"
import "net/http"
import "encoding/json"

// Client interacts with the services provided by VirusTotal.
type Client struct {
	Apikey string // private API key
	Url    string // VT URL, probably ends with .../v2/. Must end in '/'.
}

// Scan is defined by VT.
type Scan struct {
	Detected bool   `json:"detected"`
	Version  string `json:"version"`
	Result   string `json:"result"`
	Update   string `json:"update"`
}

// Report is defined by VT.
type Report struct {
	ResponseCode int             `json:"response_code"`
	VerboseMsg   string          `json:"verbose_msg"`
	Resource     string          `json:"resource"`
	ScanId       string          `json:"scan_id"`
	Md5          string          `json:"md5"`
	Sha1         string          `json:"sha1"`
	Sha256       string          `json:"sha256"`
	ScanDate     string          `json:"scan_date"`
	Positives    uint16          `json:"positives"`
	Total        uint16          `json:"total"`
	Scans        map[string]Scan `json:"scans"`
	Permalink    string          `json:"permalink"`
}

// ClientError is a generic error specific to the `govt` package.
type ClientError struct {
	msg string
}

// Error returns a string representation of the error condition.
func (self ClientError) Error() string {
	return self.msg
}

// UseDefaultUrl configures a `Client` to use the default public
//   VT URL published on their website.
func (self *Client) UseDefaultUrl() {
	self.Url = "https://www.virustotal.com/vtapi/v2/"
}

// checkApiKey ensures that the user configured her API key,
//   or returns an error.
func (self *Client) checkApiKey() (err error) {
	if self.Apikey == "" {
		return ClientError{msg: "Empty API key is invalid"}
	} else {
		return nil
	}
}

// GetReport fetches the AV scan reports tracked by VT given an MD5 hash value.
func (self *Client) GetReport(md5 string) (r *Report, err error) {
	if err = self.checkApiKey(); err != nil {
		log.Println("Invalid API Key: ", err.Error())
		return &Report{}, err
	}

	var fullurl string = self.Url + "file/report?"
	r = &Report{}

	values := url.Values{}
	values.Set("apikey", self.Apikey)
	values.Add("resource", md5)

	resp, err := http.Get(fullurl + values.Encode())
	if err != nil {
		log.Println("Failed to get ", fullurl+values.Encode(), ": ", err.Error())
		return &Report{}, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	if err = dec.Decode(r); err != nil {
		log.Println("Failed to parse response: ", err.Error())
		return &Report{}, err
	}

	return r, nil
}
