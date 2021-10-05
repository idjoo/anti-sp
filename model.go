package main

type User struct {
    Username    string `json:"username"`
    Password    string `json:"password"`
}

type AuthResponse struct {
	Status  bool   `json:"Status"`
	Message string `json:"Message"`
	URL     string `json:"URL"`
}

type ViconSchedule struct {
	StartDate               string      `json:"StartDate"`
	DisplayStartDate        string      `json:"DisplayStartDate"`
	StartTime               string      `json:"StartTime"`
	EndTime                 string      `json:"EndTime"`
	SsrComponentDescription string      `json:"SsrComponentDescription"`
	ClassCode               string      `json:"ClassCode"`
	Room                    string      `json:"Room"`
	Campus                  string      `json:"Campus"`
	DeliveryMode            string      `json:"DeliveryMode"`
	CourseCode              string      `json:"CourseCode"`
	CourseTitleEn           string      `json:"CourseTitleEn"`
	ClassType               string      `json:"ClassType"`
	WeekSession             int         `json:"WeekSession"`
	CourseSessionNumber     int         `json:"CourseSessionNumber"`
	MeetingID               string      `json:"MeetingId"`
	MeetingPassword         string      `json:"MeetingPassword"`
	MeetingURL              string      `json:"MeetingUrl"`
	UserFlag                string      `json:"UserFlag"`
	BinusianID              string      `json:"BinusianId"`
	PersonCode              string      `json:"PersonCode"`
	FullName                string      `json:"FullName"`
	AcademicCareer          string      `json:"AcademicCareer"`
	ClassMeetingID          string      `json:"ClassMeetingId"`
	Location                string      `json:"Location"`
	MeetingStartDate        string      `json:"MeetingStartDate"`
	ID                      interface{} `json:"Id"`
}
