package milenage

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type f1TestSet struct {
	K           string
	RAND        string
	SQN         string
	AMF         string
	OP          string
	ExpectedOPc string
	f1          string
	f1Star      string
}

type f2f5f3TestSet struct {
	K           string
	RAND        string
	OP          string
	ExpectedOPc string
	ExpectedRES string
	ExpectedAK  string
	ExpectedCK  string
}

type f4f5StarTestSet struct {
	K              string
	RAND           string
	OP             string
	ExpectedOPc    string
	ExpectedIK     string
	ExpectedAKStar string
}

type ConformanceTestSet struct {
	K              string
	RAND           string
	SQN            string
	AMF            string
	OP             string
	ExpectedOPc    string
	ExpectedMACA   string // f1
	ExpectedMACS   string // f1*
	ExpectedRES    string // f2
	ExpectedCK     string // f3
	ExpectedIK     string // f4
	ExpectedAK     string // f5
	ExpectedAKStar string // f5*
}

// TestF1ImplementorsDataSet according TS35.207 Test Set
func TestF1ImplementorsDataSet(t *testing.T) {
	f1TestCases := []f1TestSet{
		{
			K:           "465b5ce8b199b49faa5f0a2ee238a6bc",
			RAND:        "23553cbe9637a89d218ae64dae47bf35",
			SQN:         "ff9bb4d0b607",
			AMF:         "b9b9",
			OP:          "cdc202d5123e20f62b6d676ac72cb318",
			ExpectedOPc: "cd63cb71954a9f4e48a5994e37a02baf",
			f1:          "4a9ffac354dfafb3",
			f1Star:      "01cfaf9ec4e871e9",
		},
		{
			K:           "0396eb317b6d1c36f19c1c84cd6ffd16",
			RAND:        "c00d603103dcee52c4478119494202e8",
			SQN:         "fd8eef40df7d",
			AMF:         "af17",
			OP:          "ff53bade17df5d4e793073ce9d7579fa",
			ExpectedOPc: "53c15671c60a4b731c55b4a441c0bde2",
			f1:          "5df5b31807e258b0",
			f1Star:      "a8c016e51ef4a343",
		},
		{
			K:           "fec86ba6eb707ed08905757b1bb44b8f",
			RAND:        "9f7c8d021accf4db213ccff0c7f71a6a",
			SQN:         "9d0277595ffc",
			AMF:         "725c",
			OP:          "dbc59adcb6f9a0ef735477b7fadf8374",
			ExpectedOPc: "1006020f0a478bf6b699f15c062e42b3",
			f1:          "9cabc3e99baf7281",
			f1Star:      "95814ba2b3044324",
		},
		{
			K:           "9e5944aea94b81165c82fbf9f32db751",
			RAND:        "ce83dbc54ac0274a157c17f80d017bd6",
			SQN:         "0b604a81eca8",
			AMF:         "9e09",
			OP:          "223014c5806694c007ca1eeef57f004f",
			ExpectedOPc: "a64a507ae1a2a98bb88eb4210135dc87",
			f1:          "74a58220cba84c49",
			f1Star:      "ac2cc74a96871837",
		},
		{
			K:           "4ab1deb05ca6ceb051fc98e77d026a84",
			RAND:        "74b0cd6031a1c8339b2b6ce2b8c4a186",
			SQN:         "e880a1b580b6",
			AMF:         "9f07",
			OP:          "2d16c5cd1fdf6b22383584e3bef2a8d8",
			ExpectedOPc: "dcf07cbd51855290b92a07a9891e523e",
			f1:          "49e785dd12626ef2",
			f1Star:      "9e85790336bb3fa2",
		},
		{
			K:           "6c38a116ac280c454f59332ee35c8c4f",
			RAND:        "ee6466bc96202c5a557abbeff8babf63",
			SQN:         "414b98222181",
			AMF:         "4464",
			OP:          "1ba00a1a7c6700ac8c3ff3e96ad08725",
			ExpectedOPc: "3803ef5363b947c6aaa225e58fae3934",
			f1:          "078adfb488241a57",
			f1Star:      "80246b8d0186bcf1",
		},
	}

	for i, tc := range f1TestCases {
		t.Run(fmt.Sprintf("Set_%d", i+1), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")
			SQN, err := hex.DecodeString(tc.SQN)
			require.NoError(t, err, "decode SQN fail")
			AMF, err := hex.DecodeString(tc.AMF)
			require.NoError(t, err, "decode AMF fail")
			OP, err := hex.DecodeString(tc.OP)
			require.NoError(t, err, "decode OP fail")

			OPC, err := GenerateOPc(K, OP)
			require.NoError(t, err, "calculate OPc fail")
			require.Equal(t, tc.ExpectedOPc, hex.EncodeToString(OPC))
			MAC_A, MAC_S, err := f1(OPC, K, RAND, SQN, AMF)
			require.NoError(t, err, "calculate F1 fail")

			require.Equal(t, tc.f1, hex.EncodeToString(MAC_A))
			require.Equal(t, tc.f1Star, hex.EncodeToString(MAC_S))
		})
	}
}

func TestF2F5F3ImplementorsDataSet(t *testing.T) {
	f2f5f3TestCases := []f2f5f3TestSet{
		{
			K:           "465b5ce8b199b49faa5f0a2ee238a6bc",
			RAND:        "23553cbe9637a89d218ae64dae47bf35",
			OP:          "cdc202d5123e20f62b6d676ac72cb318",
			ExpectedOPc: "cd63cb71954a9f4e48a5994e37a02baf",
			ExpectedRES: "a54211d5e3ba50bf",
			ExpectedAK:  "aa689c648370",
			ExpectedCK:  "b40ba9a3c58b2a05bbf0d987b21bf8cb",
		},
		{
			K:           "0396eb317b6d1c36f19c1c84cd6ffd16",
			RAND:        "c00d603103dcee52c4478119494202e8",
			OP:          "ff53bade17df5d4e793073ce9d7579fa",
			ExpectedOPc: "53c15671c60a4b731c55b4a441c0bde2",
			ExpectedRES: "d3a628ed988620f0",
			ExpectedAK:  "c47783995f72",
			ExpectedCK:  "58c433ff7a7082acd424220f2b67c556",
		},
		{
			K:           "fec86ba6eb707ed08905757b1bb44b8f",
			RAND:        "9f7c8d021accf4db213ccff0c7f71a6a",
			OP:          "dbc59adcb6f9a0ef735477b7fadf8374",
			ExpectedOPc: "1006020f0a478bf6b699f15c062e42b3",
			ExpectedRES: "8011c48c0c214ed2",
			ExpectedAK:  "33484dc2136b",
			ExpectedCK:  "5dbdbb2954e8f3cde665b046179a5098",
		},
		{
			K:           "9e5944aea94b81165c82fbf9f32db751",
			RAND:        "ce83dbc54ac0274a157c17f80d017bd6",
			OP:          "223014c5806694c007ca1eeef57f004f",
			ExpectedOPc: "a64a507ae1a2a98bb88eb4210135dc87",
			ExpectedRES: "f365cd683cd92e96",
			ExpectedAK:  "f0b9c08ad02e",
			ExpectedCK:  "e203edb3971574f5a94b0d61b816345d",
		},
		{
			K:           "4ab1deb05ca6ceb051fc98e77d026a84",
			RAND:        "74b0cd6031a1c8339b2b6ce2b8c4a186",
			OP:          "2d16c5cd1fdf6b22383584e3bef2a8d8",
			ExpectedOPc: "dcf07cbd51855290b92a07a9891e523e",
			ExpectedRES: "5860fc1bce351e7e",
			ExpectedAK:  "31e11a609118",
			ExpectedCK:  "7657766b373d1c2138f307e3de9242f9",
		},
		{
			K:           "6c38a116ac280c454f59332ee35c8c4f",
			RAND:        "ee6466bc96202c5a557abbeff8babf63",
			OP:          "1ba00a1a7c6700ac8c3ff3e96ad08725",
			ExpectedOPc: "3803ef5363b947c6aaa225e58fae3934",
			ExpectedRES: "16c8233f05a0ac28",
			ExpectedAK:  "45b0f69ab06c",
			ExpectedCK:  "3f8c7587fe8e4b233af676aede30ba3b",
		},
	}

	for i, tc := range f2f5f3TestCases {
		t.Run(fmt.Sprintf("Set_%d", i+1), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")
			OP, err := hex.DecodeString(tc.OP)
			require.NoError(t, err, "decode OP fail")

			OPc, err := GenerateOPc(K, OP)
			require.NoError(t, err, "calculate OPc fail")
			require.Equal(t, tc.ExpectedOPc, hex.EncodeToString(OPc))

			RES, CK, _, AK, _, err := f2345(OPc, K, RAND)
			require.NoError(t, err, "calculate F2345 fail")
			require.Equal(t, tc.ExpectedAK, hex.EncodeToString(AK))
			require.Equal(t, tc.ExpectedCK, hex.EncodeToString(CK))
			require.Equal(t, tc.ExpectedRES, hex.EncodeToString(RES))
		})
	}
}

func TestF4F5StarImplementorsDataSet(t *testing.T) {
	f4f5StarTestCases := []f4f5StarTestSet{
		{
			K:              "465b5ce8b199b49faa5f0a2ee238a6bc",
			RAND:           "23553cbe9637a89d218ae64dae47bf35",
			OP:             "cdc202d5123e20f62b6d676ac72cb318",
			ExpectedOPc:    "cd63cb71954a9f4e48a5994e37a02baf",
			ExpectedIK:     "f769bcd751044604127672711c6d3441",
			ExpectedAKStar: "451e8beca43b",
		},
		{
			K:              "0396eb317b6d1c36f19c1c84cd6ffd16",
			RAND:           "c00d603103dcee52c4478119494202e8",
			OP:             "ff53bade17df5d4e793073ce9d7579fa",
			ExpectedOPc:    "53c15671c60a4b731c55b4a441c0bde2",
			ExpectedIK:     "21a8c1f929702adb3e738488b9f5c5da",
			ExpectedAKStar: "30f1197061c1",
		},
		{
			K:              "fec86ba6eb707ed08905757b1bb44b8f",
			RAND:           "9f7c8d021accf4db213ccff0c7f71a6a",
			OP:             "dbc59adcb6f9a0ef735477b7fadf8374",
			ExpectedOPc:    "1006020f0a478bf6b699f15c062e42b3",
			ExpectedIK:     "59a92d3b476a0443487055cf88b2307b",
			ExpectedAKStar: "deacdd848cc6",
		},
		{
			K:              "9e5944aea94b81165c82fbf9f32db751",
			RAND:           "ce83dbc54ac0274a157c17f80d017bd6",
			OP:             "223014c5806694c007ca1eeef57f004f",
			ExpectedOPc:    "a64a507ae1a2a98bb88eb4210135dc87",
			ExpectedIK:     "0c4524adeac041c4dd830d20854fc46b",
			ExpectedAKStar: "6085a86c6f63",
		},
		{
			K:              "4ab1deb05ca6ceb051fc98e77d026a84",
			RAND:           "74b0cd6031a1c8339b2b6ce2b8c4a186",
			OP:             "2d16c5cd1fdf6b22383584e3bef2a8d8",
			ExpectedOPc:    "dcf07cbd51855290b92a07a9891e523e",
			ExpectedIK:     "1c42e960d89b8fa99f2744e0708ccb53",
			ExpectedAKStar: "fe2555e54aa9",
		},
		{
			K:              "6c38a116ac280c454f59332ee35c8c4f",
			RAND:           "ee6466bc96202c5a557abbeff8babf63",
			OP:             "1ba00a1a7c6700ac8c3ff3e96ad08725",
			ExpectedOPc:    "3803ef5363b947c6aaa225e58fae3934",
			ExpectedIK:     "a7466cc1e6b2a1337d49d3b66e95d7b4",
			ExpectedAKStar: "1f53cd2b1113",
		},
	}

	for i, tc := range f4f5StarTestCases {
		t.Run(fmt.Sprintf("Set_%d", i+1), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")
			OP, err := hex.DecodeString(tc.OP)
			require.NoError(t, err, "decode OP fail")

			OPc, err := GenerateOPc(K, OP)
			require.NoError(t, err, "calculate OPc fail")
			require.Equal(t, tc.ExpectedOPc, hex.EncodeToString(OPc))

			_, _, IK, _, AKstar, err := f2345(OPc, K, RAND)
			require.NoError(t, err, "calculate F2345 fail")
			require.Equal(t, tc.ExpectedAKStar, hex.EncodeToString(AKstar))
			require.Equal(t, tc.ExpectedIK, hex.EncodeToString(IK))
		})
	}
}

func TestGenerateOPC(t *testing.T) {
	testCases := []struct {
		inOP        string
		inK         string
		expectedOPc string
	}{
		{ // CHT case
			inOP:        "00112233445566778899aabbccddeeff",
			inK:         "000102030405060708090a0b0c0d0e0f",
			expectedOPc: "69d5c2eb2e2e624750541d3bbc692ba5",
		},
		{ // Landslide Re-synch
			inOP:        "63bfa50ee6523365ff14c1f45f88737d",
			inK:         "00000000000000000000000000000000",
			expectedOPc: "c8ffd2aa7a43c926bf2b2826205b9030",
		},
	}

	for tcIdx, tc := range testCases {
		t.Run(fmt.Sprintf("Case_%d", tcIdx), func(t *testing.T) {
			K, err := hex.DecodeString(tc.inK)
			require.NoError(t, err, "decode inK fail")

			OP, err := hex.DecodeString(tc.inOP)
			require.NoError(t, err, "decode inOP fail")

			OPc, err := GenerateOPc(K, OP)
			require.NoError(t, err, "GenerateOPC fail")

			require.Equal(t, tc.expectedOPc, hex.EncodeToString(OPc))
		})
	}
}

// TestConformanceTestSet from TS 35.208
func TestConformanceTestSet(t *testing.T) {
	conformanceTestCases := []ConformanceTestSet{
		{ // Test Set 1
			K:              "465b5ce8b199b49faa5f0a2ee238a6bc",
			RAND:           "23553cbe9637a89d218ae64dae47bf35",
			SQN:            "ff9bb4d0b607",
			AMF:            "b9b9",
			OP:             "cdc202d5123e20f62b6d676ac72cb318",
			ExpectedOPc:    "cd63cb71954a9f4e48a5994e37a02baf",
			ExpectedMACA:   "4a9ffac354dfafb3",
			ExpectedMACS:   "01cfaf9ec4e871e9",
			ExpectedRES:    "a54211d5e3ba50bf",
			ExpectedCK:     "b40ba9a3c58b2a05bbf0d987b21bf8cb",
			ExpectedIK:     "f769bcd751044604127672711c6d3441",
			ExpectedAK:     "aa689c648370",
			ExpectedAKStar: "451e8beca43b",
		},
		{ // Test Set 3
			K:              "fec86ba6eb707ed08905757b1bb44b8f",
			RAND:           "9f7c8d021accf4db213ccff0c7f71a6a",
			SQN:            "9d0277595ffc",
			AMF:            "725c",
			OP:             "dbc59adcb6f9a0ef735477b7fadf8374",
			ExpectedOPc:    "1006020f0a478bf6b699f15c062e42b3",
			ExpectedMACA:   "9cabc3e99baf7281",
			ExpectedMACS:   "95814ba2b3044324",
			ExpectedRES:    "8011c48c0c214ed2",
			ExpectedCK:     "5dbdbb2954e8f3cde665b046179a5098",
			ExpectedIK:     "59a92d3b476a0443487055cf88b2307b",
			ExpectedAK:     "33484dc2136b",
			ExpectedAKStar: "deacdd848cc6",
		},
		{ // Test Set 4
			K:              "9e5944aea94b81165c82fbf9f32db751",
			RAND:           "ce83dbc54ac0274a157c17f80d017bd6",
			SQN:            "0b604a81eca8",
			AMF:            "9e09",
			OP:             "223014c5806694c007ca1eeef57f004f",
			ExpectedOPc:    "a64a507ae1a2a98bb88eb4210135dc87",
			ExpectedMACA:   "74a58220cba84c49",
			ExpectedMACS:   "ac2cc74a96871837",
			ExpectedRES:    "f365cd683cd92e96",
			ExpectedCK:     "e203edb3971574f5a94b0d61b816345d",
			ExpectedIK:     "0c4524adeac041c4dd830d20854fc46b",
			ExpectedAK:     "f0b9c08ad02e",
			ExpectedAKStar: "6085a86c6f63",
		},
		{ // Test Set 5
			K:              "4ab1deb05ca6ceb051fc98e77d026a84",
			RAND:           "74b0cd6031a1c8339b2b6ce2b8c4a186",
			SQN:            "e880a1b580b6",
			AMF:            "9f07",
			OP:             "2d16c5cd1fdf6b22383584e3bef2a8d8",
			ExpectedOPc:    "dcf07cbd51855290b92a07a9891e523e",
			ExpectedMACA:   "49e785dd12626ef2",
			ExpectedMACS:   "9e85790336bb3fa2",
			ExpectedRES:    "5860fc1bce351e7e",
			ExpectedCK:     "7657766b373d1c2138f307e3de9242f9",
			ExpectedIK:     "1c42e960d89b8fa99f2744e0708ccb53",
			ExpectedAK:     "31e11a609118",
			ExpectedAKStar: "fe2555e54aa9",
		},
		{ // Test Set 6
			K:              "6c38a116ac280c454f59332ee35c8c4f",
			RAND:           "ee6466bc96202c5a557abbeff8babf63",
			SQN:            "414b98222181",
			AMF:            "4464",
			OP:             "1ba00a1a7c6700ac8c3ff3e96ad08725",
			ExpectedOPc:    "3803ef5363b947c6aaa225e58fae3934",
			ExpectedMACA:   "078adfb488241a57",
			ExpectedMACS:   "80246b8d0186bcf1",
			ExpectedRES:    "16c8233f05a0ac28",
			ExpectedCK:     "3f8c7587fe8e4b233af676aede30ba3b",
			ExpectedIK:     "a7466cc1e6b2a1337d49d3b66e95d7b4",
			ExpectedAK:     "45b0f69ab06c",
			ExpectedAKStar: "1f53cd2b1113",
		},
		{ // Test Set 7
			K:              "2d609d4db0ac5bf0d2c0de267014de0d",
			RAND:           "194aa756013896b74b4a2a3b0af4539e",
			SQN:            "6bf69438c2e4",
			AMF:            "5f67",
			OP:             "460a48385427aa39264aac8efc9e73e8",
			ExpectedOPc:    "c35a0ab0bcbfc9252caff15f24efbde0",
			ExpectedMACA:   "bd07d3003b9e5cc3",
			ExpectedMACS:   "bcb6c2fcad152250",
			ExpectedRES:    "8c25a16cd918a1df",
			ExpectedCK:     "4cd0846020f8fa0731dd47cbdc6be411",
			ExpectedIK:     "88ab80a415f15c73711254a1d388f696",
			ExpectedAK:     "7e6455f34cf3",
			ExpectedAKStar: "dc6dd01e8f15",
		},
		{ // Test Set 8
			K:              "a530a7fe428fad1082c45eddfce13884",
			RAND:           "3a4c2b3245c50eb5c71d08639395764d",
			SQN:            "f63f5d768784",
			AMF:            "b90e",
			OP:             "511c6c4e83e38c89b1c5d8dde62426fa",
			ExpectedOPc:    "27953e49bc8af6dcc6e730eb80286be3",
			ExpectedMACA:   "53761fbd679b0bad",
			ExpectedMACS:   "21adfd334a10e7ce",
			ExpectedRES:    "a63241e1ffc3e5ab",
			ExpectedCK:     "10f05bab75a99a5fbb98a9c287679c3b",
			ExpectedIK:     "f9ec0865eb32f22369cade40c59c3a44",
			ExpectedAK:     "88196c47986f",
			ExpectedAKStar: "c987a3d23115",
		},
		{ // Test Set 9
			K:              "d9151cf04896e25830bf2e08267b8360",
			RAND:           "f761e5e93d603feb730e27556cb8a2ca",
			SQN:            "47ee0199820a",
			AMF:            "9113",
			OP:             "75fc2233a44294ee8e6de25c4353d26b",
			ExpectedOPc:    "c4c93effe8a08138c203d4c27ce4e3d9",
			ExpectedMACA:   "66cc4be44862af1f",
			ExpectedMACS:   "7a4b8d7a8753f246",
			ExpectedRES:    "4a90b2171ac83a76",
			ExpectedCK:     "71236b7129f9b22ab77ea7a54c96da22",
			ExpectedIK:     "90527ebaa5588968db41727325a04d9e",
			ExpectedAK:     "82a0f5287a71",
			ExpectedAKStar: "527dbf41f35f",
		},
		{ // Test Set 10
			K:              "a0e2971b6822e8d354a18cc235624ecb",
			RAND:           "08eff828b13fdb562722c65c7f30a9b2",
			SQN:            "db5c066481e0",
			AMF:            "716b",
			OP:             "323792faca21fb4d5d6f13c145a9d2c1",
			ExpectedOPc:    "82a26f22bba9e9488f949a10d98e9cc4",
			ExpectedMACA:   "9485fe24621cb9f6",
			ExpectedMACS:   "bce325ce03e2e9b9",
			ExpectedRES:    "4bc2212d8624910a",
			ExpectedCK:     "08cef6d004ec61471a3c3cda048137fa",
			ExpectedIK:     "ed0318ca5deb9206272f6e8fa64ba411",
			ExpectedAK:     "a2f858aa9e5d",
			ExpectedAKStar: "74e76fbbec38",
		},
		{ // Test Set 11
			K:              "0da6f7ba86d5eac8a19cf563ac58642d",
			RAND:           "679ac4dbacd7d233ff9d6806f4149ce3",
			SQN:            "6e2331d692ad",
			AMF:            "224a",
			OP:             "4b9a26fa459e3acbff36f4015de3bdc1",
			ExpectedOPc:    "0db1071f8767562ca43a0a64c41e8d08",
			ExpectedMACA:   "2831d7ae9088e492",
			ExpectedMACS:   "9b2e16951135d523",
			ExpectedRES:    "6fc30fee6d123523",
			ExpectedCK:     "69b1cae7c7429d975e245cacb05a517c",
			ExpectedIK:     "74f24e8c26df58e1b38d7dcd4f1b7fbd",
			ExpectedAK:     "4c539a26e1fa",
			ExpectedAKStar: "07861e126928",
		},
		{ // Test Set 12
			K:              "77b45843c88e58c10d202684515ed430",
			RAND:           "4c47eb3076dc55fe5106cb2034b8cd78",
			SQN:            "fe1a8731005d",
			AMF:            "ad25",
			OP:             "bf3286c7a51409ce95724d503bfe6e70",
			ExpectedOPc:    "d483afae562409a326b5bb0b20c4d762",
			ExpectedMACA:   "08332d7e9f484570",
			ExpectedMACS:   "ed41b734489d5207",
			ExpectedRES:    "aefa357beac2a87a",
			ExpectedCK:     "908c43f0569cb8f74bc971e706c36c5f",
			ExpectedIK:     "c251df0d888dd9329bcf46655b226e40",
			ExpectedAK:     "30ff25cdadf6",
			ExpectedAKStar: "e84ed0d4677e",
		},
		{ // Test Set 13
			K:              "729b17729270dd87ccdf1bfe29b4e9bb",
			RAND:           "311c4c929744d675b720f3b7e9b1cbd0",
			SQN:            "c85c4cf65916",
			AMF:            "5bb2",
			OP:             "d04c9c35bd2262fa810d2924d036fd13",
			ExpectedOPc:    "228c2f2f06ac3268a9e616ee16db4ba1",
			ExpectedMACA:   "ff794fe2f827ebf8",
			ExpectedMACS:   "24fe4dc61e874b52",
			ExpectedRES:    "98dbbd099b3b408d",
			ExpectedCK:     "44c0f23c5493cfd241e48f197e1d1012",
			ExpectedIK:     "0c9fb81613884c2535dd0eabf3b440d8",
			ExpectedAK:     "5380d158cfe3",
			ExpectedAKStar: "87ac3b559fb6",
		},
		{ // Test Set 14
			K:              "d32dd23e89dc662354ca12eb79dd32fa",
			RAND:           "cf7d0ab1d94306950bf12018fbd46887",
			SQN:            "484107e56a43",
			AMF:            "b5e6",
			OP:             "fe75905b9da47d356236d0314e09c32e",
			ExpectedOPc:    "d22a4b4180a5325708a5ff70d9f67ec7",
			ExpectedMACA:   "cf19d62b6a809866",
			ExpectedMACS:   "5d269537e45e2ce6",
			ExpectedRES:    "af4a411e1139f2c2",
			ExpectedCK:     "5af86b80edb70df5292cc1121cbad50c",
			ExpectedIK:     "7f4d6ae7440e18789a8b75ad3f42f03a",
			ExpectedAK:     "217af49272ad",
			ExpectedAKStar: "900e101c677e",
		},
		{ // Test Set 15
			K:              "af7c65e1927221de591187a2c5987a53",
			RAND:           "1f0f8578464fd59b64bed2d09436b57a",
			SQN:            "3d627b01418d",
			AMF:            "84f6",
			OP:             "0c7acb8d95b7d4a31c5aca6d26345a88",
			ExpectedOPc:    "a4cf5c8155c08a7eff418e5443b98e55",
			ExpectedMACA:   "c37cae7805642032",
			ExpectedMACS:   "68cd09a452d8db7c",
			ExpectedRES:    "7bffa5c2f41fbc05",
			ExpectedCK:     "3f8c3f3ccf7625bf77fc94bcfd22fd26",
			ExpectedIK:     "abcbae8fd46115e9961a55d0da5f2078",
			ExpectedAK:     "837fd7b74419",
			ExpectedAKStar: "56e97a6090b1",
		},
		{ // Test Set 16
			K:              "5bd7ecd3d3127a41d12539bed4e7cf71",
			RAND:           "59b75f14251c75031d0bcbac1c2c04c7",
			SQN:            "a298ae8929dc",
			AMF:            "d056",
			OP:             "f967f76038b920a9cd25e10c08b49924",
			ExpectedOPc:    "76089d3c0ff3efdc6e36721d4fceb747",
			ExpectedMACA:   "c3f25cd94309107e",
			ExpectedMACS:   "b0c8ba343665afcc",
			ExpectedRES:    "7e3f44c7591f6f45",
			ExpectedCK:     "d42b2d615e49a03ac275a5aef97af892",
			ExpectedIK:     "0b3f8d024fe6bfafaa982b8f82e319c2",
			ExpectedAK:     "5be11495525d",
			ExpectedAKStar: "4d6a34a1e4eb",
		},
		{ // Test Set 17
			K:              "6cd1c6ceb1e01e14f1b82316a90b7f3d",
			RAND:           "f69b78f300a0568bce9f0cb93c4be4c9",
			SQN:            "b4fce5feb059",
			AMF:            "e4bb",
			OP:             "078bfca9564659ecd8851e84e6c59b48",
			ExpectedOPc:    "a219dc37f1dc7d66738b5843c799f206",
			ExpectedMACA:   "69a90869c268cb7b",
			ExpectedMACS:   "2e0fdcf9fd1cfa6a",
			ExpectedRES:    "70f6bdb9ad21525f",
			ExpectedCK:     "6edaf99e5bd9f85d5f36d91c1272fb4b",
			ExpectedIK:     "d61c853c280dd9c46f297baec386de17",
			ExpectedAK:     "1c408a858b3e",
			ExpectedAKStar: "aa4ae52daa30",
		},
		{ // Test Set 18
			K:              "b73a90cbcf3afb622dba83c58a8415df",
			RAND:           "b120f1c1a0102a2f507dd543de68281f",
			SQN:            "f1e8a523a36d",
			AMF:            "471b",
			OP:             "b672047e003bb952dca6cb8af0e5b779",
			ExpectedOPc:    "df0c67868fa25f748b7044c6e7c245b8",
			ExpectedMACA:   "ebd70341bcd415b0",
			ExpectedMACS:   "12359f5d82220c14",
			ExpectedRES:    "479dd25c20792d63",
			ExpectedCK:     "66195dbed0313274c5ca7766615fa25e",
			ExpectedIK:     "66bec707eb2afc476d7408a8f2927b36",
			ExpectedAK:     "aefdaa5ddd99",
			ExpectedAKStar: "12ec2b87fbb1",
		},
		{ // Test Set 19
			K:              "5122250214c33e723a5dd523fc145fc0",
			RAND:           "81e92b6c0ee0e12ebceba8d92a99dfa5",
			SQN:            "16f3b3f70fc2",
			AMF:            "c3ab",
			OP:             "c9e8763286b5b9ffbdf56e1297d0887b",
			ExpectedOPc:    "981d464c7c52eb6e5036234984ad0bcf",
			ExpectedMACA:   "2a5c23d15ee351d5",
			ExpectedMACS:   "62dae3853f3af9d2",
			ExpectedRES:    "28d7b0f2a2ec3de5",
			ExpectedCK:     "5349fbe098649f948f5d2e973a81c00f",
			ExpectedIK:     "9744871ad32bf9bbd1dd5ce54e3e2e5a",
			ExpectedAK:     "ada15aeb7bb8",
			ExpectedAKStar: "d461bc15475d",
		},
		{ // Test Set 20
			K:              "90dca4eda45b53cf0f12d7c9c3bc6a89",
			RAND:           "9fddc72092c6ad036b6e464789315b78",
			SQN:            "20f813bd4141",
			AMF:            "61df",
			OP:             "3ffcfe5b7b1111589920d3528e84e655",
			ExpectedOPc:    "cb9cccc4b9258e6dca4760379fb82581",
			ExpectedMACA:   "09db94eab4f8149e",
			ExpectedMACS:   "a29468aa9775b527",
			ExpectedRES:    "a95100e2760952cd",
			ExpectedCK:     "b5f2da03883b69f96bf52e029ed9ac45",
			ExpectedIK:     "b4721368bc16ea67875c5598688bb0ef",
			ExpectedAK:     "83cfd54db913",
			ExpectedAKStar: "4f2039392ddc",
		},
	}
	for i, tc := range conformanceTestCases {
		t.Run(fmt.Sprintf("Set_%d", i+1), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")
			SQN, err := hex.DecodeString(tc.SQN)
			require.NoError(t, err, "decode SQN fail")
			AMF, err := hex.DecodeString(tc.AMF)
			require.NoError(t, err, "decode AMF fail")
			OP, err := hex.DecodeString(tc.OP)
			require.NoError(t, err, "decode OP fail")

			OPC, err := GenerateOPc(K, OP)
			require.NoError(t, err, "calculate OPc fail")
			require.Equal(t, tc.ExpectedOPc, hex.EncodeToString(OPC))
			MAC_A, MAC_S, err := f1(OPC, K, RAND, SQN, AMF)
			require.NoError(t, err, "calculate F1 fail")
			require.Equal(t, tc.ExpectedMACA, hex.EncodeToString(MAC_A))
			require.Equal(t, tc.ExpectedMACS, hex.EncodeToString(MAC_S))

			RES, CK, IK, AK, AKstar, err := f2345(OPC, K, RAND)
			require.NoError(t, err, "calculate F2345 fail")
			require.Equal(t, tc.ExpectedRES, hex.EncodeToString(RES))
			require.Equal(t, tc.ExpectedCK, hex.EncodeToString(CK))
			require.Equal(t, tc.ExpectedIK, hex.EncodeToString(IK))
			require.Equal(t, tc.ExpectedAK, hex.EncodeToString(AK))
			require.Equal(t, tc.ExpectedAKStar, hex.EncodeToString(AKstar))
		})
	}
}

func TestGenerateAKAParameters(t *testing.T) {
	testCases := []struct {
		K            string
		OPc          string
		RAND         string
		AMF          string
		SQN          string
		expectedRES  string
		expectedCK   string
		expectedIK   string
		expectedAUTN string
	}{
		{
			K:            "5122250214c33e723a5dd523fc145fc0",
			OPc:          "981d464c7c52eb6e5036234984ad0bcf",
			RAND:         "81e92b6c0ee0e12ebceba8d92a99dfa5",
			AMF:          "c3ab",
			SQN:          "16f3b3f70fc2",
			expectedRES:  "28d7b0f2a2ec3de5",
			expectedCK:   "5349fbe098649f948f5d2e973a81c00f",
			expectedIK:   "9744871ad32bf9bbd1dd5ce54e3e2e5a",
			expectedAUTN: "bb52e91c747ac3ab2a5c23d15ee351d5",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Set_%d", i), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			OPc, err := hex.DecodeString(tc.OPc)
			require.NoError(t, err, "decode OPc fail")
			AMF, err := hex.DecodeString(tc.AMF)
			require.NoError(t, err, "decode AMF fail")
			SQN, err := hex.DecodeString(tc.SQN)
			require.NoError(t, err, "decode SQN fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")

			IK, CK, RES, AUTN, err := GenerateAKAParameters(OPc, K, RAND, SQN, AMF)
			require.NoError(t, err, "calculate AKA parameters fail")
			require.Equal(t, tc.expectedIK, hex.EncodeToString(IK))
			require.Equal(t, tc.expectedCK, hex.EncodeToString(CK))
			require.Equal(t, tc.expectedRES, hex.EncodeToString(RES))
			require.Equal(t, tc.expectedAUTN, hex.EncodeToString(AUTN))
		})
	}
}

func TestGenerateKeysWithAUTN(t *testing.T) {
	testCases := []struct {
		K           string
		OPc         string
		RAND        string
		AUTN        string
		expectedSQN string
		expectedAK  string
		expectedCK  string
		expectedIK  string

		expectedRES string
	}{
		{
			K:           "5122250214c33e723a5dd523fc145fc0",
			OPc:         "981d464c7c52eb6e5036234984ad0bcf",
			RAND:        "81e92b6c0ee0e12ebceba8d92a99dfa5",
			AUTN:        "bb52e91c747ac3ab2a5c23d15ee351d5",
			expectedRES: "28d7b0f2a2ec3de5",
			expectedAK:  "ada15aeb7bb8",
			expectedCK:  "5349fbe098649f948f5d2e973a81c00f",
			expectedIK:  "9744871ad32bf9bbd1dd5ce54e3e2e5a",
			expectedSQN: "16f3b3f70fc2",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Set_%d", i), func(t *testing.T) {
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			OPc, err := hex.DecodeString(tc.OPc)
			require.NoError(t, err, "decode OPc fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")
			AUTN, err := hex.DecodeString(tc.AUTN)
			require.NoError(t, err, "decode AUTN fail")

			SQN, AK, IK, CK, RES, err := GenerateKeysWithAUTN(OPc, K, RAND, AUTN)
			require.NoError(t, err, "calculate AKA parameters fail")
			require.Equal(t, tc.expectedAK, hex.EncodeToString(AK))
			require.Equal(t, tc.expectedIK, hex.EncodeToString(IK))
			require.Equal(t, tc.expectedCK, hex.EncodeToString(CK))
			require.Equal(t, tc.expectedRES, hex.EncodeToString(RES))
			require.Equal(t, tc.expectedSQN, hex.EncodeToString(SQN))
		})
	}
}

func TestValidateAUTS(t *testing.T) {
	ValidateAUTSTestSetTestCases := []struct {
		K             string
		OPc           string
		AUTS          string
		RAND          string
		ExpectedSQNms string
	}{
		{
			K:             "00000000000000000000000000000000",
			OPc:           "c8ffd2aa7a43c926bf2b2826205b9030",
			AUTS:          "797d7a19ca27f99f4363d3ca24be",
			RAND:          "01000000000000002e6f0eb33b7ffde7",
			ExpectedSQNms: "00052c8c338e",
		},
		{
			K:             "00000000000000000012340000000000",
			OPc:           "c8ffd2aa7a43c926bf2b2826205b9030",
			AUTS:          "797d7a19ca27f99f4363d3ca24be",
			RAND:          "01000000000000002e6f0eb33b7ffde7",
			ExpectedSQNms: "",
		},
	}
	for i, tc := range ValidateAUTSTestSetTestCases {
		t.Run(fmt.Sprintf("Set_%d", i+1), func(t *testing.T) {
			OPc, err := hex.DecodeString(tc.OPc)
			require.NoError(t, err, "decode OPc fail")
			K, err := hex.DecodeString(tc.K)
			require.NoError(t, err, "decode K fail")
			AUTS, err := hex.DecodeString(tc.AUTS)
			require.NoError(t, err, "decode AUTS fail")
			RAND, err := hex.DecodeString(tc.RAND)
			require.NoError(t, err, "decode RAND fail")

			SQNms, err := ValidateAUTS(OPc, K, RAND, AUTS)
			if tc.ExpectedSQNms != "" {
				require.NoError(t, err, "validate AUTS fail")
				require.Equal(t, tc.ExpectedSQNms, hex.EncodeToString(SQNms), "SQNms not eqaul")
			} else { // expected fail
				macFail := &MACFailureError{}
				require.ErrorAs(t, err, &macFail, "not return \"MACFailure\"")
			}
		})
	}
}
