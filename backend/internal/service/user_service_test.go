package service

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/city-roam/backend/internal/model"
	"gorm.io/gorm"
)

type fakeProfileUserRepo struct {
	user   *model.User
	err    error
	fields map[string]any
}

func (r *fakeProfileUserRepo) FindByID(context.Context, int64) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.user, nil
}

func (r *fakeProfileUserRepo) FindByUsername(context.Context, string) (*model.User, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *fakeProfileUserRepo) Create(context.Context, *model.User) error {
	return nil
}

func (r *fakeProfileUserRepo) UpdateProfile(_ context.Context, _ int64, fields map[string]any) error {
	r.fields = fields
	return nil
}

func (r *fakeProfileUserRepo) UpdateAccount(_ context.Context, _ int64, fields map[string]any) error {
	r.fields = fields
	return nil
}

func TestUserServiceUpdateProfile(t *testing.T) {
	repo := &fakeProfileUserRepo{user: &model.User{ID: 1, AnonymousID: "anon"}}
	svc := NewUserService(repo)
	age := 28
	nickname := "  北京游客  "
	region := "广东"

	got, err := svc.UpdateProfile(context.Background(), 1, UpdateUserProfileRequest{
		Nickname:   &nickname,
		Age:        &age,
		HomeRegion: &region,
	})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if got.Nickname == nil || *got.Nickname != "北京游客" || got.Age == nil || *got.Age != 28 || got.HomeRegion == nil || *got.HomeRegion != "广东" {
		t.Fatalf("UpdateProfile() = %#v, want normalized profile", got)
	}
	if repo.fields["nickname"] != "北京游客" || repo.fields["age"] != 28 || repo.fields["home_region"] != "广东" {
		t.Fatalf("fields = %#v, want profile updates", repo.fields)
	}
}

func TestUserServiceUpdateProfileValidation(t *testing.T) {
	repo := &fakeProfileUserRepo{user: &model.User{ID: 1, AnonymousID: "anon"}}
	svc := NewUserService(repo)
	age := 130
	_, err := svc.UpdateProfile(context.Background(), 1, UpdateUserProfileRequest{Age: &age})
	if !errors.Is(err, ErrInvalidParam) {
		t.Fatalf("UpdateProfile() error = %v, want ErrInvalidParam", err)
	}
}

func TestUserServiceProfileNotFound(t *testing.T) {
	svc := NewUserService(&fakeProfileUserRepo{err: gorm.ErrRecordNotFound})
	_, err := svc.Profile(context.Background(), 1)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Profile() error = %v, want ErrNotFound", err)
	}
}
