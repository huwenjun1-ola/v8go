package gots

import "gitee.com/hasika/v8go"

func ConvertJsUint8ArrayToByteSlice(uint8ArrObj *v8go.Object) []byte {
	//arrayLenV, err := uint8ArrObj.Get("length")
	//if err != nil {
	//	panic(err)
	//}
	//arrayLen := arrayLenV.Integer()
	//var byteArray []byte = nil
	//var index int64
	//for index = 0; index < arrayLen; index++ {
	//	v, _ := uint8ArrObj.GetIdx(uint32(index))
	//	b := v.Integer()
	//	byteArray = append(byteArray, byte(b))
	//}
	//return byteArray
	if !uint8ArrObj.IsUint8Array() {
		return nil
	}
	byteArray := uint8ArrObj.GetCopiedArrayBufferViewContents()
	return byteArray
}
