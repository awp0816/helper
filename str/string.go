package str

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type CheckStringTypeI interface {
	Slice2String() string
	JoinStringsInASCII() string
}

type StringType struct {
	sliceType  *SliceType
	linkString *LinkString
	locker     sync.Mutex
}

func NewStringType(sliceType *SliceType, linkString *LinkString) *StringType {
	return &StringType{
		sliceType:  _NewSliceType(sliceType),
		linkString: _NewLinkString(linkString),
	}
}

type SliceType struct {
	inData     interface{} //传入的切片
	linkStr    string      //切片内各元素连接使用的字符串
	isUse      bool        //切片内各元素是否添加单引号
	exceptKeys []string    //需要忽略的数据

}

func _NewSliceType(in *SliceType) *SliceType {
	if in == nil {
		return &SliceType{}
	}
	return &SliceType{
		inData:     in.inData,
		linkStr:    in.linkStr,
		isUse:      in.isUse,
		exceptKeys: in.exceptKeys,
	}
}

func (e *StringType) Slice2String() string {
	defer e.locker.Unlock()
	e.locker.Lock()
	if reflect.ValueOf(e.sliceType.inData).Len() == 0 {
		return ""
	}
	var (
		result      string
		m           = make(map[string]int)
		inArray     []string
		stringRunes []rune
	)
	if len(e.sliceType.exceptKeys) > 0 {
		for _, key := range e.sliceType.exceptKeys {
			m[key] = 1
		}
	}
	if reflect.TypeOf(e.sliceType.inData).Kind() != reflect.Slice {
		return ""
	}
	inData := fmt.Sprintf("%v", e.sliceType.inData)
	for i, v := range inData {
		if i == 0 || i == len(inData)-1 {
			continue
		}
		stringRunes = append(stringRunes, v)
	}
	inArray = strings.Split(string(stringRunes), " ")
	for i, v := range inArray {
		if _, ok := m[v]; ok {
			continue
		}
		if e.sliceType.isUse {
			if i == len(inArray)-1 {
				result += fmt.Sprintf("'%s'", v)
			} else {
				result += fmt.Sprintf("'%s'%s", v, e.sliceType.linkStr)
			}
		} else {
			if i == len(inArray)-1 {
				result += fmt.Sprintf("%s", v)
			} else {
				result += fmt.Sprintf("%s%s", v, e.sliceType.linkStr)
			}
		}
	}
	return result
}

type LinkString struct {
	Data         map[string]string //需要组装的数据
	Sep          string            //组装数据以哪个字符隔开
	OnlyValues   bool              //是否只要 value
	IncludeEmpty bool              //是否包含空
	ExceptKeys   []string          //无需组装数据的key
	LinkStr      string            //数据以哪个字符串组装
}

func _NewLinkString(in *LinkString) *LinkString {
	if in == nil {
		return &LinkString{}
	}
	return &LinkString{
		Data:         in.Data,
		Sep:          in.Sep,
		OnlyValues:   in.OnlyValues,
		IncludeEmpty: in.IncludeEmpty,
		ExceptKeys:   in.ExceptKeys,
		LinkStr:      in.LinkStr,
	}
}

func (e *StringType) JoinStringsInASCII() string {
	defer e.locker.Unlock()
	e.locker.Lock()
	var (
		list, keyList []string
		m             = make(map[string]int)
	)
	if len(e.linkString.ExceptKeys) > 0 {
		for _, except := range e.linkString.ExceptKeys {
			m[except] = 1
		}
	}
	for k := range e.linkString.Data {
		if _, ok := m[k]; ok {
			continue
		}
		value := e.linkString.Data[k]
		if !e.linkString.IncludeEmpty && value == "" {
			continue
		}
		if e.linkString.OnlyValues {
			keyList = append(keyList, k)
		} else {
			list = append(list, fmt.Sprintf("%s%s%s", k, e.linkString.LinkStr, value))
		}
	}
	if e.linkString.OnlyValues {
		sort.Strings(keyList)
		for _, v := range keyList {
			list = append(list, e.linkString.Data[v])
		}
	} else {
		sort.Strings(list)
	}
	return strings.Join(list, e.linkString.Sep)
}

func Interface2String(value interface{}) string {
	if value == nil {
		return ""
	}
	switch value.(type) {
	case float64:
		return strconv.FormatFloat(value.(float64), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 64)
	case int:
		return strconv.Itoa(value.(int))
	case uint:
		return strconv.Itoa(int(value.(uint)))
	case int8:
		return strconv.Itoa(int(value.(int8)))
	case uint8:
		return strconv.Itoa(int(value.(uint8)))
	case int16:
		return strconv.Itoa(int(value.(int16)))
	case uint16:
		return strconv.Itoa(int(value.(uint16)))
	case int32:
		return strconv.Itoa(int(value.(int32)))
	case uint32:
		return strconv.Itoa(int(value.(uint32)))
	case int64:
		return strconv.FormatInt(value.(int64), 10)
	case uint64:
		return strconv.FormatUint(value.(uint64), 10)
	case string:
		return value.(string)
	case []byte:
		return string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		return string(newValue)
	}
}
