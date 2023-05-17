package wrapper

import (
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
