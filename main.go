package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/gosuri/uilive"
)

var (
	baseURL = "https://myclass.apps.binus.ac.id"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		log.Println("Your platform is unsupported!")
	}
}

func Login(client *http.Client, user User) AuthResponse {
	var authResponse AuthResponse

	loginURL := baseURL + "/Auth/Login"
	data := url.Values{
		"Username": {user.Username},
		"Password": {user.Password},
	}
	payload := bytes.NewBufferString(data.Encode())

	request, err := http.NewRequest("POST", loginURL, payload)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&authResponse); err != nil {
		log.Fatal(err)
	}
	return authResponse
}

func GetViconSchedule(client *http.Client, authResponse AuthResponse) []ViconSchedule {
	var schedules []ViconSchedule

	url := baseURL + "/Home/GetViconSchedule"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	for i := 0; i < 3; i++ {
		if err := json.NewDecoder(response.Body).Decode(&schedules); err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	return schedules
}

func OpenInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func InTimeSpan(start, end time.Time) bool {
	check, err := time.Parse("15:04:05", time.Now().Format("15:04:05"))
	if err != nil {
		log.Fatal(err)
	}

	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}

func ParseDate(date string) time.Time {
	dateLayout := "02 Jan 2006"

	dateStamp, err := time.Parse(dateLayout, date)
	if err != nil {
		log.Fatal(err)
	}

	return dateStamp
}

func ParseHour(hour string) time.Time {
	hourLayout := "15:04:05"

	hourStamp, err := time.Parse(hourLayout, hour)
	if err != nil {
		log.Fatal(err)
	}

	return hourStamp
}

func GetNextMeeting(schedules []ViconSchedule, currTimeStamp time.Time) (res ViconSchedule) {
	for i := 0; i < len(schedules); i++ {
		classDate := schedules[i].DisplayStartDate
		classDateStamp := ParseDate(classDate)

		endHour := schedules[i].StartTime
		endHourStamp := ParseHour(endHour)

		if classDateStamp.Month() == currTimeStamp.Month() {
			if classDateStamp.Day() == currTimeStamp.Day() {
				if endHourStamp.Hour() >= currTimeStamp.Hour() {
					res = schedules[i]
					return
				}
			} else if classDateStamp.Day() > currTimeStamp.Day() {
				res = schedules[i]
				return
			}
		}
	}
	return
}

func main() {
	var schedules []ViconSchedule
	var schedule ViconSchedule
	var user User

	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%s/anti-sp", configPath)
	fpath := path + "/credential.json"

	if CheckCredentials(user, path, fpath) == false {
		return
	}

	file, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.Unmarshal([]byte(file), &user); err != nil {
		log.Fatal(err)
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
	}

	CallClear()

	fmt.Println("Logging in...")
	authResponse := Login(client, user)

	if authResponse.Status == true {
		fmt.Println("Login success...")

		writer := uilive.New()
		writer.Start()
		isValid := false
		openBrowser := true
		fmt.Fprintf(writer, "Fetching schedule...\n")
		schedules = GetViconSchedule(client, authResponse)
		for {
			currTime := time.Now()
			currLayout := "Mon, 02 Jan 2006 15:04 WIB"
			currTimeStampString := currTime.Format(currLayout)
			currTimeStamp, err := time.Parse(currLayout, currTimeStampString)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(writer, "Current Time:\n%v\n", currTimeStampString)

			if isValid == false && openBrowser == false {
				openBrowser = true
			}

			schedule = GetNextMeeting(schedules, currTimeStamp)
			if schedule != (ViconSchedule{}) {
				fmt.Fprintf(writer, "\nNext Class:\n")
				fmt.Fprintf(writer, "%v\n", schedule.DisplayStartDate)
				fmt.Fprintf(writer, "(%v - %v)\n", schedule.StartTime, schedule.EndTime)
				fmt.Fprintf(writer, "%v\n", schedule.CourseCode)
				fmt.Fprintf(writer, "%v\n", schedule.CourseTitleEn)
				fmt.Fprintf(writer, "%v\n", schedule.SsrComponentDescription)
				fmt.Fprintf(writer, "%v\n", schedule.ClassCode)
				fmt.Fprintf(writer, "%v\n", schedule.MeetingURL)

				start := ParseHour(schedule.StartTime)
				end := ParseHour(schedule.EndTime)
				isValid = InTimeSpan(start, end)

				if isValid == true && openBrowser == true {
					fmt.Println("Opening link in browser...")
					OpenInBrowser(schedule.MeetingURL)
					openBrowser = false
				}
			}
			time.Sleep(1 * time.Minute)
		}
	} else {
		fmt.Println(authResponse)
		if err := os.Remove(fpath); err != nil {
			log.Fatal(err)
		}
	}
}
