package milenage

import (
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

const (
	K_LEN            = 16
	OP_LEN           = 16
	OPC_LEN          = 16
	SQN_LEN          = 6
	RAND_LEN         = 16
	AMF_LEN          = 2
	MAC_LEN          = 8
	RES_LEN          = 8
	SRES_LEN         = 4
	KC_LEN           = 8
	CK_LEN           = 16
	IK_LEN           = 16
	AK_LEN           = 6
	AUTN_LEN         = 16
	AUTS_LEN         = 14
	CIPHER_BLOCK_LEN = 16
)

// rotate parameter
const (
	r1 = 8
	r3 = 4
	r4 = 8
	r5 = 12
)

// NOTE: The block cipher selected is Rijndael Algorithm with 128-bit block size (AES) in MILENAGE algorithm.

/**
 * f1 - Milenage f1 and f1* algorithms
 * @opc: OPc = 128-bit value derived from OP and K
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @sqn: SQN = 48-bit sequence number
 * @amf: AMF = 16-bit authentication management field
 * Returns:
 * @mac_a: MAC-A = 64-bit network authentication code
 * @mac_s: MAC-S = 64-bit resync authentication code
 */
func f1(opc, k, _rand, sqn, amf []byte) (mac_a, mac_s []byte, err error) {
	mac_a, mac_s = make([]byte, MAC_LEN), make([]byte, MAC_LEN)

	rijndaelInput := make([]uint8, CIPHER_BLOCK_LEN)
	/* tmp1 = TEMP = E_K(RAND XOR OP_C) */
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		rijndaelInput[i] = _rand[i] ^ opc[i]
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, nil, err
	}

	tmp1 := make([]byte, block.BlockSize())
	block.Encrypt(tmp1, rijndaelInput)

	/* tmp2 = IN1 = SQN || AMF || SQN */
	tmp2 := make([]uint8, CIPHER_BLOCK_LEN)
	copy(tmp2[0:], sqn[0:6])
	copy(tmp2[6:], amf[0:2])
	copy(tmp2[8:], tmp2[0:8])

	/* OUT1 = E_K(TEMP XOR rot(IN1 XOR OP_C, r1) XOR c1) XOR OP_C */

	tmp3 := make([]uint8, CIPHER_BLOCK_LEN)
	/* rotate (tmp2 XOR OP_C) by r1 (= 0x40 = 8 bytes) */
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp3[(i+(CIPHER_BLOCK_LEN-r1))%CIPHER_BLOCK_LEN] = tmp2[i] ^ opc[i]
	}

	/* XOR with TEMP = E_K(RAND XOR OP_C) */
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp3[i] ^= tmp1[i]
	}

	/* XOR with c1 (= ..00, i.e., NOP) */
	/* f1 || f1* = E_K(tmp3) XOR OP_c */

	tmp1 = make([]byte, CIPHER_BLOCK_LEN)
	block.Encrypt(tmp1, tmp3)

	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp1[i] ^= opc[i]
	}

	copy(mac_a[0:], tmp1[0:8])
	copy(mac_s[0:], tmp1[8:16])

	return mac_a, mac_s, nil
}

/* AUTS = (SQNms ^ AK-S) || MAC-S */
func CutAUTS(auts []byte) ([]byte, []byte) {
	if len(auts) != AUTS_LEN {
		return nil, nil
	}
	concSQNms := make([]byte, SQN_LEN)
	macS := make([]byte, MAC_LEN)

	copy(concSQNms, auts[0:SQN_LEN])
	copy(macS, auts[SQN_LEN:SQN_LEN+MAC_LEN])

	return concSQNms, macS
}

/* AUTN = (SQNms ^ AK) || AMF || MAC-A */
func CutAUTN(autn []byte) ([]byte, []byte, []byte) {
	if len(autn) != AUTN_LEN {
		return nil, nil, nil
	}
	SQN := make([]byte, SQN_LEN)
	AMF := make([]byte, AMF_LEN)
	MAC := make([]byte, MAC_LEN)

	copy(SQN, autn[0:SQN_LEN])
	copy(AMF, autn[SQN_LEN:SQN_LEN+AMF_LEN])
	copy(MAC, autn[SQN_LEN+AMF_LEN:SQN_LEN+AMF_LEN+MAC_LEN])

	return SQN, AMF, MAC
}

/**
 * f2345 - Milenage f2, f3, f4, f5, f5* algorithms
 * @opc: OPc = 128-bit value derived from OP and K
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * Returns:
 * @res: RES = 64-bit signed response (f2)
 * @ck: CK = 128-bit confidentiality key (f3)
 * @ik: IK = 128-bit integrity key (f4)
 * @ak: AK = 48-bit anonymity key (f5)
 * @akstar: AK = 48-bit anonymity key (f5*)
 */
func f2345(opc, k, _rand []uint8) (res, ck, ik, ak, akstar []byte, err error) {
	res = make([]byte, RES_LEN)
	ck, ik = make([]byte, CK_LEN), make([]byte, IK_LEN)
	ak, akstar = make([]byte, AK_LEN), make([]byte, AK_LEN)

	/* tmp2 = TEMP = E_K(RAND XOR OP_C) */
	tmp1 := make([]uint8, CIPHER_BLOCK_LEN)
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp1[i] = _rand[i] ^ opc[i]
	}

	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	tmp2 := make([]byte, CIPHER_BLOCK_LEN)
	block.Encrypt(tmp2, tmp1)

	/* OUT2 = E_K(rot(TEMP XOR OP_C, r2) XOR c2) XOR OP_C */
	/* OUT3 = E_K(rot(TEMP XOR OP_C, r3) XOR c3) XOR OP_C */
	/* OUT4 = E_K(rot(TEMP XOR OP_C, r4) XOR c4) XOR OP_C */
	/* OUT5 = E_K(rot(TEMP XOR OP_C, r5) XOR c5) XOR OP_C */

	/* f2 and f5 */
	/* rotate by r2 (= 0, i.e., NOP) */
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp1[i] = tmp2[i] ^ opc[i]
	}
	tmp1[15] ^= 1 // XOR c2 (= ..01)

	/* f5 || f2 = E_K(tmp1) XOR OP_c */
	tmp3 := make([]byte, block.BlockSize())
	block.Encrypt(tmp3, tmp1)

	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp3[i] ^= opc[i]
	}

	/* f2 */
	copy(res[0:], tmp3[8:16])

	/* f5 */
	copy(ak[0:], tmp3[0:6])

	/* f3 */
	// rotate by r3 = 0x20 = 4 bytes
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp1[(i+(CIPHER_BLOCK_LEN-r3))%CIPHER_BLOCK_LEN] = tmp2[i] ^ opc[i]
	}
	tmp1[15] ^= 2 // XOR c3 (= ..02)

	block.Encrypt(ck, tmp1)

	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		ck[i] ^= opc[i]
	}

	/* f4 */
	// rotate by r4 = 0x40 = 8 bytes
	for i := 0; i < IK_LEN; i++ {
		tmp1[(i+(CIPHER_BLOCK_LEN-r4))%CIPHER_BLOCK_LEN] = tmp2[i] ^ opc[i]
	}
	tmp1[15] ^= 4 // XOR c4 (= ..04)

	block.Encrypt(ik, tmp1)

	for i := 0; i < IK_LEN; i++ {
		ik[i] ^= opc[i]
	}

	/* f5* */
	// rotate by r5 = 0x60 = 12 bytes
	for i := 0; i < CIPHER_BLOCK_LEN; i++ {
		tmp1[(i+(CIPHER_BLOCK_LEN-r5))%CIPHER_BLOCK_LEN] = tmp2[i] ^ opc[i]
	}
	tmp1[15] ^= 8 // XOR c5 (= ..08)

	block.Encrypt(tmp1, tmp1)

	for i := 0; i < AK_LEN; i++ {
		akstar[i] = tmp1[i] ^ opc[i]
	}

	return res, ck, ik, ak, akstar, nil
}

/**
 * gsm_milenage - Generate GSM-Milenage (3GPP TS 55.205) authentication triplet
 * @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
 * @k: K = 128-bit subscriber key
 * @_rand: RAND = 128-bit random challenge
 * @sres: Buffer for SRES = 32-bit SRES
 * @kc: Buffer for Kc = 64-bit Kc
 * Returns: 0 on success, -1 on failure
 */

/**
* generateAKAParameters - Generate AKA AUTN,IK,CK,RES
* @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
* @k: K = 128-bit subscriber key
* @rand: RAND = 128-bit random challenge
* @sqn: SQN = 48-bit sequence number
* @amf: AMF = 16-bit authentication management field
* Returns:
* @ik: IK = 128-bit integrity key (f4)
* @ck: CK = 128-bit confidentiality key (f3)
* @res: RES = 64-bit signed response (f2)
* @autn: AUTN = 128-bit authentication token
* @err: errors
 */
func generateAKAParameters(opc, k, rand, sqn, amf []byte) (ik, ck, xres, autn []byte, err error) {
	mac, _, err := f1(opc, k, rand, sqn, amf)
	if err != nil {
		err = errors.Wrap(err, "calculate F1 failed")
		return nil, nil, nil, nil, err
	}

	xres, ck, ik, ak, _, err := f2345(opc, k, rand)
	if err != nil {
		err = errors.Wrap(err, "calculate F2345 failed")
		return nil, nil, nil, nil, err
	}

	consSQNhe := xor(sqn, ak)
	autn = append(consSQNhe, append(amf, mac...)...)

	return ik, ck, xres, autn, nil
}

func GenerateAKAParameters(opc, k, rand, sqn, amf []byte) (ik, ck, xres, autn []byte, err error) {
	err = validateArg(opc, "OPc", OPC_LEN)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = validateArg(k, "K", K_LEN)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = validateArg(rand, "RAND", RAND_LEN)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = validateArg(sqn, "SQN", SQN_LEN)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	err = validateArg(amf, "AMF", AMF_LEN)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return generateAKAParameters(opc, k, rand, sqn, amf)
}

/**
* generateKeysWithAUTN - Generate AKA AK,IK,CK,RES from AUTN
* @opc: OPc = 128-bit operator variant algorithm configuration field (encr.)
* @k: K = 128-bit subscriber key
* @rand: RAND = 128-bit random challenge
* @autn: AUTN = 128-bit authentication token
* Returns:
* @sqnhe: SQN = 48-bit sequence number from HM
* @ak: AK = 48-bit anonymity key (f5)
* @ik: IK = 128-bit integrity key (f4)
* @ck: CK = 128-bit confidentiality key (f3)
* @res: RES = 64-bit signed response (f2)
* @err: errors
 */
func generateKeysWithAUTN(opc, k, rand, autn []byte) (sqnhe, ak, ik, ck, res []byte, err error) {
	res, ck, ik, ak, _, err = f2345(opc, k, rand)
	if err != nil {
		err = errors.Wrap(err, "calculate F2345 failed")
		return nil, nil, nil, nil, nil, err
	}

	concSQNhe, amf, xmac := CutAUTN(autn)

	sqnhe = xor(concSQNhe, ak)

	mac, _, err := f1(opc, k, rand, sqnhe, amf)
	if err != nil {
		err = errors.Wrap(err, "calculate F1 failed")
		return nil, nil, nil, nil, nil, err
	}

	if !reflect.DeepEqual(xmac, mac) {
		return nil, nil, nil, nil, nil,
			&MACFailureError{MACName: "MAC-A", ExpectedMAC: xmac, ExactMAC: mac}
	}

	return sqnhe, ak, ik, ck, res, nil
}

func GenerateKeysWithAUTN(opc, k, rand, autn []byte) (sqnhe, ak, ik, ck, res []byte, err error) {
	err = validateArg(opc, "OPc", OPC_LEN)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	err = validateArg(k, "K", K_LEN)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	err = validateArg(rand, "RAND", RAND_LEN)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	err = validateArg(autn, "AUTN", AUTN_LEN)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return generateKeysWithAUTN(opc, k, rand, autn)
}

var _ error = &MACFailureError{}

type MACFailureError struct {
	MACName     string
	ExpectedMAC []byte
	ExactMAC    []byte
}

func (m *MACFailureError) Error() string {
	return fmt.Sprintf("X%s[%x] not match %s[%x]", m.MACName, m.ExpectedMAC, m.MACName, m.ExactMAC)
}

// The AMF used to calculate MAC-S assumes a dummy value of all zeros
var resynchAMF = []byte{0x00, 0x00}

func validateAUTS(opc, k, rand, auts []byte) (sqnms []byte, err error) {
	ConcSQNms, MACS := CutAUTS(auts)
	// nolint:dogsled
	_, _, _, _, AKstar, err := f2345(opc, k, rand)
	if err != nil {
		return nil, errors.Wrap(err, "calculate F2345 Fail")
	}
	SQNms := xor(ConcSQNms, AKstar)

	_, XMACS, err := f1(opc, k, rand, SQNms, resynchAMF)
	if err != nil {
		return nil, errors.Wrap(err, "calculate F1 Fail")
	}

	if !reflect.DeepEqual(XMACS, MACS) {
		return nil, &MACFailureError{MACName: "MAC-S", ExpectedMAC: XMACS, ExactMAC: MACS}
	}

	return SQNms, nil
}

func ValidateAUTS(opc, k, rand, auts []byte) (sqnms []byte, err error) {
	err = validateArg(opc, "OPc", OPC_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(k, "K", K_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(rand, "RAND", RAND_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(auts, "AUTS", AUTS_LEN)
	if err != nil {
		return nil, err
	}

	return validateAUTS(opc, k, rand, auts)
}

func generateAUTS(opc, k, rand, sqnms []byte) (auts []byte, err error) {
	var AKstar []byte
	// nolint:dogsled
	_, _, _, _, AKstar, err = f2345(opc, k, rand)
	if err != nil {
		return nil, errors.Wrap(err, "calculate f2345 Fail")
	}

	ConcSQNms := xor(sqnms, AKstar)

	_, MACS, err := f1(opc, k, rand, sqnms, resynchAMF)
	if err != nil {
		return nil, errors.Wrap(err, "calculate F1 Fail")
	}

	AUTS := append(ConcSQNms, MACS...)

	return AUTS, nil
}

func GenerateAUTS(opc, k, rand, sqnms []byte) (auts []byte, err error) {
	err = validateArg(opc, "OPc", OPC_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(k, "K", K_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(rand, "RAND", RAND_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(sqnms, "SQNms", SQN_LEN)
	if err != nil {
		return nil, err
	}

	return generateAUTS(opc, k, rand, sqnms)
}

// xor return (a xor b) (copy)
func xor(a, b []byte) []byte {
	var outLen int
	if len(a) > len(b) {
		outLen = len(b)
	} else {
		outLen = len(a)
	}

	out := make([]byte, outLen)
	for i := 0; i < outLen; i++ {
		out[i] = a[i] ^ b[i]
	}

	return out
}

type ParameterLengthError struct {
	Name     string
	Exact    int
	Expected int
}

func (e *ParameterLengthError) Error() string {
	return fmt.Sprintf("parameter[%s] length should be %d byte(s), not %d byte(s)",
		e.Name, e.Expected, e.Exact)
}

// validateArg that validate the args is a valid length
func validateArg(arg []byte, argName string, expectedLen int) error {
	if len(arg) != expectedLen {
		return &ParameterLengthError{Name: argName, Exact: len(arg), Expected: expectedLen}
	}
	return nil
}

func validateHexArg(hexArg string, argName string, expectedLen int) ([]byte, error) {
	arg, err := hex.DecodeString(hexArg)
	if err != nil {
		return nil, errors.Wrapf(err, "decode arg[%s]=[%s] fail", argName, hexArg)
	}
	if len(arg) != expectedLen {
		return nil, &ParameterLengthError{Name: argName, Exact: len(arg), Expected: expectedLen}
	}
	return arg, nil
}

func GenerateOPc(k, op []uint8) (opc []uint8, err error) {
	err = validateArg(k, "K", K_LEN)
	if err != nil {
		return nil, err
	}
	err = validateArg(op, "OP", OP_LEN)
	if err != nil {
		return nil, err
	}

	return generateOPc(k, op)
}

func GenerateOPcFromHex(k, op string) (opc string, err error) {
	K, err := validateHexArg(k, "K", K_LEN)
	if err != nil {
		return "", err
	}
	OP, err := validateHexArg(op, "OP", OP_LEN)
	if err != nil {
		return "", err
	}

	OPc, err := generateOPc(K, OP)
	return hex.EncodeToString(OPc), err
}

func generateOPc(k, op []uint8) ([]uint8, error) {
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	opc := make([]byte, block.BlockSize())

	block.Encrypt(opc, op)

	for i := 0; i < OPC_LEN; i++ {
		opc[i] ^= op[i]
	}

	return opc, nil
}
