package main

type Teacher struct {
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Course struct {
	Id        int      `json:"id"`
	User_id   string   `json:"user_id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	Img       string   `json:"img"`
	Desc      string   `json:"desc"`
	Date      string   `json:"date"`
	Timestamp string   `json:"timestamp"`
	Teacher   Teacher  `json:"teacher"`
}
