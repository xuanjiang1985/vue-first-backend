package controllers

type Users struct {
	Id         int
	Name       string
	Phone      string
	Password   string
	Header     string
	Sex        int
	Admin      int
	Created_at int
	Updated_at int
}

type Messages struct {
	Id         int
	Content    string
	Created_at string
}
