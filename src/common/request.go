package common

import (
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

const (
	// OffsetDefault offset default
	OffsetDefault = 0
	// OffsetMin offset min
	OffsetMin = 0
	// LimitDefault limit default
	LimitDefault = 10
	// LimitMax limit max
	LimitMax = 100
)

// RequestBody request body
type RequestBody interface {
	// IsLegal 是否合法
	IsLegal() bool
}

// RangeBody range body
type RangeBody struct {
	PageNum  int64
	PageSize int
}

// IsLegal is legal
func (rangeBody *RangeBody) IsLegal() bool {
	if rangeBody.PageSize > LimitMax || rangeBody.PageNum < OffsetMin {
		return false
	}
	return true
}

// ParameterParams parameter params
type ParameterParams struct {
	Key   string
	Value string
}

type InvokeContractListParams struct {
	ContractName string
	ContractAddr string
	MethodName   string
	MethodFunc   string
	Parameters   []*ParameterParams
}

func ConvertToPbKeyValues(keyValues []*ParameterParams) []*pbcommon.KeyValuePair {
	// keyValues := body.Parameters
	if len(keyValues) > 0 {
		pbKvs := make([]*pbcommon.KeyValuePair, 0)
		for _, kv := range keyValues {
			pbKvs = append(pbKvs, &pbcommon.KeyValuePair{
				Key:   kv.Key,
				Value: []byte(kv.Value),
			})
		}
		return pbKvs
	}
	return []*pbcommon.KeyValuePair{}
}
