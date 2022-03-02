package qtum

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestSearchLogsRequestFiltersTopicsIfAllNull(t *testing.T) {
	expected := `[1,2,{"addresses":["0x1","0x2"]},null,1]`
	minConfs := uint(1)
	request := &SearchLogsRequest{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []string{"0x1", "0x2"},
		Topics: []SearchLogsTopic{
			{"0x0", "0x1"},
			{"0x2", "0x3"},
		},
		MinimumConfirmations: &minConfs,
	}

	result, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != expected {
		t.Errorf(
			"error\nwant: %s\ngot: %s",
			expected,
			string(result),
		)
	}
}

func TestSearchLogsRequestGeneratesNulls(t *testing.T) {
	expected := `[1,2,{"addresses":["0x1","0x2"]},{"topics":[null,"0x3"]},1]`
	minConfs := uint(1)
	request := &SearchLogsRequest{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []string{"0x1", "0x2"},
		Topics: []SearchLogsTopic{
			{},
			{"0x3"},
		},
		MinimumConfirmations: &minConfs,
	}

	result, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != expected {
		t.Errorf(
			"error\nwant: %s\ngot: %s",
			expected,
			string(result),
		)
	}
}

func TestSearchLogsRequestFiltersTopicsIfOnlyOneNull(t *testing.T) {
	expected := `[1,2,{"addresses":["0x1","0x2"]},null,1]`
	minConfs := uint(1)
	request := &SearchLogsRequest{
		FromBlock: big.NewInt(1),
		ToBlock:   big.NewInt(2),
		Addresses: []string{"0x1", "0x2"},
		Topics: []SearchLogsTopic{
			{"0x3", "0x4"},
		},
		MinimumConfirmations: &minConfs,
	}

	result, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	if string(result) != expected {
		t.Errorf(
			"error\nwant: %s\ngot: %s",
			expected,
			string(result),
		)
	}
}

// Test extraction of address from script with OP_SENDER
func TestGetOpSenderAddress(t *testing.T) {

	testData := DecodedRawTransactionResponse{
		Vouts: []*DecodedRawTransactionOutV{{
			ScriptPubKey: struct {
				ASM       string   `json:"asm"`
				Hex       string   `json:"hex"`
				ReqSigs   int64    `json:"reqSigs"`
				Type      string   `json:"type"`
				Addresses []string `json:"addresses"`
			}{
				ASM: "1 81e872329e767a0487de7e970992b13b644f1f4f 6b483045022100b83ef90bc808569fb00e29a0f6209d32c1795207c95a554c091401ac8fa8ab920220694b7ec801efd2facea2026d12e8eb5de7689c637f539a620f24c6da8fff235f0121021104b7672c2e08fe321f1bfaffc3768c2777adeedb857b4313ed9d2f15fc8ce4 OP_SENDER 4 55000 40 a9059cbb000000000000000000000000710e94d7f8a5d7a1e5be52bd783370d6e3008a2a0000000000000000000000000000000000000000000000000000000005f5e100 af1ae4e29253ba755c723bca25e883b8deb777b8 OP_CALL",
				Hex: "01011481e872329e767a0487de7e970992b13b644f1f4f4c6c6b483045022100b83ef90bc808569fb00e29a0f6209d32c1795207c95a554c091401ac8fa8ab920220694b7ec801efd2facea2026d12e8eb5de7689c637f539a620f24c6da8fff235f0121021104b7672c2e08fe321f1bfaffc3768c2777adeedb857b4313ed9d2f15fc8ce4c4010403d8d600012844a9059cbb000000000000000000000000710e94d7f8a5d7a1e5be52bd783370d6e3008a2a0000000000000000000000000000000000000000000000000000000005f5e10014af1ae4e29253ba755c723bca25e883b8deb777b8c2",
			},
		}},
	}

	expected := "81e872329e767a0487de7e970992b13b644f1f4f"

	result, err := testData.GetOpSenderAddress()

	if err != nil {
		t.Fatal(err)
	}

	if string(result) != expected {
		t.Errorf(
			"error\n\n-----\n\nwant: %s\n\n-----\n\ngot: %s",
			expected,
			string(result),
		)
	}
}
