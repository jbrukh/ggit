package api

// import (
// 	//	"fmt"
// 	"github.com/jbrukh/ggit/test"
// 	"testing"
// )

// const varRepo = "../test/var"

// var (
// 	linear_history_test_blob_1  = api.OidNow("fd3c81a4d763121f1827c8b4bbdd0cef674c30f9")
// 	linear_history_test_blob_2  = api.OidNow("07d39ab342a978a45857534f30e4a4e2d2ddbf25")
// 	linear_history_test_blob_3  = api.OidNow("083edaac24891a2acc9c1d384cdda1ce7dd3eada")
// 	linear_history_test_blob_4  = api.OidNow("0cfbf08886fca9a91cb753ec8734c84fcbe52c9f")
// 	linear_history_test_blob_5  = api.OidNow("00750edc07d6415dcc07ae0351e9397b0222b7ba")
// 	linear_history_test_blob_6  = api.OidNow("166be640db574d2513aecfde810718f324f529b2")
// 	linear_history_test_blob_7  = api.OidNow("1e8b314962144c26d5e0e50fd29d2ca327864913")
// 	linear_history_test_blob_8  = api.OidNow("1f242fa6f000425d17a7f6c74f77c4908e6b4ef4")
// 	linear_history_test_blob_9  = api.OidNow("269192240915488e73f90552c8c8d83a10cba5df")
// 	linear_history_test_blob_10 = api.OidNow("2e435a26e08d0cb483ca192b3003d030f6e501ee")
// 	linear_history_test_blob_11 = api.OidNow("7ed6ff82de6bcc2a78243fc9c54d3ef5ac14da69")
// 	linear_history_test_blob_12 = api.OidNow("45a4fb75db864000d01701c0f7a51864bd4daabf")
// 	linear_history_test_blob_13 = api.OidNow("49019db807899bc5793047943ce0fbb1a09b2e14")
// 	linear_history_test_blob_14 = api.OidNow("51993f072d5832f20b98b6bd0cf763fb8b4c8a1b")
// 	linear_history_test_blob_15 = api.OidNow("51fdf048b8ac6b09906d589e385475f2ceae18bf")
// 	linear_history_test_blob_16 = api.OidNow("226aaf8af79f3f0bd0ba7766db38233fb0406a58")
// 	linear_history_test_blob_17 = api.OidNow("60276f120ddf8cd0f76acdf0e7681c1c91536e07")
// 	linear_history_test_blob_18 = api.OidNow("6ed281c757a969ffe22f3dcfa5830c532479c726")
// 	linear_history_test_blob_19 = api.OidNow("7290ba859f4adbf90d68526fe0ab1f8cbcf65098")
// 	linear_history_test_blob_20 = api.OidNow("7f8f011eb73d6043d2e6db9d2c101195ae2801f2")
// 	linear_history_test_blob_21 = api.OidNow("91dea2c76e7518c2fa1a625da1cc314bd46f7f05")
// 	linear_history_test_blob_22 = api.OidNow("a5c8806279fa7d6b7d04418a47e21b7e89ab18f8")
// 	linear_history_test_blob_23 = api.OidNow("b2f7f08c17074991c47ce9b2475c3fe58fc26247")
// 	linear_history_test_blob_24 = api.OidNow("b62923296e54bec1a66a8cb71fe025d4166cb9b4")
// 	linear_history_test_blob_25 = api.OidNow("b8626c4cff2849624fb67f87cd0ad72b163671ad")
// 	linear_history_test_blob_26 = api.OidNow("d00491fd7e5bb6fa28c517a0bb32b8b506539d4d")
// 	linear_history_test_blob_27 = api.OidNow("e8183f05f5db68b3934e93f4bf6bed2bb664e0b5")
// 	linear_history_test_blob_28 = api.OidNow("ec635144f60048986bc560c5576355344005e6e7")
// 	linear_history_test_blob_29 = api.OidNow("f599e28b8ab0d8c9c57a486c89c4a5132dcbd3b2")
// 	linear_history_test_blob_30 = api.OidNow("f6ba75da254caa70f552693c14cbbec11e637ad3")
// )

// func Test_readSimpleBlobs(t *testing.T) {
// 	const (
// 		blob2 = "00750edc07d6415dcc07ae0351e9397b0222b7ba"
// 		blob3 = "00750edc07d6415dcc07ae0351e9397b0222b7ba"
// 	)

// 	dir, err := test.Repo(varRepo, "../test/cases/linear_history.sh")
// 	assertNoErr(t, err)

// 	repo := Open(dir)

// 	var o Object
// 	o, err = repo.ObjectFromShortOid(blob2)
// 	assertNoErr(t, err)

// 	assert(t, o.Header().Type() == ObjectBlob)
// 	assert(t, o.Header().Size() == 2)
// 	assert(t, o.(*Blob).ObjectId().String() == blob2)

// 	err = repo.Destroy()
// 	assertNoErr(t, err)
// }
