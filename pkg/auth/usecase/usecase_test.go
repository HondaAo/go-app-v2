package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/HondaAo/video-app/config"
	"github.com/HondaAo/video-app/log"
	"github.com/HondaAo/video-app/pkg/auth/model"
	authUtils "github.com/HondaAo/video-app/pkg/auth/utils"
)

func Test_authUC_Login(t *testing.T) {
	type fields struct {
		cfg       *config.Config
		authRepo  Repository
		redisRepo RedisRepository
		logger    log.Logger
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secretkey",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
		},
	}
	token, _ := authUtils.GenerateJWTToken(&model.User{
		UserID:   "1",
		Email:    "test@test.com",
		Password: "Test",
	}, cfg)
	apiLogger := log.NewApiLogger(cfg)
	wantUser := NewUserModel()
	wantUser.Password = ""
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserWithToken
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				cfg:       cfg,
				logger:    apiLogger,
				authRepo:  &AuthRepositoryMock{},
				redisRepo: &RedisRepositoryMock{},
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					UserID:   "1",
					Email:    "test@test.com",
					Password: "Test",
				},
			},
			want: &model.UserWithToken{
				User:  wantUser,
				Token: token,
			},
		},
		{
			name: "error",
			fields: fields{
				cfg:       cfg,
				logger:    apiLogger,
				authRepo:  &AuthRepositoryMockError{},
				redisRepo: &RedisRepositoryMock{},
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					UserID:   "1",
					Email:    "test@test.com",
					Password: "Test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &authUC{
				cfg:       tt.fields.cfg,
				authRepo:  tt.fields.authRepo,
				redisRepo: tt.fields.redisRepo,
				logger:    tt.fields.logger,
			}
			got, err := u.Login(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("authUC.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authUC.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authUC_GetByID(t *testing.T) {
	type fields struct {
		cfg       *config.Config
		authRepo  Repository
		redisRepo RedisRepository
		logger    log.Logger
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secret",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
			Level:             "info",
		},
	}
	apiLogger := log.NewApiLogger(cfg)
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				cfg:       cfg,
				logger:    apiLogger,
				authRepo:  &AuthRepositoryMock{},
				redisRepo: &RedisRepositoryMockNil{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
			},
			want: &model.User{
				UserID: "1",
				Email:  "test@test.com",
			},
		},
		{
			name: "ok redis",
			fields: fields{
				cfg:       cfg,
				logger:    apiLogger,
				authRepo:  &AuthRepositoryMock{},
				redisRepo: &RedisRepositoryMock{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
			},
			want: &model.User{
				UserID: "1",
				Email:  "test@test.com",
			},
		},
		{
			name: "get error",
			fields: fields{
				cfg:       cfg,
				logger:    apiLogger,
				authRepo:  &AuthRepositoryMock{},
				redisRepo: &RedisRepositoryMockError{},
			},
			args: args{
				ctx:    context.Background(),
				userID: "1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &authUC{
				cfg:       tt.fields.cfg,
				authRepo:  tt.fields.authRepo,
				redisRepo: tt.fields.redisRepo,
				logger:    tt.fields.logger,
			}
			got, err := u.GetByID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("authUC.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authUC.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authUC_Register(t *testing.T) {
	type fields struct {
		cfg       *config.Config
		authRepo  Repository
		redisRepo RedisRepository
		logger    log.Logger
	}
	type args struct {
		ctx  context.Context
		user *model.User
	}
	cfg := &config.Config{
		Server: config.ServerConfig{
			JwtSecretKey: "secretkey",
		},
		Logger: config.Logger{
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
			Level:             "info",
		},
	}
	token, _ := authUtils.GenerateJWTToken(&model.User{
		UserID:   "1",
		Email:    "test@test.com",
		Password: "Test",
	}, cfg)
	apiLogger := log.NewApiLogger(cfg)
	wantUser := NewUserModel()
	wantUser.Password = ""
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.UserWithToken
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				cfg:       cfg,
				authRepo:  &AuthRepositoryMock{},
				redisRepo: &RedisRepositoryMock{},
				logger:    apiLogger,
			},
			args: args{
				ctx: context.Background(),
				user: &model.User{
					UserID:   "1",
					Email:    "test@test.com",
					Password: "Test",
				},
			},
			want: &model.UserWithToken{
				User:  wantUser,
				Token: token,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &authUC{
				cfg:       tt.fields.cfg,
				authRepo:  tt.fields.authRepo,
				redisRepo: tt.fields.redisRepo,
				logger:    tt.fields.logger,
			}
			got, err := u.Register(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("authUC.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authUC.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}
