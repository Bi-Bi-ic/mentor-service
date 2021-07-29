package factory

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rgrs-x/service/api/models"
	uuid "github.com/satori/go.uuid"
)

// Partnerable ...
type Partnerable struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	UserName         string    `json:"username"`
	PartnerName      string    `json:"name"`
	Token            string    `json:"token" sql:"-"`
	RefreshToken     string    `json:"refresh_token" sql:"-"`
	Avatar           string    `json:"avatar"`
	PostAvailable    string    `json:"post_available"`
	TotalLike        uint      `json:"total_like"`
	Phone            string    `json:"phone"`
	MailContact      string    `json:"mail_contact"`
	models.Address   `json:"address"`
	models.WorkSpace `json:"company"`
}

// PartnerPublishable ...
type PartnerPublishable struct {
	ID               uuid.UUID `json:"id"`
	Email            string    `json:"email"`
	UserName         string    `json:"username"`
	PartnerName      string    `json:"name"`
	Token            string    `json:"token" sql:"-"`
	RefreshToken     string    `json:"refresh_token" sql:"-"`
	Avatar           string    `json:"avatar"`
	PostAvailable    string    `json:"post_available"`
	TotalLike        uint      `json:"total_like"`
	Phone            string    `json:"phone"`
	MailContact      string    `json:"mail_contact"`
	models.Address   `json:"address"`
	models.WorkSpace `json:"company"`
}

// PartnerInfoFactory this object create for anything what if you want about partner
type PartnerInfoFactory struct{}

// Create is a list of 'Partnerable' fixed from Partner entity
func (factory PartnerInfoFactory) Create(partner interface{}) Partnerable {
	customer := models.Partner{}

	err := mapstructure.Decode(partner, &customer)
	if err != nil {
		panic(err)
	}

	return Partnerable{
		ID:            customer.ID,
		Email:         customer.Email,
		UserName:      customer.UserName,
		PartnerName:   customer.PartnerName,
		Token:         customer.Token,
		RefreshToken:  customer.RefreshToken,
		Avatar:        URL_SERVER + customer.Avatar,
		PostAvailable: customer.PostAvailable,
		TotalLike:     customer.TotalLike,
		Phone:         customer.Phone,
		MailContact:   customer.MailContact,
		Address:       customer.Address,
		WorkSpace:     customer.WorkSpace,
	}
}

func (factory PartnerInfoFactory) CreateFromModel(partner models.Partner) PartnerPublishable {

	return PartnerPublishable{
		ID:            partner.ID,
		Email:         partner.Email,
		UserName:      partner.UserName,
		PartnerName:   partner.PartnerName,
		Token:         partner.Token,
		RefreshToken:  partner.RefreshToken,
		Avatar:        URL_SERVER + partner.Avatar,
		PostAvailable: partner.PostAvailable,
		TotalLike:     partner.TotalLike,
		Phone:         partner.Phone,
		MailContact:   partner.MailContact,
		Address:       partner.Address,
		WorkSpace:     partner.WorkSpace,
	}
}

// CreateDetail ...
func (factory PartnerInfoFactory) CreateDetail(partner interface{}) models.Partner {
	partnerFactory := models.Partner{}

	err := mapstructure.Decode(partner, &partnerFactory)
	if err != nil {
		panic(err)
	}

	if partnerFactory.Avatar != "" {
		partnerFactory.Avatar = URL_SERVER + partnerFactory.Avatar
	}

	if partnerFactory.Cover != "" {
		partnerFactory.Cover = URL_SERVER + partnerFactory.Cover
	}

	return partnerFactory
}
