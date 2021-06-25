// @Author: abbeymart | Abi Akindele | @Created: 2020-12-24 | @Updated: 2020-12-24
// @Company: mConnect.biz | @License: MIT
// @Description: go: mConnect

package helper

import (
	"fmt"
	"github.com/abbeymart/mcorm/types"
	"github.com/abbeymart/mcresponse"
)

func GetParamsMessage(msgObject types.MessageObject) mcresponse.ResponseMessage {
	var messages = ""

	for key, val := range msgObject {
		if messages != "" {
			messages = fmt.Sprintf("%v | %v : %v", messages, key, val)
		} else {
			messages = fmt.Sprintf("%v : %v", key, val)
		}
	}
	return mcresponse.GetResMessage("validateError", mcresponse.ResponseMessageOptions{
		Message: messages,
		Value:   nil,
	})
}
