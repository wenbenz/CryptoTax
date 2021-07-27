package nicehash

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	NICEHASH_API = "https://api2.nicehash.com"
)

type Client struct {
	Credentials Credentials
	httpClient  http.Client
}

type Credentials struct {
	Key    string
	Secret string
	OrgId  string
}

func NewClient(cred Credentials) Client {
	return Client{
		Credentials: cred,
		httpClient:  http.Client{},
	}
}

func NewClientReadFrom(path string) (Client, error) {
	credentials, err := getCredentials(path)
	if err != nil {
		return Client{}, err
	}
	return Client{
		Credentials: credentials,
		httpClient:  http.Client{},
	}, nil
}

func (c *Client) Do(method, endUrl string, queryMap map[string]string, destination interface{}) error {
	endpoint := NICEHASH_API + endUrl
	request, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		return err
	}
	c.populateHeaders(*request, queryMap, true)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	return json.NewDecoder(response.Body).Decode(destination)
}

func getCredentials(path string) (Credentials, error) {
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return Credentials{}, err
	}
	fields := strings.Split(string(fileBytes), "\n")
	return Credentials{
		Key:    fields[0],
		Secret: fields[1],
		OrgId:  fields[2],
	}, err
}

func (c *Client) populateHeaders(request http.Request, queryMap map[string]string, includeAuth bool) {
	// make random nonce
	nonceBytes := make([]byte, 18) // becomes 36 byte hex
	rand.Read(nonceBytes)
	nonce := hex.EncodeToString(nonceBytes)

	request.Header.Add("X-Time", strconv.FormatInt(time.Now().Unix()*1000, 10))
	request.Header.Add("X-Nonce", nonce)
	request.Header.Add("X-Organization-Id", c.Credentials.OrgId)
	request.Header.Add("X-Request-Id", nonce)
	request.Header.Add("Accept", "application/json")
	query := request.URL.Query()
	for key, value := range queryMap {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()
	if includeAuth {
		c.populateAuth(request)
	}
}

func (c *Client) populateAuth(request http.Request) {
	input := bytes.NewBuffer([]byte{})
	input.WriteString(c.Credentials.Key)
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Time"))
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Nonce"))
	input.WriteByte(0x00)
	input.WriteByte(0x00)
	input.WriteString(request.Header.Get("X-Organization-Id"))
	input.WriteByte(0x00)
	input.WriteByte(0x00)
	input.WriteString(request.Method)
	input.WriteByte(0x00)
	input.WriteString(request.URL.Path)
	input.WriteByte(0x00)
	input.WriteString(request.URL.RawQuery)
	// body?

	inputString := input.String()

	// hash
	mac := hmac.New(sha256.New, []byte(c.Credentials.Secret))
	mac.Write([]byte(inputString))
	request.Header.Add("X-Auth", c.Credentials.Key+":"+hex.EncodeToString(mac.Sum(nil)))
}
