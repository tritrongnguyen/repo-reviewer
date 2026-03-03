package handler

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     int    `json:"role:"`
}

type UserProfile struct {
	FullName  string `json:"full_name"`
	AvatarUrl string `json:"avatar_url"`
	Phone     string `json:"phone"`
}

type UserConfig struct {
	APIKey string `json:"string,omitempty"`
}

type UserDetails struct {
	User
	Profile UserProfile `json:"profile"`
	Config  UserConfig  `json:"config"`
}

// func Register(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request){

// 	}
// }
