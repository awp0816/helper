package helper

import (
	"testing"
)

func TestStringType_Slice2String(t *testing.T) {
	type fields struct {
		inData     interface{}
		linkStr    string
		isUse      bool
		exceptKeys []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test Slice2String with nil []string",
			fields: fields{
				inData:  []string{},
				linkStr: ",",
				isUse:   true,
			},
			want: "",
		},
		{
			name: "test Slice2String with []string",
			fields: fields{
				inData:  []string{"1", "2", "3"},
				linkStr: ",",
				isUse:   true,
			},
			want: "'1','2','3'",
		},
		{
			name: "test Slice2String with []int",
			fields: fields{
				inData:     []int{1, 2, 3},
				linkStr:    "&",
				isUse:      false,
				exceptKeys: []string{"1"},
			},
			want: "2&3",
		},
		{
			name: "test Slice2String with []float64",
			fields: fields{
				inData:     []float64{1.1, 2.1, 3.2},
				linkStr:    "&",
				isUse:      false,
				exceptKeys: []string{"1.1"},
			},
			want: "2.1&3.2",
		},
		{
			name: "test Slice2String with []float64",
			fields: fields{
				inData:     []float64{1.1, 2.1, 3.2},
				linkStr:    "&",
				isUse:      true,
				exceptKeys: []string{"1.1"},
			},
			want: "'2.1'&'3.2'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringType(&SliceType{
				tt.fields.inData,
				tt.fields.linkStr,
				tt.fields.isUse,
				tt.fields.exceptKeys,
			}, nil).Slice2String(); got != tt.want {
				t.Errorf("Slice2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringType_JoinStringsInASCII(t *testing.T) {
	type fields struct {
		Data         map[string]string //需要组装的数据
		Sep          string            //组装数据以哪个字符隔开
		OnlyValues   bool              //是否只要 value
		IncludeEmpty bool              //是否包含空
		ExceptKeys   []string          //无需组装数据的key
		LinkStr      string            //数据以哪个字符串组装
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test JoinStringsInASCII with string map",
			fields: fields{
				Data: map[string]string{
					"lwj": "98",
					"xtn": "97",
				},
				Sep:     "&",
				LinkStr: "=",
			},
			want: "lwj=98&xtn=97",
		},
		{
			name: "test JoinStringsInASCII with string exceptKeys is empty",
			fields: fields{
				Data: map[string]string{
					"lwj": "98",
					"xtn": "97",
					"lmc": "100",
				},
				Sep:     "&",
				LinkStr: "=",
			},
			want: "lmc=100&lwj=98&xtn=97",
		},
		{
			name: "test JoinStringsInASCII with string exceptKeys",
			fields: fields{
				Data: map[string]string{
					"lwj": "98",
					"xtn": "97",
					"lmc": "100",
				},
				Sep:        "&",
				ExceptKeys: []string{"lwj"},
				LinkStr:    "=",
			},
			want: "lmc=100&xtn=97",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStringType(nil, &LinkString{
				Data:         tt.fields.Data,
				Sep:          tt.fields.Sep,
				OnlyValues:   tt.fields.OnlyValues,
				IncludeEmpty: tt.fields.IncludeEmpty,
				ExceptKeys:   tt.fields.ExceptKeys,
				LinkStr:      tt.fields.LinkStr,
			}).JoinStringsInASCII(); got != tt.want {
				t.Errorf("JoinStringsInASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrConvert(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test DataConvert by int",
			args: args{
				value: 1,
			},
			want: "1",
		},
		{
			name: "test DataConvert by string",
			args: args{
				value: "1",
			},
			want: "1",
		},
		{
			name: "test DataConvert by float64",
			args: args{
				value: 3.12,
			},
			want: "3.12",
		},
		{
			name: "test DataConvert by float32",
			args: args{
				value: 3.0,
			},
			want: "3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Interface2String(tt.args.value); got != tt.want {
				t.Errorf("Strval() = %v, want %v", got, tt.want)
			}
		})
	}
}
