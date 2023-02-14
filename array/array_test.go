package array

import (
	"reflect"
	"testing"
)

func TestNewArrayType(t *testing.T) {
	tests := []struct {
		name string
		want *Type
	}{
		{
			name: "test new array type",
			want: NewArrayType(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewArrayType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArrayType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayType_GetValueType(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "check array type test with []string",
			args: args{
				in: []string{},
			},
			want: "[]string",
		},
		{
			name: "check array type test with []int",
			args: args{
				in: []int{},
			},
			want: "[]int",
		},
		{
			name: "check array type test with string",
			args: args{
				in: "",
			},
			want: "string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewArrayType()
			if got := e.GetValueType(tt.args.in); got != tt.want {
				t.Errorf("GetArrayType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayType_Len(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{
			name: "test array type len",
			want: 2,
		},
	}
	e := NewArrayType()
	e.Add(1, "2")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.Len(); got != tt.want {
				t.Errorf("Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayType_Check(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test array type with int 1",
			args: args{
				in: 1,
			},
			want: true,
		},
		{
			name: "test array type with string 1",
			args: args{
				in: "1",
			},
			want: true,
		},
		{
			name: "test array type with string 2",
			args: args{
				in: "2",
			},
			want: false,
		},
	}
	e := NewArrayType()
	e.Add(1, "1")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := e.Check(tt.args.in); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
