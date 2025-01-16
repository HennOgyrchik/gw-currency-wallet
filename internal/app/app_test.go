package app

import "testing"

func Test_getTokenFromString(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test 1",
			args:    args{raw: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.1p-Kr4pbE_ntrvJMSZGYjekyCMTw8CoMR6z30jUuc9M"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.1p-Kr4pbE_ntrvJMSZGYjekyCMTw8CoMR6z30jUuc9M",
			wantErr: false,
		},
		{
			name:    "test 2",
			args:    args{raw: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.1p-Kr4pbE_ntrvJMSZGYjekyCMTw8CoMR6z30jUuc9M"},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.1p-Kr4pbE_ntrvJMSZGYjekyCMTw8CoMR6z30jUuc9M",
			wantErr: false,
		},
		{
			name:    "test 3",
			args:    args{raw: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTokenFromString(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTokenFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTokenFromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adder(t *testing.T) {
	type args struct {
		before     float32
		amount     float32
		multiplier float32
	}
	tests := []struct {
		name    string
		args    args
		want    float32
		wantErr bool
	}{
		{
			name: "пополнить 1:1",
			args: args{
				before:     0,
				amount:     10,
				multiplier: 1,
			},
			want:    10,
			wantErr: false,
		},
		{
			name: "списать больше чем есть 1:1",
			args: args{
				before:     0,
				amount:     10,
				multiplier: -1,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "списать допустимую сумму 1:1",
			args: args{
				before:     10,
				amount:     10,
				multiplier: -1,
			},
			want:    -10,
			wantErr: false,
		},
		{
			name: "пополнить 1:n",
			args: args{
				before:     10,
				amount:     4,
				multiplier: 1.5,
			},
			want:    6,
			wantErr: false,
		},
		{
			name: "списать допустимую сумму 1:n",
			args: args{
				before:     10,
				amount:     10,
				multiplier: -0.5,
			},
			want:    -5,
			wantErr: false,
		},
		{
			name: "списать больше чем есть 1:n",
			args: args{
				before:     10,
				amount:     6,
				multiplier: -2,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := adder(tt.args.before, tt.args.amount, tt.args.multiplier)
			if (err != nil) != tt.wantErr {
				t.Errorf("adder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("adder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
