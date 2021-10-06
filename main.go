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

    if err := json.NewDecoder(response.Body).Decode(&schedules); err != nil {
        log.Fatal(err)
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

func CheckTime(date string, hourStart string, hourEnd string) bool {
    today := time.Now()

    tStart, err := time.Parse("15:04:05", hourStart)
    if err != nil {
        log.Fatal(err)
    }
    tEnd, err := time.Parse("15:04:05", hourEnd)
    if err != nil {
        log.Fatal(err)
    }
    tNow, err := time.Parse("15:04:05", time.Now().Format("15:04:05"))
    if err != nil {
        log.Fatal(err)
    }

    if date == today.Format("02 Jan 2006") {
        if tNow.Add(time.Minute * 30).After(tStart) == true && tNow.Before(tEnd) {
            return true
        }
    }
    return false
}

func ParseScedhule(date string, hour string) (time.Time, time.Time) {
    dateLayout := "02 Jan 2006"
    hourLayout := "15:04:05"

    dateStamp, err := time.Parse(dateLayout, date)
    if err != nil {
        log.Fatal(err)
    }

    hourStamp, err := time.Parse(hourLayout, hour)
    if err != nil {
        log.Fatal(err)
    }

    return dateStamp, hourStamp
}

func GetNextMeeting(schedules []ViconSchedule, currTimeStamp time.Time) (res ViconSchedule) {
    for i:=0; i<len(schedules); i++ {
        date := schedules[i].DisplayStartDate
        hour := schedules[i].EndTime
        dateStamp, hourStamp := ParseScedhule(date, hour)
        if dateStamp.Month() == currTimeStamp.Month() {
            if dateStamp.Day() == currTimeStamp.Day(){
                if hourStamp.Hour() >= currTimeStamp.Hour() {
                    res = schedules[i]
                    return
                }
            } else if dateStamp.Day() > currTimeStamp.Day() {
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
        getData := true
        isValid := false
        openBrowser := true
        for {
            currTime := time.Now()
            currLayout := "Mon, 02 Jan 2006 15:04 WIB"
            currTimeStampString := currTime.Format(currLayout)
            currTimeStamp, err := time.Parse(currLayout, currTimeStampString)
            // min := currTimeStamp.Minute()
            if err != nil {
                log.Fatal(err)
            }
            fmt.Fprintf(writer, "Current Time:\n%v\n", currTimeStampString)

            if isValid == false && getData == false && openBrowser == false {
                getData = true
                openBrowser = true
            }

            if getData == true {
                fmt.Fprintf(writer, "Fetching schedule...\n")
                schedules = GetViconSchedule(client, authResponse)
                schedule = GetNextMeeting(schedules, currTimeStamp)
                getData = false
            }

            if schedule != (ViconSchedule{}) {
                fmt.Fprintf(writer, "\nNext Class:\n")
                fmt.Fprintf(writer, "%v\n", schedule.DisplayStartDate)
                fmt.Fprintf(writer, "(%v - %v)\n", schedule.StartTime, schedule.EndTime)
                fmt.Fprintf(writer, "%v\n", schedule.CourseCode)
                fmt.Fprintf(writer, "%v\n", schedule.CourseTitleEn)
                fmt.Fprintf(writer, "%v\n", schedule.SsrComponentDescription)
                fmt.Fprintf(writer, "%v\n", schedule.ClassCode)
                fmt.Fprintf(writer, "%v\n", schedule.MeetingURL)

                isValid = CheckTime(schedule.DisplayStartDate, schedule.StartTime, schedule.EndTime)
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

