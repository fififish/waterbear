// basicalgorithm
package threshprf

import (
	ecc "crypto/elliptic"
	"crypto/rand"
	sha256 "crypto/sha256"
	"fmt"
	"math/big"
	word "cryptolib/word"
)

func BiginttoU64_256(in *big.Int) [4]uint64 {
	//	out := make([]uint64, 4)
	var out [4]uint64
	in64 := in.Bits()
	if len(in64) < 4 {
		for i := 0; i < len(in64); i++ {
			out[i] = uint64(in64[i])
		}
		for j := len(in64); j < 4; j++ {
			out[j] = uint64(0)
		}
	} else {
		out[0] = uint64(in64[0])
		out[1] = uint64(in64[1])
		out[2] = uint64(in64[2])
		out[3] = uint64(in64[3])
	}

	return out
}

func U64toBigint_256(in [4]uint64) *big.Int {
	out := new(big.Int)
	tmp := new(big.Int)
	out.SetUint64(in[3])
	tmp.SetUint64(in[2])
	out.Lsh(out, 64)
	out.Add(out, tmp)
	tmp.SetUint64(in[1])
	out.Lsh(out, 64)
	out.Add(out, tmp)
	tmp.SetUint64(in[0])
	out.Lsh(out, 64)
	out.Add(out, tmp)
	return out
}

func Bytetostring_256(a [32]byte) []byte {
	var a64 [4]uint64
	a64[3] = uint64(a[0])<<56 | uint64(a[1])<<48 | uint64(a[2])<<40 | uint64(a[3])<<32 | uint64(a[4])<<24 | uint64(a[5])<<16 | uint64(a[6])<<8 | uint64(a[7])
	a64[2] = uint64(a[8])<<56 | uint64(a[9])<<48 | uint64(a[10])<<40 | uint64(a[11])<<32 | uint64(a[12])<<24 | uint64(a[13])<<16 | uint64(a[14])<<8 | uint64(a[15])
	a64[1] = uint64(a[16])<<56 | uint64(a[17])<<48 | uint64(a[18])<<40 | uint64(a[19])<<32 | uint64(a[20])<<24 | uint64(a[21])<<16 | uint64(a[22])<<8 | uint64(a[23])
	a64[0] = uint64(a[24])<<56 | uint64(a[25])<<48 | uint64(a[26])<<40 | uint64(a[27])<<32 | uint64(a[28])<<24 | uint64(a[29])<<16 | uint64(a[30])<<8 | uint64(a[31])
	out := word.U64toByte_256(a64)
	return out
}

/*产生n个参与者的密钥组，每个参与者私钥(SK)长度为n*4 uint64，公钥(VK)长度为n*512 uint64.
  参与者0的公钥为VK[0:7],私钥为SK[0:3]
  参与者1的公钥为VK[8:15],私钥为SK[4:7],以此类推。
//*/
func Gen_key_dealer(n int64, k int64) (VK, SK []uint64) {
	vk64 := make([]uint64, int(n)*8)
	sk64 := make([]uint64, int(n)*4)
	ra64 := make([]uint64, int(k)*4)
	rabyte := make([]byte, 32)
	var tp64 [4]uint64

	p256 := ecc.P256()

	for i := 0; i < int(k); i++ {
		rand.Read(rabyte)
		tp64 = word.BytetoU64_256(rabyte)
		ra64[i*4] = tp64[0]
		ra64[i*4+1] = tp64[1]
		ra64[i*4+2] = tp64[2]
		ra64[i*4+3] = tp64[3]
	}
	//fmt.Printf("\n ra64 is %x, %d", ra64, len(ra64))

	vkbig_x := new(big.Int)
	vkbig_y := new(big.Int)
	skbig := new(big.Int)
	tpbig := new(big.Int)
	ibig := new(big.Int)
	jbig := new(big.Int)

	var i, j int
	for i = 0; i < int(n); i++ {
		tp64[0] = ra64[0]
		tp64[1] = ra64[1]
		tp64[2] = ra64[2]
		tp64[3] = ra64[3]
		skbig = U64toBigint_256(tp64)
		for j = 1; j < int(k); j++ {
			tp64[0] = ra64[4*j]
			tp64[1] = ra64[4*j+1]
			tp64[2] = ra64[4*j+2]
			tp64[3] = ra64[4*j+3]
			tpbig = U64toBigint_256(tp64)
			ibig.SetUint64(uint64(i))
			jbig.SetUint64(uint64(j))
			ibig.Exp(ibig, jbig, p256.Params().N)
			tpbig.Mul(tpbig, ibig)
			tpbig.Mod(tpbig, p256.Params().N)
			skbig.Add(skbig, tpbig)
			skbig.Mod(skbig, p256.Params().N)
		}
		skbig.Mod(skbig, p256.Params().N)
		tp64 = BiginttoU64_256(skbig)
		sk64[4*i] = tp64[0]
		sk64[4*i+1] = tp64[1]
		sk64[4*i+2] = tp64[2]
		sk64[4*i+3] = tp64[3]
		vkbig_x, vkbig_y = p256.ScalarBaseMult(word.U64toByte_256(tp64))
		tp64 = BiginttoU64_256(vkbig_x)
		vk64[8*i] = tp64[0]
		vk64[8*i+1] = tp64[1]
		vk64[8*i+2] = tp64[2]
		vk64[8*i+3] = tp64[3]
		tp64 = BiginttoU64_256(vkbig_y)
		vk64[8*i+4] = tp64[0]
		vk64[8*i+5] = tp64[1]
		vk64[8*i+6] = tp64[2]
		vk64[8*i+7] = tp64[3]
	}
	//	fmt.Printf("\n sk64 is %x, %d", sk64, len(sk64))
	//	fmt.Printf("\n vk64 is %x, %d", vk64, len(vk64))
	return vk64, sk64
}

/*计算拉格朗日插值，idarr为用户id序列索引，k为序列个数，id为目标，lagrange为id插值结果；
lagrange长度为k*4 uint64//*/
func Compute_Lagrangeinter(idarr []int64, k int64, id int64) (lagrang []uint64) {
	if len(idarr) < int(k) {
		return nil
	}
	lag := make([]uint64, 4*k)

	p256 := ecc.P256()

	top := new(big.Int)
	bottom := new(big.Int)

	idbig := new(big.Int)
	idbig = big.NewInt(id)

	sub_top := new(big.Int)
	sub_bott := new(big.Int)
	arrbig := new(big.Int)
	fixidbig := new(big.Int)
	var tp64 [4]uint64

	var j int
	for i := 0; i < int(k); i++ {

		top = big.NewInt(1)
		bottom = big.NewInt(1)

		for j = 0; j < int(k); j++ {
			if i != j {
				arrbig = big.NewInt(idarr[j])
				fixidbig = big.NewInt(idarr[i])

				sub_top.Sub(idbig, arrbig)
				top.Mul(top, sub_top)
				top.Mod(top, p256.Params().N)

				sub_bott.Sub(fixidbig, arrbig)
				bottom.Mul(bottom, sub_bott)
				bottom.Mod(bottom, p256.Params().N)
			}
		}
		bottom.ModInverse(bottom, p256.Params().N)
		top.Mul(top, bottom)
		top.Mod(top, p256.Params().N)
		//fmt.Printf("\n id is %d", idarr[i])
		//fmt.Printf("\n lag is %x", top)
		tp64 = BiginttoU64_256(top)
		//fmt.Printf("\n tp64 is %x", tp64)
		lag[4*i] = tp64[0]
		lag[4*i+1] = tp64[1]
		lag[4*i+2] = tp64[2]
		lag[4*i+3] = tp64[3]
	}
	return lag
}

func Verify_lagrange(idarr []int64, k int64, id_obj int64, lagrang []uint64) {
	if len(idarr) != int(k) || len(lagrang) != 4*int(k) {
		return
	}

	p256 := ecc.P256()
	var id int64
	var vk_x, vk_y []byte
	var tp64 [4]uint64
	var lagb []byte
	vkx := new(big.Int)
	vky := new(big.Int)
	sumx := new(big.Int)
	sumy := new(big.Int)

	for i := 0; i < int(k); i++ {
		id = idarr[i]
		vk_x, vk_y = LoadvkFromFiles(id)
		tp64 = word.BytetoU64_256(vk_x)
		vkx = U64toBigint_256(tp64)
		tp64 = word.BytetoU64_256(vk_y)
		vky = U64toBigint_256(tp64)

		tp64[0] = lagrang[4*i]
		tp64[1] = lagrang[4*i+1]
		tp64[2] = lagrang[4*i+2]
		tp64[3] = lagrang[4*i+3]
		lagb = word.U64toByte_256(tp64)

		vkx, vky = p256.ScalarMult(vkx, vky, lagb)
		if i == 0 {
			sumx.Set(vkx)
			sumy.Set(vky)
		} else {
			sumx, sumy = p256.Add(sumx, sumy, vkx, vky)
		}
	}
	fmt.Printf("\n vkx is %x", sumx)
	fmt.Printf("\n vky is %x", sumy)
}

func Hashmap_point(C []byte) (px [4]uint64, py [4]uint64) {
	if len(C) == 0 {
		return [4]uint64{0, 0, 0, 0}, [4]uint64{0, 0, 0, 0}
	}

	p256 := ecc.P256()


	hashv := sha256.Sum256(C)
	hashs := Bytetostring_256(hashv)

	//fmt.Printf("\n hashv is %x", hashv)

	//fmt.Printf("\n hashs is %x", hashs)

	pxbig, pybig := p256.ScalarBaseMult(hashs)
	px = BiginttoU64_256(pxbig)
	py = BiginttoU64_256(pybig)
	return
}

func Compute_share_own(C []byte) []byte{
	return Compute_share(C, nsk, nvkx, nvky)
}

/*参与者的参数共享，其中share为g||c||z
//*/
func Compute_share(C []byte, sk, vkx, vky []byte) (share []byte) {
	if len(C) == 0 {
		return nil
	}

	p256 := ecc.P256()

	gx, gy := Hashmap_point(C)
	gxb := U64toBigint_256(gx)
	gyb := U64toBigint_256(gy)
	share_x, share_y := p256.ScalarMult(gxb, gyb, sk)
	tp64 := BiginttoU64_256(share_x)
	share = word.U64toByte_256(tp64)
	tp64 = BiginttoU64_256(share_y)
	tpb := word.U64toByte_256(tp64)
	share = append(share, tpb...)
	//fmt.Printf("\n g is %x", share)

	s := make([]byte, 32)
	rand.Read(s)
	hx, hy := p256.ScalarBaseMult(s)
	hx_cap, hy_cap := p256.ScalarMult(gxb, gyb, s)

	tp64 = BiginttoU64_256(p256.Params().Gx)
	h_in := word.U64toByte_256(tp64)
	tp64 = BiginttoU64_256(p256.Params().Gy)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)

	h_in = append(h_in, vkx...)
	h_in = append(h_in, vky...)

	tp64 = BiginttoU64_256(hx)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)
	tp64 = BiginttoU64_256(hy)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)

	tpb = word.U64toByte_256(gx)
	h_in = append(h_in, tpb...)
	tpb = word.U64toByte_256(gy)
	h_in = append(h_in, tpb...)

	h_in = append(h_in, share...)

	tp64 = BiginttoU64_256(hx_cap)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)
	tp64 = BiginttoU64_256(hy_cap)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)

	//fmt.Printf("\n h_in is %x, %d", h_in, len(h_in))
	cv := sha256.Sum256(h_in)
	cin := Bytetostring_256(cv)

	//fmt.Printf("\n cin is %x", cin)

	tp64 = word.BytetoU64_256(s)
	sbig := U64toBigint_256(tp64)
	tp64 = word.BytetoU64_256(sk)
	skbig := U64toBigint_256(tp64)
	tp64 = word.BytetoU64_256(cin)
	cbig := U64toBigint_256(tp64)
	cbig.Mul(cbig, skbig)
	cbig.Add(cbig, sbig)
	cbig.Mod(cbig, p256.Params().N)
	tp64 = BiginttoU64_256(cbig)
	zin := word.U64toByte_256(tp64)
	//fmt.Printf("\n zin is %x", zin)

	share = append(share, cin...)
	share = append(share, zin...)

	//fmt.Printf("\n sharein is %x", share)
	//	fmt.Printf("\n cin is %x", zin)
	return share
}

func Verify_share_node(C []byte, nodeid int64, share []byte) bool{
	vk_x, vk_y := LoadvkFromFiles(nodeid)
	return Verify_share(C,vk_x,vk_y,share)
}

func Verify_share(C, vkx, vky, share []byte) bool {
	if len(C) == 0 || len(vkx) == 0 || len(vky) == 0 || len(share) == 0 {
		return false
	}

	p256 := ecc.P256()


	tp64 := BiginttoU64_256(p256.Params().Gx)
	h_in := word.U64toByte_256(tp64)
	tp64 = BiginttoU64_256(p256.Params().Gy)
	tpb := word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)

	h_in = append(h_in, vkx...)
	h_in = append(h_in, vky...)

	cb := share[64:96]
	zb := share[96:128]
	// fmt.Printf("\n cb is %x", cb)
	// fmt.Printf("\n zb is %x", zb)

	tpx, tpy := p256.ScalarBaseMult(zb)
	tp64 = word.BytetoU64_256(vkx)
	tmx := U64toBigint_256(tp64)
	tp64 = word.BytetoU64_256(vky)
	tmy := U64toBigint_256(tp64)
	tmx, tmy = p256.ScalarMult(tmx, tmy, cb)
	tmy.Sub(p256.Params().P, tmy)
	tmx, tmy = p256.Add(tpx, tpy, tmx, tmy)
	tp64 = BiginttoU64_256(tmx)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)
	tp64 = BiginttoU64_256(tmy)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)

	gx, gy := Hashmap_point(C)
	tpb = word.U64toByte_256(gx)
	h_in = append(h_in, tpb...)
	tpb = word.U64toByte_256(gy)
	h_in = append(h_in, tpb...)

	h_in = append(h_in, share[0:64]...)

	tpx = U64toBigint_256(gx)
	tpy = U64toBigint_256(gy)
	tpx, tpy = p256.ScalarMult(tpx, tpy, zb)
	tp64 = word.BytetoU64_256(share[0:32])
	tmx = U64toBigint_256(tp64)
	tp64 = word.BytetoU64_256(share[32:64])
	tmy = U64toBigint_256(tp64)
	tmx, tmy = p256.ScalarMult(tmx, tmy, cb)
	tmy.Sub(p256.Params().P, tmy)
	tmx, tmy = p256.Add(tpx, tpy, tmx, tmy)
	tp64 = BiginttoU64_256(tmx)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)
	tp64 = BiginttoU64_256(tmy)
	tpb = word.U64toByte_256(tp64)
	h_in = append(h_in, tpb...)
	//fmt.Printf("\n hin is %x, %d", h_in, len(h_in))

	cv := sha256.Sum256(h_in)
	cin := Bytetostring_256(cv)

	// fmt.Printf("\n cb is %x", cb)
	// fmt.Printf("\n cin is %x", cin)

	cb64 := word.BytetoU64_256(cb)
	cin64 := word.BytetoU64_256(cin)

	if cb64 != cin64 {
		return false
	}
	return true
}


func Compute_prf_from_shares(idarr []int64, k int64, shares [][]byte) (prf []byte) {

	lagrang := Compute_Lagrangeinter(idarr, k, int64(0))
	//var id int64
	var share_x, share_y, share_all []byte
	var tp64 [4]uint64
	var lambda []byte

	p256 := ecc.P256()


	gx := new(big.Int)
	gy := new(big.Int)
	tpx := new(big.Int)
	tpy := new(big.Int)

	for i := 0; i < int(k); i++ {
		//id = idarr[i]
		share_all = shares[i]
		
		share_x = share_all[0:32]
		share_y = share_all[32:64]
		tp64 = word.BytetoU64_256(share_x)
		gx = U64toBigint_256(tp64)
		tp64 = word.BytetoU64_256(share_y)
		gy = U64toBigint_256(tp64)

		tp64[0] = lagrang[4*i]
		tp64[1] = lagrang[4*i+1]
		tp64[2] = lagrang[4*i+2]
		tp64[3] = lagrang[4*i+3]
		lambda = word.U64toByte_256(tp64)

		gx, gy = p256.ScalarMult(gx, gy, lambda)
		if i == 0 {
			tpx.Set(gx)
			tpy.Set(gy)
		} else {
			tpx, tpy = p256.Add(gx, gy, tpx, tpy)
		}
	}

	tp64 = BiginttoU64_256(tpx)
	h_in := word.U64toByte_256(tp64)
	tp64 = BiginttoU64_256(tpy)
	h_in = append(h_in, word.U64toByte_256(tp64)...)
	//fmt.Printf("\n h_in is %x", h_in)

	out := sha256.Sum256(h_in)
	outs := Bytetostring_256(out)

	//fmt.Printf("\n outs is %x", outs)
	return outs
}


func Compute_prf(idarr []int64, k int64) (prf []byte) {

	lagrang := Compute_Lagrangeinter(idarr, k, int64(0))
	var id int64
	var share_x, share_y, share_all []byte
	var tp64 [4]uint64
	var lambda []byte

	p256 := ecc.P256()


	gx := new(big.Int)
	gy := new(big.Int)
	tpx := new(big.Int)
	tpy := new(big.Int)

	for i := 0; i < int(k); i++ {
		id = idarr[i]
		share_all = LoadshareFromFiles(id)
		share_x = share_all[0:32]
		share_y = share_all[32:64]
		tp64 = word.BytetoU64_256(share_x)
		gx = U64toBigint_256(tp64)
		tp64 = word.BytetoU64_256(share_y)
		gy = U64toBigint_256(tp64)

		tp64[0] = lagrang[4*i]
		tp64[1] = lagrang[4*i+1]
		tp64[2] = lagrang[4*i+2]
		tp64[3] = lagrang[4*i+3]
		lambda = word.U64toByte_256(tp64)

		gx, gy = p256.ScalarMult(gx, gy, lambda)
		if i == 0 {
			tpx.Set(gx)
			tpy.Set(gy)
		} else {
			tpx, tpy = p256.Add(gx, gy, tpx, tpy)
		}
	}

	tp64 = BiginttoU64_256(tpx)
	h_in := word.U64toByte_256(tp64)
	tp64 = BiginttoU64_256(tpy)
	h_in = append(h_in, word.U64toByte_256(tp64)...)
	//fmt.Printf("\n h_in is %x", h_in)

	out := sha256.Sum256(h_in)
	outs := Bytetostring_256(out)

	//fmt.Printf("\n outs is %x", outs)
	return outs
}

func Init_key_dealer(n, k int64) {
	vk, sk := Gen_key_dealer(int64(n), int64(k))
	Store_key_dealer(vk, sk, int64(n))
}
