package repository

import (
	"a21hc3NpZ25tZW50/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SessionRepository interface {
	AddSessions(session model.Session) error
	DeleteSession(token string) error
	UpdateSessions(session model.Session) error
	SessionAvailEmail(email string) (model.Session, error)
	SessionAvailToken(token string) (model.Session, error)
	TokenExpired(session model.Session) bool
}

type sessionsRepo struct {
	db *gorm.DB
}

func NewSessionsRepo(db *gorm.DB) *sessionsRepo {
	return &sessionsRepo{db}
}

func (u *sessionsRepo) AddSessions(session model.Session) error {
	result := u.db.Create(&session)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *sessionsRepo) DeleteSession(token string) error {
	result := u.db.Delete(&model.Session{}, "token = ?", token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (u *sessionsRepo) UpdateSessions(session model.Session) error {
	result := u.db.Model(&model.Session{}).Where("email = ?", session.Email).Update("token", session.Token)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (u *sessionsRepo) SessionAvailEmail(email string) (model.Session, error) {
	var sessionemail model.Session
	result := u.db.Where("email = ?", email).First(&sessionemail)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	return sessionemail, nil
}

func (u *sessionsRepo) SessionAvailToken(token string) (model.Session, error) {
	var sessiontoken model.Session
	result := u.db.Where("token = ?", token).First(&sessiontoken)
	if result.Error != nil {
		return model.Session{}, result.Error
	}
	return sessiontoken, nil
}

func (u *sessionsRepo) TokenValidity(token string) (model.Session, error) {
	session, err := u.SessionAvailToken(token)
	if err != nil {
		return model.Session{}, err
	}

	if u.TokenExpired(session) {
		err := u.DeleteSession(token)
		if err != nil {
			return model.Session{}, err
		}
		return model.Session{}, fmt.Errorf("session expired")
	}

	return session, nil
}

func (u *sessionsRepo) TokenExpired(session model.Session) bool {
	return session.Expiry.Before(time.Now())
}
