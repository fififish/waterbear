// test
package threshprf

import (
	"fmt"
)

//5/3个用户的测试
func Test_course() {
	// dealer 生成并存储密钥，公共参数C
	n := 5
	k := 3
	Init_key_dealer(int64(n), int64(k))
	C := []byte("0123456789abcdef")

	//==========   与参与者1的交互  =======================
	//参与者1发布共享
	id_1 := 1
	skb, vkx, vky := LoadkeyFromFiles(int64(id_1))
	share := Compute_share(C, skb, vkx, vky)

	//user检测参与者1共享是否正确，如果正确，则存储，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_1))
	} else {
		fmt.Printf("\n share for party 1 is false")
	}
	//====================================================

	//==========   与参与者2的交互  =======================
	//参与者2发布共享
	id_2 := 2
	skb, vkx, vky = LoadkeyFromFiles(int64(id_2))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者2共享是否正确，如果正确，则存储，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_2))
	} else {
		fmt.Printf("\n share for party 2 is false")
	}
	//====================================================

	//==========   与参与者3的交互  =======================
	//参与者3发布共享
	id_3 := 3
	skb, vkx, vky = LoadkeyFromFiles(int64(id_3))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_3))
	} else {
		fmt.Printf("\n share for party 3 is false")
	}
	//====================================================

	//=============  用户生成随机数 ========================
	idarr := []int64{int64(id_3), int64(id_2), int64(id_1)}
	prf := Compute_prf(idarr, int64(k))
	fmt.Printf("\n prf is %x", prf)
	//===================================================

	//============== 验证正确型，此步骤仅为验证算法正确型 ======
	id_0 := 0
	skb, vkx, vky = LoadkeyFromFiles(int64(id_0))
	share = Compute_share(C, skb, vkx, vky)
	fmt.Printf("\n share_0 is %x", share[0:64])
	//====================================================

}

//10/6个用户的测试
func Test_course_10() {
	// dealer 生成并存储密钥，公共参数C
	n := 10
	k := 6
	Init_key_dealer(int64(n), int64(k))
	C := []byte("0123456789abcdef")

	//==========   与参与者1的交互  =======================
	//参与者1发布共享
	id_1 := 1
	skb, vkx, vky := LoadkeyFromFiles(int64(id_1))
	share := Compute_share(C, skb, vkx, vky)

	//user检测参与者1共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_1))
	} else {
		fmt.Printf("\n share for party 1 is false")
	}
	//====================================================

	//==========   与参与者2的交互  =======================
	//参与者2发布共享
	id_2 := 2
	skb, vkx, vky = LoadkeyFromFiles(int64(id_2))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者2共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_2))
	} else {
		fmt.Printf("\n share for party 2 is false")
	}
	//====================================================

	//==========   与参与者3的交互  =======================
	//参与者3发布共享
	id_3 := 3
	skb, vkx, vky = LoadkeyFromFiles(int64(id_3))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_3))
	} else {
		fmt.Printf("\n share for party 3 is false")
	}
	//====================================================

	//==========   与参与者4的交互  =======================
	//参与者4发布共享
	id_4 := 4
	skb, vkx, vky = LoadkeyFromFiles(int64(id_4))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_4))
	} else {
		fmt.Printf("\n share for party 4 is false")
	}
	//====================================================

	//==========   与参与者5的交互  =======================
	//参与者5发布共享
	id_5 := 6
	skb, vkx, vky = LoadkeyFromFiles(int64(id_5))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_5))
	} else {
		fmt.Printf("\n share for party 5 is false")
	}
	//====================================================

	//==========   与参与者5的交互  =======================
	//参与者5发布共享
	id_6 := 7
	skb, vkx, vky = LoadkeyFromFiles(int64(id_6))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_6))
	} else {
		fmt.Printf("\n share for party 6 is false")
	}

	id_7 := 8
	skb, vkx, vky = LoadkeyFromFiles(int64(id_7))
	share = Compute_share(C, skb, vkx, vky)

	//user检测参与者3共享是否正确，如果正确，则存储share，否则丢弃。
	if Verify_share(C, vkx, vky, share) == true {
		Store_share(share, int64(id_7))
	} else {
		fmt.Printf("\n share for party 6 is false")
	}
	//====================================================

	//=============  用户生成随机数 ========================
	idarr := []int64{int64(id_1), int64(id_2), int64(id_3), int64(id_4), int64(id_5), int64(id_6)}
	prf := Compute_prf(idarr, int64(k))
	fmt.Printf("\n prf is %x", prf)
	//===================================================

	idarr = []int64{int64(id_2), int64(id_7), int64(id_1), int64(id_4), int64(id_5), int64(id_6)}
	prf1 := Compute_prf(idarr, int64(k))
	fmt.Printf("\n prf1 is %x", prf1)

	//============== 验证正确型，此步骤仅为验证算法正确型 ======
	id_0 := 0
	skb, vkx, vky = LoadkeyFromFiles(int64(id_0))
	share = Compute_share(C, skb, vkx, vky)
	fmt.Printf("\n share_0 is %x", share[0:64])
	//====================================================

}
