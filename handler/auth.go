package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mbaitar/levenue-assignment/db"
	"github.com/mbaitar/levenue-assignment/pkg/kit/validate"
	"github.com/mbaitar/levenue-assignment/pkg/metrics"
	"github.com/mbaitar/levenue-assignment/pkg/sb"
	"github.com/mbaitar/levenue-assignment/pkg/strp"
	"github.com/mbaitar/levenue-assignment/types"
	"github.com/mbaitar/levenue-assignment/view/auth"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/oauth"
	"github.com/uptrace/bun"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/nedpals/supabase-go"
)

const (
	sessionAccessTokenKey = "accessToken"
)

func HandleAccountSetupIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.AccountSetup().Render(r.Context(), w)
}

func HandleAccountSetupCreate(w http.ResponseWriter, r *http.Request) error {
	subscribeStr := r.FormValue("stripeOnboard")
	stripeOnboard, err := strconv.ParseBool(subscribeStr)
	if err != nil {
		stripeOnboard = false
	}

	params := auth.AccountSetupParams{
		StripeConnected: stripeOnboard,
	}

	user := getAuthenticatedUser(r)
	account := types.Account{
		UserID:          user.ID,
		StripeConnected: params.StripeConnected,
	}
	if err := db.CreateAccount(&account); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}

func HandleAccountSetupTypeIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.AccountTypeSetup().Render(r.Context(), w)
}

func HandleAccountSetupTypeCreate(w http.ResponseWriter, r *http.Request) error {
	accountType := r.FormValue("accountType")
	fmt.Println(accountType)
	user := getAuthenticatedUser(r)
	account, err := db.GetAccountByUserID(user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		acc := types.Account{
			UserID: user.ID,
			Type:   accountType,
		}
		if err := db.CreateAccount(&acc); err != nil {
			return err
		}
		if accountType == "Seller" && !account.StripeConnected {
			return hxRedirect(w, r, "/account/setup")
		}
		return hxRedirect(w, r, "/dashboard")
	}

	account.Type = accountType
	if err := db.UpdateAccount(&account); err != nil {
		return err
	}

	if accountType == "Seller" && !account.StripeConnected {
		return hxRedirect(w, r, "/account/setup")
	}

	return hxRedirect(w, r, "/dashboard")
}

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Login())
}

func HandleSignupIndex(w http.ResponseWriter, r *http.Request) error {
	return render(r, w, auth.Signup())
}

func HandleSignupCreate(w http.ResponseWriter, r *http.Request) error {
	params := auth.SignupParams{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirmPassword"),
	}
	errors := auth.SignupErrors{}
	if ok := validate.New(&params, validate.Fields{
		"Email":    validate.Rules(validate.Email),
		"Password": validate.Rules(validate.Password),
		"ConfirmPassword": validate.Rules(
			validate.Equal(params.Password),
			validate.Message("passwords do not match"),
		),
	}).Validate(&errors); !ok {
		return render(r, w, auth.SignupForm(params, errors))
	}
	user, err := sb.Client.Auth.SignUp(r.Context(), supabase.UserCredentials{
		Email:    params.Email,
		Password: params.Password,
	})
	if err != nil {
		return err
	}
	return render(r, w, auth.SignupSuccess(user.Email))
}

func HandleLoginWithGoogle(w http.ResponseWriter, r *http.Request) error {
	resp, err := sb.Client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
		Provider:   "google",
		RedirectTo: "http://localhost:3000/auth/callback",
	})
	if err != nil {
		slog.Error("error on provider", "err", err)
		return err
	}
	http.Redirect(w, r, resp.URL, http.StatusSeeOther)
	return nil
}

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) error {
	accessToken := r.URL.Query().Get("access_token")
	if len(accessToken) == 0 {
		return render(r, w, auth.CallbackScript())
	}
	if err := setAuthSession(w, r, accessToken); err != nil {
		return err
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func HandleLoginCreate(w http.ResponseWriter, r *http.Request) error {
	credentials := supabase.UserCredentials{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	resp, err := sb.Client.Auth.SignIn(r.Context(), credentials)
	if err != nil {
		slog.Error("login error", "err", err)
		return render(r, w, auth.LoginForm(credentials, auth.LoginErrors{
			InvalidCredentials: "The credentials you have entered are invalid",
		}))
	}
	if err := setAuthSession(w, r, resp.AccessToken); err != nil {
		return err
	}
	return hxRedirect(w, r, "/")
}

func HandleStripeAuth(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(types.UserContextKey).(types.AuthenticatedUser)
	url := fmt.Sprintf("https://connect.stripe.com/oauth/authorize?response_type=code&client_id=%s&scope=read_write&redirect_uri=http://localhost:7331/account/stripe/callback&state=%s", os.Getenv("STRIPE_CLIENT_ID"), user.ID)
	return hxRedirect(w, r, url)
}

func HandleStripeConnectCompleted(w http.ResponseWriter, r *http.Request) error {
	user := getAuthenticatedUser(r)
	acc := types.Account{
		UserID:          user.ID,
		Type:            "SELLER",
		StripeConnected: true,
	}

	account, err := db.GetAccountByUserID(user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err := db.CreateAccount(&acc); err != nil {
				return err
			}
		}
		return err
	} else {
		account.StripeConnected = true
		if err := db.UpdateAccount(&account); err != nil {
			return err
		}
	}

	token, err := db.GetTokenByUserID(user.ID)
	if err != nil {
		return err
	}
	if err := strp.FetchAllSubscriptions(token, user.ID); err != nil {
		return err
	}

	if err := metrics.CalculateMetrics(user.ID); err != nil {
		return err
	}

	return hxRedirect(w, r, "/dashboard")
}

func HandleStripeAuthCallback(w http.ResponseWriter, r *http.Request) error {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	userID, err := uuid.Parse(state)
	if err != nil {
		slog.Error("Invalid UUID:", "err", err)
		return err
	}
	stripe.Key = os.Getenv("STRIPE_API_KEY")
	params := &stripe.OAuthTokenParams{
		GrantType: stripe.String("authorization_code"),
		Code:      stripe.String(code),
	}
	token, _ := oauth.New(params)
	err = db.Bun.RunInTx(r.Context(), &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		tkn := types.IntegrationToken{
			AccessToken:        token.AccessToken,
			RefreshToken:       token.RefreshToken,
			AccountID:          userID,
			ConnectedAccountId: token.StripeUserID,
		}
		if err := db.CreateIntegrationTokens(&tkn); err != nil {
			slog.Error("token save failed: ", "err", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	http.Redirect(w, r, "/account/stripe/completed", http.StatusSeeOther)
	return nil
}

func HandleLogoutCreate(w http.ResponseWriter, r *http.Request) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, types.UserContextKey)
	session.Values[sessionAccessTokenKey] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}

func setAuthSession(w http.ResponseWriter, r *http.Request, accessToken string) error {
	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	session, _ := store.Get(r, types.UserContextKey)
	session.Values[sessionAccessTokenKey] = accessToken
	return session.Save(r, w)
}
