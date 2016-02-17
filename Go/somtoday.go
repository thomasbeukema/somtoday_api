package somtoday

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Somtoday object
type Somtoday struct {
	PupilID  string
	Username string
	Password string
	School   string
	Brin     string
	BaseURL  string
	FullName string
}

type loginResponse struct {
	Leerlingen []map[string]interface{}
	Verzorger  bool
	Error      string
}

type somResponse struct {
	Data []map[string]interface{}
}

// SetCredentials : Sets the credentials of the Somtoday object
func (som *Somtoday) SetCredentials(credentials [4]string) {
	som.Username = credentials[0]
	som.Password = credentials[1]
	som.School = credentials[2]
	som.Brin = credentials[3]
}

// Login on somtoday servers
func (som *Somtoday) Login() error {
	username := som.Username
	password := som.Password
	school := som.School
	brin := som.Brin

	username = base64.StdEncoding.EncodeToString([]byte(username))

	hasher := sha1.New()
	hasher.Write([]byte(password))
	password = fmt.Sprintf("%s", hasher.Sum(nil))
	password = base64.StdEncoding.EncodeToString([]byte(password))
	password = hex.EncodeToString([]byte(password))

	loginURL := "http://somtoday.nl/" + school + "/services/mobile/v10/Login/CheckMultiLoginB64/" + username + "/" + password + "/" + brin

	response, err := http.Get(loginURL)

	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(response.Body)
	r := fmt.Sprintf("%s", body)
	defer response.Body.Close()
	var jsonContent loginResponse
	err = json.Unmarshal([]byte(r), &jsonContent)
	if err != nil {
		fmt.Printf("Error parsing JSON at login: %s\n", err)
		os.Exit(1)
	}

	if jsonContent.Error == "SUCCESS" {
		som.Username = username
		som.Password = password
		som.School = school
		som.Brin = brin
		pupilID := fmt.Sprintf("%.0f", jsonContent.Leerlingen[0]["leerlingId"].(float64))
		som.PupilID = pupilID
		som.BaseURL = "http://somtoday.nl/" + school + "/services/mobile/v10/"
		som.FullName = jsonContent.Leerlingen[0]["fullName"].(string)

		return nil
	} else if jsonContent.Error == "FEATURE_NOT_ACTIVATED" {
		return errors.New("FEATURE_NOT_ACTIVATED")
	} else if jsonContent.Error == "FAILED_AUTENTICATION" {
		return errors.New("FAILED_AUTENTICATION")
	} else if jsonContent.Error == "FAILED_OTHER_TYPE" {
		return errors.New("FAILED_OTHER_TYPE")
	} else {
		return errors.New("SOMTODAY_ERROR_UNKNOWN")
	}
}

// GetTimetable returns timetable
func (som *Somtoday) GetTimetable(days int64) (*somResponse, error) {
	var daysToGo []byte
	daysToGo = strconv.AppendInt(daysToGo, (int64(time.Now().Unix()+(days*86400)) * 1000), 10)

	URL := som.BaseURL + "Agenda/GetMultiStudentAgendaB64/" + som.Username + "/" + som.Password + "/" + som.Brin + "/" + string(daysToGo) + "/" + som.PupilID
	response, err := http.Get(URL)
	if err != nil {
		return new(somResponse), errors.New("SOMTODAY_URL_FETCH_ERROR")
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonObj somResponse
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		fmt.Printf("%s", err)
		return new(somResponse), nil
	}
	return &jsonObj, nil
}

// GetHW returns the homework for a given period (days int)
func (som *Somtoday) GetHW(days int) (*somResponse, error) {
	URL := som.BaseURL + "Agenda/GetMultiStudentAgendaHuiswerkMetMaxB64/" + som.Username + "/" + som.Password + "/" + som.Brin + "/" + strconv.Itoa(days) + "/" + som.PupilID
	response, err := http.Get(URL)
	if err != nil {
		return nil, errors.New("SOMTODAY_URL_FETCH_ERROR")
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonObj somResponse
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, nil
	}

	return &jsonObj, nil
}

// GetGrades returns grades (recent/all, depending on mode (0/1))
func (som *Somtoday) GetGrades(mode int) (*somResponse, error) {
	var URL string
	if mode == 0 {
		URL = som.BaseURL + "Cijfers/GetMultiCijfersRecentB64/" + som.Username + "/" + som.Password + "/" + som.Brin + "/" + som.PupilID
	} else {
		URL = som.BaseURL + "Cijfers/GetCijfersRecentMetMaxB64/" + som.Username + "/" + som.Password + "/" + som.Brin + "/1000"
	}

	response, err := http.Get(URL)
	if err != nil {
		return nil, errors.New("SOMTODAY_URL_FETCH_ERROR")
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var jsonObj somResponse
	err = json.Unmarshal(body, &jsonObj)
	if err != nil {
		fmt.Printf("%s", err)
		return nil, nil
	}

	return &jsonObj, nil
}

// ChangeHwDone changes the status of the homework, appointmentID and hwID could be retrieved from method GetHW
func (som *Somtoday) ChangeHwDone(appointmentID, hwID string, done bool) bool {
	URL := som.BaseURL + "Agenda/HuiswerkDoneB64/" + som.Username + "/" + som.Password + "/" + som.Brin + "/" + appointmentID + "/" + hwID

	if done {
		URL += "/1"
	} else {
		URL += "/0"
	}

	fmt.Println(URL)

	response, err := http.Get(URL)
	if err != nil {
		return false
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if strings.Contains(string(body), "\"status\":\"OK\"") {
		return true
	}

	return false
}

func (som *Somtoday) String() string {
	username := som.Username
	password := som.Password
	school := som.School
	brin := som.Brin
	pupilID := som.PupilID
	fn := som.FullName
	return fmt.Sprintf("Username: %s\nPassword: %s\nSchool: %s\nBRIN: %s\npupilID: %s\nName: %s\n", username, password, school, brin, pupilID, fn)
}
