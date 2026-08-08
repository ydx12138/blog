package service

import (
	"blog/models"
	"testing"
	"time"
)

func TestUserProfileFromUser(t *testing.T) {
	phone := "13800000000"
	user := models.User{
		ID:           7,
		Email:        "a@example.com",
		Nickname:     "阿明",
		Phone:        &phone,
		Password:     "secret",
		WechatOpenID: "openid",
		CreatedAt:    time.Date(2026, 8, 8, 3, 4, 5, 0, time.UTC),
	}

	profile := UserProfileFromUser(user)
	if profile.ID != 7 || profile.Email != "a@example.com" || profile.Nickname != "阿明" {
		t.Fatalf("unexpected profile identity: %#v", profile)
	}
	if profile.Phone == nil || *profile.Phone != phone {
		t.Fatalf("unexpected profile phone: %#v", profile.Phone)
	}
	if !profile.CreatedAt.Equal(user.CreatedAt) {
		t.Fatalf("unexpected profile creation time: %v", profile.CreatedAt)
	}
}

func TestUserProfileFromUserKeepsNilPhone(t *testing.T) {
	profile := UserProfileFromUser(models.User{})
	if profile.Phone != nil {
		t.Fatal("expected nil phone")
	}
}
