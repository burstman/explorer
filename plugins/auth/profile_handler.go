package auth

import (
	"explorer/app/db"
	"explorer/app/handlers"
	"explorer/app/types"
	"fmt"

	"github.com/anthdm/superkit/kit"
	v "github.com/anthdm/superkit/validate"
)

var profileSchema = v.Schema{
	"firstName":            v.Rules(v.Min(3), v.Max(50)),
	"lastName":             v.Rules(v.Min(3), v.Max(50)),
	"phoneNumber":          v.Rules(v.Min(8), v.Max(15)),
	"socialLink":           v.Rules(v.Max(64)),
	"nationalIdentityCard": v.Rules(v.Max(8), DigitOnly),
}

// ProfileFormValues represents the form data for updating a user's profile.
// It contains fields for user identification, name, email, and a success message.
type ProfileFormValues struct {
	ID                   uint   `form:"id"`
	FirstName            string `form:"firstName"`
	LastName             string `form:"lastName"`
	Email                string
	PhoneNumber          string `form:"phoneNumber"`
	SocialLink           string `form:"socialLink"`
	NationalIdentityCard string `form:"nationalIdentityCard"`
	Success              string
}

func HandleProfileShow(kit *kit.Kit) error {
	auth := kit.Auth().(types.AuthUser)

	var user types.User
	if err := db.Get().First(&user, auth.UserID).Error; err != nil {
		return err
	}

	formValues := ProfileFormValues{
		ID:                   user.ID,
		FirstName:            user.FirstName,
		LastName:             user.LastName,
		Email:                user.Email,
		PhoneNumber:          user.PhoneNumber,
		SocialLink:           user.SocialLink,
		NationalIdentityCard: user.Cin,
	}

	return handlers.RenderWithLayout(kit, ProfileShow(formValues))
}

func HandleProfileUpdate(kit *kit.Kit) error {
	var values ProfileFormValues

	errors, ok := v.Request(kit.Request, &values, profileSchema)
	if !ok {
		return kit.Render(ProfileForm(values, errors))
	}

	if values.SocialLink != "" && !isValidSocialLink(values.SocialLink) {
		errors.Add("socialLink", "invalid social link")
	}

	auth := kit.Auth().(types.AuthUser)
	if auth.UserID != values.ID {
		return fmt.Errorf("unauthorized request for profile %d", values.ID)
	}
	//log.Println("values:", values)

	err := db.Get().Model(&types.User{}).
		Where("id = ?", auth.UserID).
		Updates(&types.User{
			FirstName:   values.FirstName,
			LastName:    values.LastName,
			PhoneNumber: values.PhoneNumber,
			SocialLink:  values.SocialLink,
			Cin:         values.NationalIdentityCard,
		}).Error
	if err != nil {
		return err
	}

	values.Success = "Profile successfully updated!"
	values.Email = auth.Email

	return kit.Render(ProfileForm(values, v.Errors{}))
}
