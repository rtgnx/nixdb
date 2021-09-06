package nixdb

import (
	"reflect"
	"testing"
)

func TestPasswdEntry_Encode(t *testing.T) {
	type fields struct {
		Name     string
		Password string
		UID      uint
		GID      uint
		Fullname string
		Home     string
		Shell    string
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "Passwd entry with shadowed password",
			fields: fields{
				Name: "user1", Password: "x", UID: 10, GID: 10,
				Fullname: "Username", Home: "/home/user1", Shell: "/bin/sh",
			},
			want: []byte("user1:x:10:10:Username:/home/user1:/bin/sh"),
		},
		{
			name: "Passwd entry with no password",
			fields: fields{
				Name: "user1", Password: "", UID: 10, GID: 10,
				Fullname: "Username", Home: "/home/user1", Shell: "/bin/sh",
			},
			want: []byte("user1::10:10:Username:/home/user1:/bin/sh"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PasswdEntry{
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
				UID:      tt.fields.UID,
				GID:      tt.fields.GID,
				Fullname: tt.fields.Fullname,
				Home:     tt.fields.Home,
				Shell:    tt.fields.Shell,
			}
			if got := p.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PasswdEntry.Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPasswdEntry_Decode(t *testing.T) {
	type fields struct {
		Name     string
		Password string
		UID      uint
		GID      uint
		Fullname string
		Home     string
		Shell    string
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Decode passwd entry with shadowed password",
			args:    args{b: []byte("user1:x:10:10:Username:/home/user1:/bin/sh")},
			wantErr: false,
		},
		{
			name:    "Decode passwd entry with no password",
			args:    args{b: []byte("user1::10:10:Username:/home/user1:/bin/sh")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PasswdEntry{
				Name:     tt.fields.Name,
				Password: tt.fields.Password,
				UID:      tt.fields.UID,
				GID:      tt.fields.GID,
				Fullname: tt.fields.Fullname,
				Home:     tt.fields.Home,
				Shell:    tt.fields.Shell,
			}
			if err := p.Decode(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("PasswdEntry.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
