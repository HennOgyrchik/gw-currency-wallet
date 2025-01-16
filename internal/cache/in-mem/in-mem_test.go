package in_mem

import (
	"testing"
	"time"
)

func TestInMem_Lifetime(t *testing.T) {
	cache := New(3 * time.Second)

	type result struct {
		key   string
		value any
		ok    bool
	}

	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name    string
		args    args
		timeout time.Duration
		want    result
	}{
		{
			name: "test 1",
			args: args{
				key:   "1",
				value: 1,
			},
			timeout: 0 * time.Second,
			want: result{
				key:   "1",
				value: 1,
				ok:    true,
			},
		},
		{
			name: "test 2",
			args: args{
				key:   "1",
				value: 1,
			},
			timeout: 5 * time.Second,
			want: result{
				key:   "1",
				value: nil,
				ok:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Set(tt.args.key, tt.args.value)

			time.Sleep(tt.timeout)

			v, ok := cache.Get(tt.args.key)
			res := result{
				key:   tt.want.key,
				value: v,
				ok:    ok,
			}

			if tt.want != res {
				t.Errorf("Result was incorrect, got: %v %t, want: %v %t.", v, ok, tt.want.value, tt.want.ok)
			}
		})
	}
}

func TestInMem_Set_Get(t *testing.T) {
	cache := New(10 * time.Second)

	type result struct {
		key   string
		value any
		ok    bool
	}

	type args struct {
		key   string
		value any
	}
	tests := []struct {
		name string
		args args
		want result
	}{
		{
			name: "test 1",
			args: args{
				key:   "1",
				value: 1,
			},
			want: result{
				key:   "1",
				value: 1,
				ok:    true,
			},
		},
		{
			name: "test 2",
			args: args{
				key:   "1",
				value: 1,
			},
			want: result{
				key:   "2",
				value: nil,
				ok:    false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Set(tt.args.key, tt.args.value)

			v, ok := cache.Get(tt.want.key)
			res := result{
				key:   tt.want.key,
				value: v,
				ok:    ok,
			}

			if tt.want != res {
				t.Errorf("Result was incorrect, got: %v %t, want: %v %t.", v, ok, tt.want.value, tt.want.ok)
			}

		})
	}
}
