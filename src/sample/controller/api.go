package controller

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/zenazn/goji/web"
)

type StructTestInner struct {
	Id    int
	Score int
	Text  string
}
type StructTestReceive struct {
	Aa          int
	Bb          int
	Float       float32
	IsBool      bool
	Text        string
	NumText     string
	Array       []int
	Struct      StructTestInner
	StructArray []StructTestInner
}

/**************************************************************************************************/
/*!
 *  サンプル実装
 */
/**************************************************************************************************/
func SampleTest(c web.C, w http.ResponseWriter, r *http.Request) {

	var rec = map[string]interface{}{}
	ew := DecryptAndUnpack(c, rec)
	if ew.HasErr() {
		ResError(InCorrectData, "decrypt and unpack error", w, ew)
		return
	}
	d := analyze(rec)
	fmt.Println(d)

	type returnData struct {
		Id          int
		Num         int
		Uid         uint
		Float       float32
		IsBool      bool
		Text        string
		Array       []int
		Struct      StructTestInner
		StructArray []StructTestInner
	}
	ret := new(returnData)

	var structData = StructTestInner{
		Id:    123,
		Score: 9876,
		Text:  "structData",
	}

	var arrays = []StructTestInner{}
	for i := 0; i < 5; i++ {
		d := StructTestInner{}
		d.Id = i
		d.Score = 100 + i
		d.Text = "abcf"
		arrays = append(arrays, d)
	}

	ret.Id = 15
	ret.Num = -98765
	ret.Uid = 4294967294
	ret.Float = 4.567
	ret.IsBool = false
	ret.Text = "てきすと"
	ret.Array = []int{127, 61234, 2147483640, -2147483640}
	ret.Struct = structData
	ret.StructArray = arrays
	fmt.Println("return : ", ret)
	ResWrite(ret, w)
}

// 変数へのアクセステスト
func mapAccessTest(data map[string]interface{}) {
	Aa := data["Aa"].(int)
	Bb := data["Bb"].(int)
	IntMax := data["IntMax"].(int)
	intMin := data["IntMin"].(int)
	Float := data["Float"].(float32)
	FloatInt := data["FloatInt"].(float32)
	IsBool := data["IsBool"].(bool)
	Text := data["Text"].(string)
	NumText := data["NumText"].(string)
	Array := data["Array"].([]interface{})
	Struct := data["Struct"].(map[string]interface{})
	StructArray := data["StructArray"].([]interface{})
	fmt.Println(Aa, Bb, Float, FloatInt, IsBool, Text, NumText, IntMax, intMin)

	testtext := ""
	for _, v := range Array {
		value := v.(int)
		testtext = fmt.Sprintf("%s %d", testtext, value)
	}
	fmt.Println(testtext)

	{
		id := Struct["Id"].(int)
		addScore := Struct["Score"].(int)
		text := Struct["Text"].(string)
		fmt.Println("struct -> ", id, addScore, text)
	}

	{
		testtext = ""
		for i, data := range StructArray {
			m := data.(map[string]interface{})
			id := m["Id"].(int)
			addScore := m["Score"].(int)
			text := m["Text"].(string)
			testtext = fmt.Sprintf("%s %d:[%d, %d, %s]", testtext, i, id, addScore, text)
		}
		fmt.Println(testtext)
	}
}

/**************************************************************************************************/
/*!
 *  受信データ解析
 */
/**************************************************************************************************/
func analyze(rec map[string]interface{}) map[string]interface{} {
	keys := rec["keySlots"].([]interface{})
	values := rec["valueSlots"].([]interface{})
	num := rec["generation"].(int64)
	//log.Debug(keys)
	//log.Debug(values)

	var data = map[string]interface{}{}

	for i := 0; i < int(num); i++ {
		key := keys[i].(string)
		value := values[i]
		// log.Debug(key, ":", value, " : ")
		data[key] = mapping(value, key)
	}
	return data
}

/**************************************************************************************************/
/*!
 *  データマッピング
 */
/**************************************************************************************************/
func mapping(value interface{}, key ...string) interface{} {

	var mv interface{}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := rv.Int()
		mv = int(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := rv.Uint()
		mv = int(v)

	case reflect.Float32, reflect.Float64:
		v := rv.Float()
		mv = float32(v)

	case reflect.String:
		mv = rv.String()

	case reflect.Struct:
		fmt.Println("struct -->", key)

	case reflect.Bool:
		mv = rv.Bool()

	case reflect.Slice, reflect.Array:
		var v []interface{}
		for i := 0; i < rv.Len(); i++ {
			iFace := rv.Index(i).Interface()
			if iFace != nil {
				v = append(v, mapping(iFace))
			}
		}
		mv = v
		return mv

	case reflect.Map:
		mm := value.(map[interface{}]interface{})

		var itemsKey interface{} = "_items"
		var sizeKey interface{} = "_size"

		// 中身が配列で構成されている場合、配列にして返す
		iFace, isArray := mm[itemsKey]
		if isArray {
			array := iFace.([]interface{})
			var v []interface{}
			size := mm[sizeKey].(int64)
			for i := int64(0); i < size; i++ {
				v = append(v, mapping(array[i]))
			}
			mv = v
			break
		}

		// mapを新規作成する
		var newMap = map[string]interface{}{}
		for k, v := range mm {
			s := k.(string)
			newMap[s] = mapping(v)
		}
		mv = newMap
	}
	return mv
}
