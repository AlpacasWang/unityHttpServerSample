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

	"github.com/labstack/echo"
	"github.com/shamoto-donuts/go/codec"
)

type AnalyzeType int

const (
	OnUserId AnalyzeType = iota
	OnDefault
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
func BodyAnalyze(c echo.Context, analyzeType AnalyzeType) err.ErrWriter {
	ew := err.NewErrWriter()

	// get body
	bodyBuf := new(bytes.Buffer)
	_, err := bodyBuf.ReadFrom(c.Request().Body)
	if err != nil {
		return ew.Write(err)
	}

	// typeごとに分析
	switch analyzeType {
	case OnUserId:
		bodyAnalyzeUserId(c, bodyBuf.Bytes())
	case OnDefault:
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
func bodyAnalyzeUserId(c echo.Context, body []byte) {
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
	c.Set(cKey.UserId, userId)
	c.Set(cKey.Iv, iv)
	c.Set(cKey.CryptData, cryptData)
}

/**************************************************************************************************/
/*!
 *  type : NewUser時のデータ分析
 *
 *  \param   c       : コンテキスト
 *  \param   body    : bodyデータ
 */
/**************************************************************************************************/
func bodyAnalyzeNewUser(c echo.Context, body []byte) {
	// 分割長
	bodyLen := len(body)

	// 暗号化データ
	iv := body[:aes.BlockSize]
	cryptData := body[aes.BlockSize:bodyLen]

	// コンテキストに登録
	c.Set(cKey.Iv, iv)
	c.Set(cKey.CryptData, cryptData)
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
func DecryptAndUnpack(c echo.Context, out interface{}, secretKey string) err.ErrWriter {
	ew := err.NewErrWriter()

	c.Set(cKey.SecretKey, secretKey)
	// 暗号化データ
	iv := c.Get(cKey.Iv).([]byte)
	cryptData := c.Get(cKey.CryptData).([]byte)

	// decrypt
	ci, err := aes.NewCipher([]byte(secretKey))
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
func PackAndEncrypt(c echo.Context, data interface{}) ([]byte, err.ErrWriter) {
	ew := err.NewErrWriter()

	secretKey := c.Get(cKey.SecretKey).(string)

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
	ci, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return []byte(""), ew.Write(err)
	}

	// padding
	in, err := pkcs7Pad(encodeData, aes.BlockSize)
	if err != nil {
		return []byte(""), ew.Write(err)
	}

	// encrypt
	iv := make([]byte, len(secretKey))
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
func ResWrite(c echo.Context, data interface{}) {

	out, ew := PackAndEncrypt(c, data)
	if ew.HasErr() {
		ResError("pack and enctypt error", c, ew)
		return
	}
	c.Blob(http.StatusOK, "application/msgpack; charset=UTF-8", out)
}

/**************************************************************************************************/
/*!
 *  エラー投げる
 */
/**************************************************************************************************/
func ResError(msg string, c echo.Context, ew err.ErrWriter) {
	v := append([]interface{}{msg, ":"}, ew.Err()...)
	fmt.Println("ERROR", v)
	c.JSON(http.StatusInternalServerError, errorResponse{"message": msg})
}
