// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import "github.com/ubiq/go-ubiq/v7/common"

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main Ubiq network.
var MainnetBootnodes = []string{
	// Ubiq Go Bootnodes
	"enode://e68e5e6e1a27c1191c09ca3b05fe4e391cfb9648e00c6d085ba4b48931345636bc76117282c2155838d98f63d03994bb88ea9e8b8ecc254da8077398af1c6710@104.156.230.85:30388",
	"enode://f0862b1210672c50f32ec7827159aedd16c8790f64083a5830662e853abb04771ff79d88b2165da8741908aff7ded653e4419f0959f52be607c15b76b318f562@45.76.112.217:30388",
	"enode://3c50be8974756f304ac0195a2a11f9b5ba826354c8617d4b58da21a36102928ddecd96395c7227e9dd1409110ec1414d25b1cfe7f9e4b40732c507d605a7b2b9@45.32.179.15:30388",
	"enode://966f1895b085bf7fdad648afed684b79de9e030a7303c1ebd2acae436e69d754e8d5d35238a08112fd049066c0d310d71ca61e94c16ec0dda4336c065674604c@45.32.117.58:30388",
	"enode://c6de80d7e2a4f8f061b9cc2956aba3edaf5a2b5b37307b6c3036833b76b950ad83085761c6133aec208dc83499425c2447ec0eb56d5ec0f484314076ee7de9e0@5.161.58.43:30388",
	"enode://b902a1538d5bbd6c676676c139e9470fcd942e0d299f5db8bd8ea690af9035f696fe3d88118fece4f74949beb4cf2ba9c3437a002f9a9d08e2b4bfc58fac490f@107.191.104.97:30388",
	"enode://7b40de3623783f6608978a929297fccd9ec5df467eee105b432e8c9b486a5b5ea3bd559854175a1f7663c3bb4f815d93e9ba9db9925c45ac6a53fdb49277742a@149.28.222.52:30388",
	"enode://ee14a9f7200b72a0818efe517779f0bbebdff15b5ae113dd7ef08962d22631f3b95da3509d9b077ec5dbbb10fd287c79958f54bf980493c4f0946cfbaf48c722@140.82.48.169:30388",
	"enode://9fdf5e5dc0b27f582178e0dc28c956459ce2849e455bbab0e17b85f6b598402a9abcd56f3897fb059254d7849e078e774472c80ef67a82fb0992ce3089c9566f@149.28.205.32:30388",
}

// RinkebyBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Rinkeby test network.
var RinkebyBootnodes = []string{}

// GoerliBootnodes are the enode URLs of the P2P bootstrap nodes running on the
// Görli test network.
var GoerliBootnodes = []string{}

var V5Bootnodes = []string{
	// Julian - bootnodes
	"enr:-Ka4QPad278UYSffANmJ-Ovx5tZTbFGGmcbVhBV_afgUO0rSduZCl87BMh0E-ftKWGo5uCD5AKhvs7a2aBKvaQf0yTOGAX6-dZ5ng2V0aMrJhJ7H9VuDHjuZgmlkgnY0gmlwhAWhOiuJc2VjcDI1NmsxoQLG3oDX4qT48GG5zClWq6Ptr1orWzcwe2wwNoM7drlQrYRzbmFwwIN0Y3CCdrSDdWRwgna0",
	"enr:-KC4QPvA5nSaIeD1G1FcV0_4uW6nfNlh9ZWhqYrHPuLqPAj9DN8gJ_YSZUE6LKOdbPxedGXXj5uUmFYg2UAZGlltQKNsg2V0aMrJhJ7H9VuDHjuZgmlkgnY0gmlwhGu_aGGJc2VjcDI1NmsxoQO5AqFTjVu9bGdmdsE56UcPzZQuDSmfXbi9jqaQr5A19oRzbmFwwIN0Y3CCdrSDdWRwgna0",
}

const dnsPrefix = "enrtree://AKA3AM6LPBYEUDMVNU3BSVQJ5AD45Y7YPOHJLEF6W26QOE4VTUDPE@"

// KnownDNSNetwork returns the address of a public DNS-based node list for the given
// genesis hash and protocol. See https://github.com/ethereum/discv4-dns-lists for more
// information.
func KnownDNSNetwork(genesis common.Hash, protocol string) string {
	var net string
	switch genesis {
	case MainnetGenesisHash:
		net = "ubiq"
	default:
		return ""
	}
	return dnsPrefix + protocol + "." + net + ".ethdisco.net"
}
