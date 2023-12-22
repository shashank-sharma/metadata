package component

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/shashank-sharma/metadata/settings"
	"github.com/shashank-sharma/metadata/util"
)

type UserRecord struct {
	Email    string `json:"email"`
	Id       string `json:"id"`
	Username string `json:"username"`
}

type AuthResponse struct {
	Record UserRecord `json:"record"`
	Token  string     `json:"token"`
}

type AuthTokenResponse struct {
	Token string `json:"token"`
}

type TrackingDeviceResponse struct {
	ProductId string `json:"id"`
}

func ResetView() {
	s := util.GetSetting()
	s.Reset()
	GenerateHomePage()
}

func HealthCheck() bool {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/api/health", util.GetConfig().ServiceHost, util.GetConfig().ServicePort), nil)
	if err != nil {
		log.Printf("Error getting token: %s", err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Request failed: %s", err)
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func PingOnlineStatus(s *settings.Settings) {
	url := fmt.Sprintf("http://%s:%s/api/track", util.GetConfig().ServiceHost, util.GetConfig().ServicePort)
	storage := util.GetStorage()
	contentType := "application/json"
	data := []byte(fmt.Sprintf(`{"userid": "%s", "token": "%s", "productid": "%s"}`, s.GetUserId(), s.GetToken(), s.GetProductId()))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		tempFp, _ := storage.FailedPing.Get()
		storage.FailedPing.Set(tempFp + 1)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Body =", string(body))
	if err != nil {
		fmt.Println(err)
		return
	}

	tempSp, _ := storage.SuccessfulPing.Get()
	storage.SuccessfulPing.Set(tempSp + 1)
	return
}

func setTrackingDevice(s *settings.Settings, token string, name string) (string, error) {
	url := fmt.Sprintf("http://%s:%s/api/track/create", util.GetConfig().ServiceHost, util.GetConfig().ServicePort)
	contentType := "application/json"
	data := []byte(fmt.Sprintf(`{"name": "%s", "hostname": "%s", "os": "%s", "arch": "%s"}`, name, s.HostName, s.OperatingSystem, s.Arch))

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Authorization", token)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("Body =", string(body))
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	trackingDeviceResponse := TrackingDeviceResponse{}
	err = json.Unmarshal(body, &trackingDeviceResponse)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return "", err
	}
	return trackingDeviceResponse.ProductId, nil
}

func GenerateHomePage() {
	s := util.GetSetting()
	w := util.GetWindow()
	box := container.NewVBox()

	box.Add(MakeToolbarTab(s, w))

	appearance := container.NewBorder(box, nil, nil, nil)
	content := container.NewVBox()
	content.Add(appearance)
	fmt.Println("Home page token =", s.GetToken())

	if s.GetToken() == "" {
		productName := widget.NewEntry()
		username := widget.NewEntry()
		password := widget.NewPasswordEntry()

		form := &widget.Form{
			Items: []*widget.FormItem{ // we can specify items in the constructor
				{Text: "Name", Widget: productName},
				{Text: "Email", Widget: username},
				{Text: "Password", Widget: password}},
			OnSubmit: func() { // optional, handle form submission
				log.Println("username:", username.Text)
				log.Println("password:", password.Text)
				params := url.Values{}
				params.Add("identity", username.Text)
				params.Add("password", password.Text)
				params.Add("userId", "1")
				resp, err := http.PostForm(fmt.Sprintf("http://%s:%s/api/collections/users/auth-with-password", util.GetConfig().ServiceHost, util.GetConfig().ServicePort),
					params)
				if err != nil {
					log.Printf("Request Failed: %s", err)
					return
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				authResponse := AuthResponse{}
				err = json.Unmarshal(body, &authResponse)
				if err != nil {
					log.Printf("Reading body failed: %s", err)
					return
				}
				fmt.Println("authResponse=", authResponse)

				req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/api/token", util.GetConfig().ServiceHost, util.GetConfig().ServicePort), nil)
				if err != nil {
					log.Printf("Error getting token: %s", err)
				}
				req.Header.Set("Authorization", authResponse.Token)
				client := &http.Client{}
				tokenResp, err := client.Do(req)
				if err != nil {
					log.Printf("Request failed: %s", err)
					return
				}
				defer tokenResp.Body.Close()
				newBody, _ := ioutil.ReadAll(tokenResp.Body)
				fmt.Println("BODY GOT =", newBody)
				authTokenResponse := AuthTokenResponse{}
				err = json.Unmarshal(newBody, &authTokenResponse)
				if err != nil {
					log.Println("Error reading body: ", err)
				}
				productId, err := setTrackingDevice(s, authResponse.Token, productName.Text)
				if err != nil {
					log.Println("Error setting tracking device: ", err)
				}
				s.SetToken(authTokenResponse.Token)
				s.SetUserId(authResponse.Record.Id)
				s.SetProductId(productId)
				log.Println("All done")
				GenerateHomePage()
			},
		}
		content.Add(form)
	} else {
		storage := util.GetStorage()
		message := widget.NewLabel("Token present")
		successfulPing := widget.NewLabelWithData(binding.IntToStringWithFormat(storage.SuccessfulPing, "Successful Ping: %d"))
		failedPing := widget.NewLabelWithData(binding.IntToStringWithFormat(storage.FailedPing, "Failed Ping: %d"))
		operatingSystem := widget.NewLabel(fmt.Sprintf("OS: %s", runtime.GOOS))
		arch := widget.NewLabel(fmt.Sprintf("Arch: %s", runtime.GOARCH))
		settingContent := container.NewVBox(successfulPing, failedPing, message, operatingSystem, arch)
		content.Add(settingContent)

		w.SetContent(content)
	}

	w.SetContent(content)
}

func GenerateSettingsPage(s *settings.Settings, w fyne.Window) {
	box := container.NewVBox()

	box.Add(MakeToolbarTab(s, w))

	appearance := container.NewBorder(box, nil, nil, nil)
	content := container.NewVBox()
	content.Add(appearance)

	boundHealthCheck := binding.NewBool()
	boundHealthCheck.Set(false)
	go func(bhCheck binding.Bool) {
		r := HealthCheck()
		bhCheck.Set(r)
	}(boundHealthCheck)
	apiDetails := widget.NewLabel(fmt.Sprintf("http://%s:%s/api/health", util.GetConfig().ServiceHost, util.GetConfig().ServicePort))
	apiStatus := widget.NewLabelWithData(binding.BoolToStringWithFormat(boundHealthCheck, "Status: %s"))
	hostname := widget.NewLabel(fmt.Sprintf("Hostname: %s", s.HostName))
	operatingSystem := widget.NewLabel(fmt.Sprintf("OS: %s", s.OperatingSystem))
	arch := widget.NewLabel(fmt.Sprintf("Arch: %s", s.Arch))
	settingContent := container.NewVBox(apiDetails, apiStatus, hostname, operatingSystem, arch)
	content.Add(settingContent)

	w.SetContent(content)
}
