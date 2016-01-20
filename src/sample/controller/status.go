package controller

type GameErrorStatus int

/**
 * NOTE : ここを更新する場合はクライアントにも反映させること
 */

// 共通のエラー
const (
	NoProblem GameErrorStatus = iota

	MisMatchGameVersion // バージョン違い
	InCorrectData       // 不正なデータ
	LogicInFailure      // 処理途中でエラーが起きた
)

// 固有のエラー
const (
	ErrorUnique GameErrorStatus = 1000 + iota

	NotFoundUuidOrShard // UUID or UserShardが見つからない
	NotFoundUserData    // UserDataが見つからない
)
