package client

import (
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

type Wrapper struct {
	client sdk.ChainClient
}

func CreateCMClientWithConfig(configPath string) (*sdk.ChainClient, error) {
	client, err := sdk.NewChainClient(
		sdk.WithConfPath(configPath),
	)

	if err != nil {
		return nil, err
	}
	// if client.GetAuthType() == sdk.PermissionedWithCert {
	// 	if err := client.EnableCertHash(); err != nil {
	// 		return nil, err
	// 	}
	// }

	return client, nil
}

func CheckProposalRequestResp(resp *common.TxResponse, needContractResult bool) error {
	if resp.Code != common.TxStatusCode_SUCCESS {
		if resp.Message == "" {
			resp.Message = resp.Code.String()
		}
		return errors.New(resp.Message)
	}

	if needContractResult && resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		return errors.New(resp.ContractResult.Message)
	}

	return nil
}
