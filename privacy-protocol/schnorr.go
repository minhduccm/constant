package privacy

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
)

//SchnPubKey denoted Schnorr Publickey
type SchnPubKey struct {
	PK, H EllipticPoint // PK = G^SK + H^R
}

//SchnPrivKey denoted Schnorr Privatekey
type SchnPrivKey struct {
	SK, R  *big.Int
	PubKey *SchnPubKey
}

//SchnSignature denoted Schnorr Signature
type SchnSignature struct {
	E, S1, S2 *big.Int
}

//---------------------------------------------------------------------------------------------------------

//SignScheme contain some algorithm for sign something
type SignScheme interface {
	KeyGen()                //Generate PriKey and PubKey
	GetPubkey() *SchnPubKey //return Publickey belong to the PrivateKey
	Sign(hash []byte) (*SchnSignature, error)
	Verify(signature *SchnSignature, hash []byte) bool
}

//KeyGen Generate PriKey and PubKey
func (priKey *SchnPrivKey) KeyGen() {
	if priKey == nil {
		priKey = new(SchnPrivKey)
	}
	xBytes := RandBytes(32)
	priKey.SK = new(big.Int).SetBytes(xBytes)
	priKey.SK.Mod(priKey.SK, Curve.Params().N)

	rBytes := RandBytes(32)
	priKey.R = new(big.Int).SetBytes(rBytes)
	priKey.R.Mod(priKey.R, Curve.Params().N)

	priKey.PubKey = new(SchnPubKey)

	genPoint := *new(EllipticPoint)
	genPoint.X = Curve.Params().Gx
	genPoint.Y = Curve.Params().Gy

	priKey.PubKey.H = *new(EllipticPoint)
	priKey.PubKey.H.X, priKey.PubKey.H.Y = Curve.ScalarBaseMult(RandBytes(32))
	rH := new(EllipticPoint)
	rH.X, rH.Y = Curve.ScalarMult(priKey.PubKey.H.X, priKey.PubKey.H.Y, priKey.R.Bytes())

	priKey.PubKey.PK = *new(EllipticPoint)
	priKey.PubKey.PK.X, priKey.PubKey.PK.Y = Curve.ScalarBaseMult(priKey.SK.Bytes())
	priKey.PubKey.PK.X, priKey.PubKey.PK.Y = Curve.Add(priKey.PubKey.PK.X, priKey.PubKey.PK.Y, rH.X, rH.Y)

}

//Sign is function which using for sign on hash array by privatekey
func (priKey SchnPrivKey) Sign(hash []byte) (*SchnSignature, error) {
	if len(hash) != 32 {
		return nil, errors.New("Hash length must be 32 bytes")
	}

	genPoint := *new(EllipticPoint)
	genPoint.X = Curve.Params().Gx
	genPoint.Y = Curve.Params().Gy

	signature := new(SchnSignature)

	k1Bytes := RandBytes(32)
	k1 := new(big.Int).SetBytes(k1Bytes)
	k1.Mod(k1, Curve.Params().N)

	k2Bytes := RandBytes(32)
	k2 := new(big.Int).SetBytes(k2Bytes)
	k2.Mod(k2, Curve.Params().N)

	t1 := new(EllipticPoint)
	t1.X, t1.Y = Curve.ScalarMult(Curve.Params().Gx, Curve.Params().Gy, k1.Bytes())

	t2 := new(EllipticPoint)
	t2.X, t2.Y = Curve.ScalarMult(priKey.PubKey.H.X, priKey.PubKey.H.Y, k2.Bytes())

	t := new(EllipticPoint)
	t.X, t.Y = Curve.Add(t1.X, t1.Y, t2.X, t2.Y)

	signature.E = Hash(*t, hash)

	xe := new(big.Int)
	xe.Mul(priKey.SK, signature.E)
	signature.S1 = new(big.Int)
	signature.S1.Sub(k1, xe)
	signature.S1.Mod(signature.S1, Curve.Params().N)

	re := new(big.Int)
	re.Mul(priKey.R, signature.E)
	signature.S2 = new(big.Int)
	signature.S2.Sub(k2, re)
	signature.S2.Mod(signature.S2, Curve.Params().N)

	return signature, nil
}

//Verify is function which using for verify that the given signature was signed by by privatekey of the public key
func (pub SchnPubKey) Verify(signature *SchnSignature, hash []byte) bool {
	if len(hash) != 32 {
		return false
	}

	if signature == nil {
		return false
	}

	rv := new(EllipticPoint)
	rv.X, rv.Y = Curve.ScalarMult(Curve.Params().Gx, Curve.Params().Gy, signature.S1.Bytes())
	tmp := new(EllipticPoint)
	tmp.X, tmp.Y = Curve.ScalarMult(pub.H.X, pub.H.Y, signature.S2.Bytes())
	rv.X, rv.Y = Curve.Add(rv.X, rv.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(pub.PK.X, pub.PK.Y, signature.E.Bytes())
	rv.X, rv.Y = Curve.Add(rv.X, rv.Y, tmp.X, tmp.Y)

	ev := Hash(*rv, hash)
	if ev.Cmp(signature.E) == 0 {
		return true
	}

	return false
}

//---------------------------------------------------------------------------------------------------------

// SchnGenPrivKey generates Schnorr private key
func SchnGenPrivKey() *SchnPrivKey {
	priv := new(SchnPrivKey)
	xBytes := RandBytes(32)
	priv.SK = new(big.Int).SetBytes(xBytes)
	priv.SK.Mod(priv.SK, Curve.Params().N)

	rBytes := RandBytes(32)
	priv.R = new(big.Int).SetBytes(rBytes)
	priv.R.Mod(priv.R, Curve.Params().N)
	priv.PubKey = SchnGenPubKey(*priv)

	return priv
}

func SchnGenPubKey(priv SchnPrivKey) *SchnPubKey {
	pub := new(SchnPubKey)

	genPoint := *new(EllipticPoint)
	genPoint.X = Curve.Params().Gx
	genPoint.Y = Curve.Params().Gy

	pub.H = *new(EllipticPoint)
	pub.H.X, pub.H.Y = Curve.ScalarBaseMult(RandBytes(32))
	rH := new(EllipticPoint)
	rH.X, rH.Y = Curve.ScalarMult(pub.H.X, pub.H.Y, priv.R.Bytes())

	pub.PK = *new(EllipticPoint)
	pub.PK.X, pub.PK.Y = Curve.ScalarBaseMult(priv.SK.Bytes())
	pub.PK.X, pub.PK.Y = Curve.Add(pub.PK.X, pub.PK.Y, rH.X, rH.Y)

	return pub
}

func SchnSign(hash []byte, priv SchnPrivKey) (*SchnSignature, error) {
	if len(hash) != 32 {
		return nil, errors.New("Hash length must be 32 bytes")
	}

	genPoint := *new(EllipticPoint)
	genPoint.X = Curve.Params().Gx
	genPoint.Y = Curve.Params().Gy

	signature := new(SchnSignature)

	k1Bytes := RandBytes(32)
	k1 := new(big.Int).SetBytes(k1Bytes)
	k1.Mod(k1, Curve.Params().N)

	k2Bytes := RandBytes(32)
	k2 := new(big.Int).SetBytes(k2Bytes)
	k2.Mod(k2, Curve.Params().N)

	t1 := new(EllipticPoint)
	t1.X, t1.Y = Curve.ScalarMult(Curve.Params().Gx, Curve.Params().Gy, k1.Bytes())

	t2 := new(EllipticPoint)
	t2.X, t2.Y = Curve.ScalarMult(priv.PubKey.H.X, priv.PubKey.H.Y, k2.Bytes())

	t := new(EllipticPoint)
	t.X, t.Y = Curve.Add(t1.X, t1.Y, t2.X, t2.Y)

	signature.E = Hash(*t, hash)

	xe := new(big.Int)
	xe.Mul(priv.SK, signature.E)
	signature.S1 = new(big.Int)
	signature.S1.Sub(k1, xe)
	signature.S1.Mod(signature.S1, Curve.Params().N)

	re := new(big.Int)
	re.Mul(priv.R, signature.E)
	signature.S2 = new(big.Int)
	signature.S2.Sub(k2, re)
	signature.S2.Mod(signature.S2, Curve.Params().N)

	return signature, nil
}

func SchnVerify(signature *SchnSignature, hash []byte, pub SchnPubKey) bool {
	if len(hash) != 32 {
		return false
	}

	if signature == nil {
		return false
	}

	rv := new(EllipticPoint)
	rv.X, rv.Y = Curve.ScalarMult(Curve.Params().Gx, Curve.Params().Gy, signature.S1.Bytes())
	tmp := new(EllipticPoint)
	tmp.X, tmp.Y = Curve.ScalarMult(pub.H.X, pub.H.Y, signature.S2.Bytes())
	rv.X, rv.Y = Curve.Add(rv.X, rv.Y, tmp.X, tmp.Y)
	tmp.X, tmp.Y = Curve.ScalarMult(pub.PK.X, pub.PK.Y, signature.E.Bytes())
	rv.X, rv.Y = Curve.Add(rv.X, rv.Y, tmp.X, tmp.Y)

	ev := Hash(*rv, hash)
	if ev.Cmp(signature.E) == 0 {
		return true
	}

	return false
}

// Hash calculates a hash concatenating a given message bytes with a given EC Point. H(p||m)
func Hash(p EllipticPoint, m []byte) *big.Int {
	var b []byte
	cXBytes := p.X.Bytes()
	cYBytes := p.Y.Bytes()
	b = append(b, cXBytes...)
	b = append(b, cYBytes...)
	b = append(b, m...)
	h := sha256.New()
	h.Write(b)
	hash := h.Sum(nil)
	r := new(big.Int).SetBytes(hash)
	return r
}

func TestSchn() {
	priv := SchnGenPrivKey()

	hash := RandBytes(32)
	fmt.Printf("Hash: %v\n", hash)

	signature, _ := SchnSign(hash, *priv)
	fmt.Printf("Signature: %+v\n", signature)

	res := SchnVerify(signature, hash, *priv.PubKey)
	fmt.Println(res)

}
