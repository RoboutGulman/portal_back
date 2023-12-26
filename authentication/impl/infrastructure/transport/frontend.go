package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"portal_back/authentication/api/frontend"
	"portal_back/authentication/impl/app/auth"
	"portal_back/authentication/impl/app/token"
	"time"
)

func NewServer(authService auth.Service, tokenService token.Service) frontendapi.ServerInterface {
	return &frontendServer{authService: authService, tokenService: tokenService}
}

type frontendServer struct {
	authService  auth.Service
	tokenService token.Service
}

func (s *frontendServer) GetSaltByEmail(w http.ResponseWriter, r *http.Request, email string) {
	salt, err := s.authService.GetSaltByEmail(r.Context(), email)
	if err == auth.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(frontendapi.SaltResponse{
		Salt: &salt,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *frontendServer) Login(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var loginReq frontendapi.LoginRequest
	err = json.Unmarshal(reqBody, &loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginData, err := s.authService.Login(
		r.Context(),
		auth.LoginData{
			Email:    *loginReq.Email,
			Password: *loginReq.Password,
		})

	if err == auth.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else if err == nil {

		s.setRefreshTokenToCookie(w, loginData.Tokens.RefreshToken)
		companyID, err := s.authService.GetCompanyByUserID(r.Context(), loginData.User.Id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := json.Marshal(frontendapi.LoginResponse{
			AccessJwtToken: loginData.Tokens.AccessToken,
			User: frontendapi.User{
				Id: loginData.User.Id,
			},
			Company: frontendapi.Company{
				Id: companyID,
			},
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (s *frontendServer) RefreshToken(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokens, err := s.tokenService.RefreshToken(r.Context(), cookie.Value)
	if err == token.ErrUserWithTokenNotFound {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.setRefreshTokenToCookie(w, tokens.RefreshToken)

	resp, err := json.Marshal(frontendapi.TokenResponse{
		AccessJwtToken: &tokens.AccessToken,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *frontendServer) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = s.authService.Logout(r.Context(), cookie.Value)

	if err == auth.ErrUserNotLogged {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *frontendServer) setRefreshTokenToCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(14 * 24 * time.Hour), // 2 weeks
	}
	http.SetCookie(w, &cookie)
}
