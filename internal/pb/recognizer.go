package pb

import (
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	msgCodeAndMsgDescMap = make(map[int16]protoreflect.MessageDescriptor)
	msgNameAndMsgCodeMap = make(map[string]int16)
)

func getMsgDescByMsgCode(code int16) (protoreflect.MessageDescriptor, error) {
	if code < 0 {
		return nil, ErrInvalidMsgCode
	}

	return msgCodeAndMsgDescMap[code], nil
}

func getMsgCodeByMsgName(name string) (int16, error) {
	if name == "" {
		return -1, ErrEmptyMsgName
	}

	return msgNameAndMsgCodeMap[name], nil
}

func InitMaps() {
	for k, v := range MsgCode_value {
		msgName := strings.ToLower(
			strings.ReplaceAll(k, "_", ""),
		)
		msgNameAndMsgCodeMap[msgName] = int16(v)
	}
	
	msgDescLst := File_api_proto_game_msg_proto.Messages()
	for i := 0; i < msgDescLst.Len(); i++ {
		msgDesc := msgDescLst.Get(i)
		msgName := strings.ToLower(
			strings.ReplaceAll(string(msgDesc.Name()), "_", ""),
		)
		msgCodeAndMsgDescMap[msgNameAndMsgCodeMap[msgName]] = msgDesc
	}
}
