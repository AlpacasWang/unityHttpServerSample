package controller

/**************************************************************************************************/
/*!
 *  base.go
 *
 *  どのcontrollerでも呼ばれる基底メソッド系
 *  あまりimportさせたくないものをこちらに書く
 *
 */
/**************************************************************************************************/
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"net/http"
	"sample/common/err"
	cKey "sample/conf/context"

	"reflect"

	"encoding/json"

	"github.com/ugorji/go/codec"
	"github.com/zenazn/goji/web"
)

/**
 * 一旦の仮置き！！
 */
const SECRET_KEY = "0123456789ABCDEF"
const USER_KEY   = "FEDCBA9876543210"

type AnalyzeType int

const (
	OnDefault AnalyzeType = iota
	OnLogin
	OnNewUser
)

// エラーはJSONで返す
type errorResponse map[string]interface{}

/**************************************************************************************************/
/*!
 *  リクエストbodyのデータ分け
 *
 *  \param   c : コンテキスト
 *  \param   r : リクエスト情報
 *  \return  エラー情報
 */
/**************************************************************************************************/
func BodyAnalyze(c *web.C, r *http.Request, analyzeType AnalyzeType) err.ErrWriter {
	ew := err.NewErrWriter()

	// get body
	bodyBuf := new(bytes.Buffer)
	_, err := bodyBuf.ReadFrom(r.Body)
	if err != nil {
		return ew.Write(err)
	}

	// typeごとに分析
	switch analyzeType {
	case OnDefault:
		bodyAnalyzeDefault(c, bodyBuf.Bytes())
	case OnLogin:
		bodyAnalyzeLogin(c, bodyBuf.Bytes())
	case OnNewUser:
		bodyAnalyzeNewUser(c, bodyBuf.Bytes())
	default:
		ew.Write("undefine analyze type!!")
	}

	return ew
}

/**************************************************************************************************/
/*!
 *  type : Default時のデータ分析
 *
 *  \param   c       : コンテキスト
 *  \param   body    : bodyデータ
 */
/**************************************************************************************************/
func bodyAnalyzeDefault(c *web.C, body []byte) {
	// 分割長
	userIdPartLen := 2
	bodyLen := len(body)

	// bodyをユーザIDと、暗号化データに分ける
	// バイト端のデータを結合
	uBytes1 := body[:userIdPartLen]
	uBytes2 := body[bodyLen-userIdPartLen : bodyLen]
	uBytes := append(uBytes2, uBytes1...)

	// convert to user_id
	userId := binary.LittleEndian.Uint32(uBytes)

	// 暗号化データ
	iv := body[userIdPartLen : userIdPartLen+aes.BlockSize]
	cryptData := body[userIdPartLen+aes.BlockSize : bodyLen-userIdPartLen]

	// コンテキストに登録
	c.Env[cKey.UserId] = userId
	c.Env[cKey.Iv] = iv
	c.Env[cKey.CryptData] = cryptData
}

/**************************************************************************************************/
/*!
 *  type : Login時のデータ分析
 *
 *  \param   c       : コンテキスト
 *  \param   body    : bodyデータ
 */
/**************************************************************************************************/
func bodyAnalyzeLogin(c *web.C, body []byte) {
	// 分割長
	userIdPartLen := 18
	bodyLen := len(body)

	// bodyをユーザIDと、暗号化データに分ける
	// バイト端のデータを結合
	uBytes1 := body[:userIdPartLen]
	uBytes2 := body[bodyLen-userIdPartLen : bodyLen]
	uBytes := append(uBytes2, uBytes1...)

	// convert to uuid
	uuid := string(uBytes)

	// 暗号化データ
	iv := body[userIdPartLen : userIdPartLen+aes.BlockSize]
	cryptData := body[userIdPartLen+aes.BlockSize : bodyLen-userIdPartLen]

	// コンテキストに登録
	c.Env[cKey.Uuid] = uuid
	c.Env[cKey.Iv] = iv
	c.Env[cKey.CryptData] = cryptData
}

/**************************************************************************************************/
/*!
 *  type : NewUser時のデータ分析
 *
 *  \param   c       : コンテキスト
 *  \param   body    : bodyデータ
 */
/**************************************************************************************************/
func bodyAnalyzeNewUser(c *web.C, body []byte) {
	// 分割長
	bodyLen := len(body)

	// 暗号化データ
	iv := body[:aes.BlockSize]
	cryptData := body[aes.BlockSize:bodyLen]

	// コンテキストに登録
	c.Env[cKey.Iv] = iv
	c.Env[cKey.CryptData] = cryptData
}

/**************************************************************************************************/
/*!
 *  暗号化されたデータの解析
 *
 *  \param   c       : コンテキスト
 *  \param   out     : データ格納先
 *  \param   keyArgs : security_key情報（予定）
 *  \return  エラー情報
 */
/**************************************************************************************************/
func DecryptAndUnpack(c web.C, out interface{}, keyArgs ...string) err.ErrWriter {
	ew := err.NewErrWriter()

	// 暗号化データ
	iv := c.Env[cKey.Iv].([]byte)
	cryptData := c.Env[cKey.CryptData].([]byte)

	// TODO : secret keyはログイン時に生成したものを利用する

	// decrypt
	ci, err := aes.NewCipher([]byte(SECRET_KEY))
	if err != nil {
		return ew.Write(err)
	}
	cbcDecrypter := cipher.NewCBCDecrypter(ci, iv)

	plain := make([]byte, len(cryptData))
	cbcDecrypter.CryptBlocks(plain, cryptData)



	// decode(codec)

		mh := &codec.MsgpackHandle{RawToString: true}
		dec := codec.NewDecoderBytes(plain, mh)
		err = dec.Decode(out)


	if err != nil {
		return ew.Write(err)
	}

	return ew
}

/**************************************************************************************************/
/*!
 *  送信するデータを暗号化する
 *
 *  \param   data    : 暗号化するデータ
 *  \param   keyArgs : security_key情報（予定）
 *  \return  暗号化済データ、エラー情報
 */
/**************************************************************************************************/
func PackAndEncrypt(data interface{}, keyArgs ...string) ([]byte, err.ErrWriter) {
	ew := err.NewErrWriter()

	// pkcs7 padding function
	pkcs7Pad := func(packed []byte, blockLength int) ([]byte, error) {
		if blockLength <= 0 {
			return nil, fmt.Errorf("invalid block-length %d", blockLength)
		}

		padLen := blockLength - (len(packed) % blockLength)

		pad := bytes.Repeat([]byte{byte(padLen)}, padLen)
		return append(packed, pad...), nil
	}

	var encodeData []byte

	// encode(codec)

		mh := &codec.MsgpackHandle{}
		mh.MapType = reflect.TypeOf(data)
		encoder := codec.NewEncoderBytes(&encodeData, mh)
		e := encoder.Encode(data)
		if e != nil {
			return []byte(""), ew.Write(e)
		}


	// new cipher
	ci, err := aes.NewCipher([]byte(SECRET_KEY))
	if err != nil {
		return []byte(""), ew.Write(err)
	}

	// padding
	in, err := pkcs7Pad(encodeData, aes.BlockSize)
	if err != nil {
		return []byte(""), ew.Write(err)
	}

	// encrypt
	iv := make([]byte, len(SECRET_KEY))
	out := make([]byte, len(in))
	cbcEncrypter := cipher.NewCBCEncrypter(ci, iv)
	cbcEncrypter.CryptBlocks(out, in)

	encryptData := append(iv, out...)
	return encryptData, ew
}

/**************************************************************************************************/
/*!
 *  データを返す
 *  想定済みのエラーもこちらで返す
 */
/**************************************************************************************************/
func ResWrite(data interface{}, w http.ResponseWriter) {

	out, ew := PackAndEncrypt(data)
	if ew.HasErr() {
		ResError(InCorrectData, "pack and enctypt error", w, ew)
		return
	}
	w.Header().Set("Content-Type", "application/msgpack; charset=UTF-8")
	w.Write(out)
}

/**************************************************************************************************/
/*!
 *  エラー投げる
 */
/**************************************************************************************************/
func ResError(status GameErrorStatus, msg string, w http.ResponseWriter, ew err.ErrWriter) {
	v := append([]interface{}{msg, ":"}, ew.Err()...)
	fmt.Println("ERROR", v)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(errorResponse{"status": status, "message": msg})
}
