package serializer

import (
  msgpack "gopkg.in/vmihailenco/msgpack.v2"
)

func EncodeMsgPack(message map[string]interface{}) ([]byte, error) {
    data, err := msgpack.Marshal(message)
    return data, err
}

func DecodeMsgPack(data []byte) (out map[string]interface{}) {
    _ = msgpack.Unmarshal(data, &out)
    return
}

func EncodeArrayMsgPack(message []map[string]interface{}) ([]byte, error) {
    data, err := msgpack.Marshal(message)
    return data, err
}


func DecodeArrayMsgPack(data []byte) (out []map[string]interface{}) {
    _ = msgpack.Unmarshal(data, &out)
    return
}
