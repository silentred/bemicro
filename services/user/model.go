package user

// User for mysql
type User struct {
	UID       uint64 `orm:"auto;pk;column(uid)"`
	Coin      int
	Recharge  uint64
	Earn      uint64
	ReplyCost uint64
	Email     string
	Username  string
	Password  string
	Avatar    string
	Cover     string

	FollowerCount  uint64
	FollowingCount uint32
}

func (user *User) TableName() string {
	return "snap_user_info"
}
