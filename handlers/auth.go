package handlers

import (
	"dreampic/pkg/sb"
	"dreampic/pkg/validate"
	"dreampic/templates/pages/auth"
	"dreampic/templates/ui"
	"fmt"
	"net/http"

	"github.com/angelofallars/htmx-go"
	"github.com/nedpals/supabase-go"
)

func HandleSignIn(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignIn())
}

func HandleSignInGoogleCreate(w http.ResponseWriter, r *http.Request) error {

	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:4000/auth/callback",
	})

	if err != nil {
		return err
	}

	http.Redirect(w, r, resp.URL, http.StatusSeeOther)

	return nil
}

func HandleSignInCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	errors := auth.SignInFormErrors{}

	if ok := validate.New(credentials, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignInForm(credentials, errors))
	}

	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)

	if err != nil {
		return render(r, w, auth.SignInForm(credentials, auth.SignInFormErrors{
			InvalidCredentials: err.Error(),
		}))
	}

	setAuthCookie(w, resp.AccessToken)

	return htmx.NewResponse().Redirect("/").Write(w)

}

func HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.SignUp())
}

func HandleSignUpCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignUpFormProps{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}

	errors := auth.SignUpFormErrors{}

	if ok := validate.New(&params, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password),
			validate.Message("Should be equal"),
		),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignUpForm(params, errors))
	}

	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})

	if err != nil {
		return render(r, w, auth.SignUpForm(params, auth.SignUpFormErrors{
			InvalidCredentials: err.Error(),
		}))
	}

	successMsg := fmt.Sprintf("Congrats your account was created succesfully!! \n We have sent you a confirmation email to: %s", user.Email)

	return render(r, w, ui.SuccessAlert(successMsg))

}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")

	if len(accessToken) == 0 {
		return render(r, w, auth.CallBackScript())
	}

	setAuthCookie(w, accessToken)

	http.Redirect(w, r, "/", http.StatusSeeOther)

	return nil
}

func setAuthCookie(w http.ResponseWriter, accessToken string) {
	cookie := &http.Cookie{
		Value:    accessToken,
		Name:     "at",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, cookie)
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {

	token, err := r.Cookie("at")

	if err != nil {
		return err
	}

	sb.Client.Auth.SignOut(r.Context(), token.Value)

	cookie := http.Cookie{
		Value:    "",
		Name:     "at",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	}

	http.SetCookie(w, &cookie)

	return htmx.NewResponse().Redirect("/signin").Write(w)

}
