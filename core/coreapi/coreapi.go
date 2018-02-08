package coreapi

import (
	core "github.com/ipfs/go-ipfs/core"
	coreiface "github.com/ipfs/go-ipfs/core/coreapi/interface"
	caopts "github.com/ipfs/go-ipfs/core/coreapi/interface/options"
)

type CoreAPI struct {
	node *core.IpfsNode
	*caopts.ApiOptions
}

// NewCoreAPI creates new instance of IPFS CoreAPI backed by go-ipfs Node.
func NewCoreAPI(n *core.IpfsNode) coreiface.CoreAPI {
	api := &CoreAPI{n, nil}
	return api
}

// Unixfs returns the UnixfsAPI interface backed by the go-ipfs node
func (api *CoreAPI) Unixfs() coreiface.UnixfsAPI {
	return (*UnixfsAPI)(api)
}

func (api *CoreAPI) Block() coreiface.BlockAPI {
	return &BlockAPI{api, nil}
}

// Dag returns the DagAPI interface backed by the go-ipfs node
func (api *CoreAPI) Dag() coreiface.DagAPI {
	return &DagAPI{api, nil}
}

// Name returns the NameAPI interface backed by the go-ipfs node
func (api *CoreAPI) Name() coreiface.NameAPI {
	return &NameAPI{api, nil}
}

// Key returns the KeyAPI interface backed by the go-ipfs node
func (api *CoreAPI) Key() coreiface.KeyAPI {
	return &KeyAPI{api, nil}
}

//Object returns the ObjectAPI interface backed by the go-ipfs node
func (api *CoreAPI) Object() coreiface.ObjectAPI {
	return &ObjectAPI{api, nil}
}

func (api *CoreAPI) Pin() coreiface.PinAPI {
	return &PinAPI{api, nil}
}
