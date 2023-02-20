package helper

import (
	"testing"
)

func TestGetLocalIp(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "test GetLocalIp",
			want: "10.8.0.94",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLocalIp(); got != tt.want {
				t.Errorf("GetLocalIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPortInUse(t *testing.T) {
	type args struct {
		host string
		port int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				host: "192.168.22.5",
				port: 99999,
			},
			want: false,
		},
		{
			args: args{
				host: "192.168.0.48",
				port: 41046,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPortInUse(tt.args.host, tt.args.port); got != tt.want {
				t.Errorf("IsPortInUse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPing(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "",
			args: args{
				ip: "192.168.0.48",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				ip: "192.168.22.5",
			},
			want: true,
		},
		{
			name: "",
			args: args{
				ip: "www.baidu.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := Ping(tt.args.ip)
			if got != tt.want {
				t.Errorf("Ping() got = %v, want %v", got, tt.want)
			}
		})
	}
}
