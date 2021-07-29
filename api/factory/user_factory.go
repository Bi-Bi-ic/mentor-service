package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// Userable ...
type Userable struct {
	ID           uuid.UUID      `json:"id"`
	UserName     string         `json:"username"`
	Fullname     string         `json:"fullname,omitempty"`
	Email        string         `json:"email"`
	MailContact  string         `json:"mail_contact,omitempty"`
	Token        string         `json:"token,omitempty"`
	RefreshToken string         `json:"refresh_token,omitempty"`
	Avatar       string         `json:"avatar,omitempty"`
	Cover        string         `json:"cover,omitempty"`
	GuestMode    bool           `json:"guestmode"`
	TimeLines    []TimeLineAble `json:"time_line"`
}

// TimeLineAble ...
type TimeLineAble struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	SubTitle    string    `json:"sub_title"`
	FromTime    int64     `json:"from_time"`
	ToTime      int64     `json:"to_time"`
	Description string    `json:"description"`
}

type TimeLineFactory struct{}

func (factory TimeLineFactory) Create(timeLineEntity models.TimeLine) TimeLineAble {
	return TimeLineAble{
		ID:          timeLineEntity.ID,
		SubTitle:    timeLineEntity.SubTitle,
		Title:       timeLineEntity.Title,
		FromTime:    timeLineEntity.FromTime,
		ToTime:      timeLineEntity.ToTime,
		Description: timeLineEntity.Description,
	}
}

// UserInfoFactory ... this object create for anything what if you want about user
type UserInfoFactory struct{}

// URL_SERVER  ...
const URL_SERVER string = "https://api.huc.com.vn"

// Create is a list of 'Userable' fixed from User entity
func (factory UserInfoFactory) Create(user interface{}) Userable {
	customer := models.User{}

	err := mapstructure.Decode(user, &customer)
	if err != nil {
		panic(err)
	}

	timeLinesAble := []TimeLineAble{}

	timeLineFactory := TimeLineFactory{}

	for _, value := range customer.TimeLines {

		timeLinesAble = append(timeLinesAble, timeLineFactory.Create(value))
	}

	avatarAdvance := ""
	coverAdvance := ""

	if customer.Avatar != "" {
		avatarAdvance = URL_SERVER + customer.Avatar
	}

	if customer.Cover != "" {
		coverAdvance = URL_SERVER + customer.Cover
	}

	return Userable{
		ID:           customer.ID,
		UserName:     customer.UserName,
		Fullname:     customer.Fullname,
		Email:        customer.Email,
		Token:        customer.Token,
		RefreshToken: customer.RefreshToken,
		Avatar:       avatarAdvance,
		Cover:        coverAdvance,
		GuestMode:    customer.GuestMode,
		TimeLines:    timeLinesAble,
	}
}

// CreateDetail ...
func (factory UserInfoFactory) CreateDetail(user interface{}) Userable {
	return factory.Create(user)
}

func (factory UserInfoFactory) CreateFromModel(user models.User) Userable {

	timeLinesAble := []TimeLineAble{}

	timeLineFactory := TimeLineFactory{}

	for _, value := range user.TimeLines {

		timeLinesAble = append(timeLinesAble, timeLineFactory.Create(value))
	}

	avatarAdvance := ""
	coverAdvance := ""

	if user.Avatar != "" {
		avatarAdvance = URL_SERVER + user.Avatar
	}

	if user.Cover != "" {
		coverAdvance = URL_SERVER + user.Cover
	}

	return Userable{
		ID:           user.ID,
		UserName:     user.UserName,
		Fullname:     user.Fullname,
		Email:        user.Email,
		Token:        user.Token,
		RefreshToken: user.RefreshToken,
		Avatar:       avatarAdvance,
		Cover:        coverAdvance,
		GuestMode:    user.GuestMode,
		TimeLines:    timeLinesAble,
	}
}

func (factory UserInfoFactory) CreateFromList(userEntities []models.User) (userables []Userable) {

	userables = make([]Userable, 0)

	for _, value := range userEntities {
		userables = append(userables, factory.CreateFromModel(value))
	}

	return

}
