package nixdb

import "testing"

func TestHostEntry_Decode(t *testing.T) {
	type fields struct {
		IPAddress string
		FQDN      string
		Aliases   []string
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
			name:    "Test 2 field input",
			args:    args{b: []byte("10.1.0.1\t\tdevice1.domain.com")},
			wantErr: false,
		},
		{
			name:    "Test 3 field input",
			args:    args{b: []byte("10.1.0.1\t\tdevice1.domain.com\tdevice1\t\tdev1")},
			wantErr: false,
		},
		{
			name:    "Test invalid input",
			args:    args{b: []byte("10.1.0.1")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HostEntry{
				IPAddress: tt.fields.IPAddress,
				FQDN:      tt.fields.FQDN,
				Aliases:   tt.fields.Aliases,
			}
			if err := h.Decode(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("HostEntry.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
