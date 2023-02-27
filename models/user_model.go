package models

type User struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    UserName  string `json:"user_name"`
    Password  string `json:"password"`
}