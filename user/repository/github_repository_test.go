package repository

import (
	"reflect"
	"testing"

	"github.com/usual2970/userhub/domain"
)

func TestGithubRepository_GetAccessToken(t *testing.T) {
	type fields struct {
		clientId     string
		clientSecret string
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.GithubAccessTokenResp
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				clientId:     "a5826b1abce58821f500",
				clientSecret: "c11a0c7264be853e1f1384edb8d3f8589cdc14a3",
			},
			args: args{
				code: "f679772be82ed67c1bc7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GithubRepository{
				clientId:     tt.fields.clientId,
				clientSecret: tt.fields.clientSecret,
			}
			got, err := r.GetAccessToken(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubRepository.GetAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubRepository.GetAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGithubRepository_UserInfo(t *testing.T) {
	type fields struct {
		clientId     string
		clientSecret string
	}
	type args struct {
		accessToken string
		openid      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.GithubUserInfoResp
		wantErr bool
	}{
		{
			name: "1",
			fields: fields{
				clientId:     "a5826b1abce58821f500",
				clientSecret: "c11a0c7264be853e1f1384edb8d3f8589cdc14a3",
			},
			args: args{
				accessToken: "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",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &GithubRepository{
				clientId:     tt.fields.clientId,
				clientSecret: tt.fields.clientSecret,
			}
			got, err := r.UserInfo(tt.args.accessToken, tt.args.openid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GithubRepository.UserInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GithubRepository.UserInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
