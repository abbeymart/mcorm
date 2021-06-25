// @Author: abbeymart | Abi Akindele | @Created: 2020-12-01 | @Updated: 2020-12-01
// @Company: mConnect.biz | @License: MIT
// @Description: get / query - stream record(s)

package mcorm

import "github.com/abbeymart/mcresponse"

func (crud Crud) GetStream() mcresponse.ResponseMessage {

	return mcresponse.GetResMessage("success", mcresponse.ResponseMessageOptions{
		Message: "success",
		Value:   nil,
	})
}