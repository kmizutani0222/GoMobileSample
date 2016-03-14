package httpUtil

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"strconv"
	"time"
)

var DEBUG_NUM = 5
var DEBUG_KBN = 0
//  API基本設定盲目
var api_device = ""
var api_build_no = ""
var api_version_no = ""
var api_key = ""

// プロフィール編集用項目
var api_profile_update_job_id = ""
var api_profile_update_home_prefecture_id = ""
var api_profile_update_prefecture_id = ""
var api_profile_update_avatar = ""
var api_profile_update_hobby = []string{}
var api_profile_update_income = ""
var api_profile_update_height = ""
var api_profile_update_smoking = ""
var api_profile_update_alcohol = ""
var api_profile_update_pr = ""
var api_profile_update_blood = ""
var api_profile_update_body_type = ""
var api_profile_update_educational = ""
var api_profile_update_housemate = ""
var api_profile_update_brothers = ""
var api_profile_update_country = ""
var api_profile_update_language = []string{}
var api_profile_update_children = ""
var api_profile_update_marriage_intention = ""
var api_profile_update_wants_children = ""
var api_profile_update_housework = ""
var api_profile_update_matching_type = ""
var api_profile_update_first_date_payment = ""
var api_profile_update_personality = []string{}
var api_profile_update_sociability = ""
var api_profile_update_holiday = ""
var api_user_search_prefectures = []string{}

type BasicAuth struct{
	Id string
	Pass string
}

const (
	api_user_badge_get_url = "/api/v5/user/badge/get"
	api_user_badge_update_url = "/api/v5/user/badge/update"
	api_user_profile_url = "/api/v5/user/profile"
	api_user_search_list_url = "/api/v5/user/search/list"
	api_user_post_list_url = "/api/v5/user/post/list"
	api_user_post_insert_url = "/api/v5/user/post/insert"
	api_user_post_update_url = "/api/v5/user/post/update"
	api_user_follow_list_url = "/api/v5/user/follow/list"
	api_user_like_list_url = "/api/v5/user/like/list"
	api_public_send_document_url = "/api/v5/public/send_document"
	api_public_prefecture_list_url = "/api/v5/public/prefecture/list"
	api_public_stations_url = "/api/v5/public/stations"
	api_public_release_url = "/api/v5/public/releases"
	api_public_generate_api_key_url = "/api/v5/public/generate_api_key"
	api_account_profile_update_url = "/api/v5/account/profile/update"
	api_account_init_url = "/api/v5/account/init"
	api_approach_like_url = "/api/v5/approach/like"
	api_approach_follow_url = "/api/v5/approach/follow"
	api_approach_accept_url = "/api/v5/approach/accept"
	api_approach_request_url = "/api/v5/approach/request"
	api_approach_block_url = "/api/v5/approach/block"
	api_approach_penalty_url = "/api/v5/approach/penalty"
	api_approach_list_url ="/api/v5/approach/list"
	api_approach_talk_list_url ="/api/v5/approach/talk/list"
	api_media_download_url = "/api/v5/media/download"
	api_media_upload_url = "/api/v5/media/upload"
	api_spot_list_url = "/api/v5/spot/list"
	api_spot_search_url = "/api/v5/spot/search"
	api_spot_detail_url = "/api/v5/spot/detail"
	api_spot_genres_url = "/api/v5/spot/genres"
	api_talk_polling_url = "/api/v5/talk/polling"
	api_talk_diff_url = "/api/v5/talk/diff"
	api_talk_send_url = "/api/v5/talk/send"
)

// フェイスシーン向けのAPIを使用する場合、最初に呼び出す
func FaceSceneInitialize(device string, build_no string) {
	api_device = device
	api_build_no = build_no
}

func SetApiKey(key string) {
	api_key = key
}

// バッヂ取得API呼び出し
func ApiUserBadgeGet() string {
	params := setDefaultParams()
	return httpPost(api_user_badge_get_url, params)
}

// バッヂ取得API呼び出し
func ApiUserBadgeUpdate(new_talks int, new_talk_id int, new_activities int, new_news int) string {
	params := setDefaultParams()
	params["new_talks"] = string(new_talks)
	params["new_talk_id"] = string(new_talk_id)
	params["new_activities"] = string(new_activities)
	params["new_news"] = string(new_news)
	return httpPost(api_user_badge_update_url, params)
}

// ユーザー検索結果一覧取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiUserSearch(api_url string, prefectures string, age_from string, age_to string, income_from string, income_to string, smoking string) string {
	params := setDefaultParams()
	if prefectures > "" {
		params["prefectures"] = prefectures
	}
	if age_from > "" {
		params["age[from]"] = age_from
	}
	if age_to > "" {
		params["age[to]"] = age_to
	}
	if income_from > "" {
		params["income[from]"] = income_from
	}
	if income_to > "" {
		params["income[to]"] = income_to
	}
	if smoking > "" {
		params["smoking"] = smoking
	}
	if api_url <= "" {
		api_url = api_user_search_list_url
	}
	return httpPost(api_url, params)
}

// 投稿一覧API呼び出し
func ApiUserPostList(api_url string, type_code string, id string) string {
	params := setDefaultParams()
	params["type"] = type_code
	params["id"] = id
	if api_url <= "" {
		api_url = api_user_post_list_url
	}
	return httpPost(api_url, params)
}

// 投稿API呼び出し
func ApiUserPostInsert(content string, media string) string {
	params := setDefaultParams()
	params["content"] = content
	params["media"] = media
	return httpPost(api_user_post_insert_url, params)
}

// 投稿更新API呼び出し
func ApiUserPostUpdate(id string, content string, media string, delete_flg bool) string {
	params := setDefaultParams()
	params["id"] = id
	params["content"] = content
	params["media"] = media
	if delete_flg {
		params["delete_media"] = "true"
	}
	return httpPost(api_user_post_update_url, params)
}

// いいね一覧取得API呼び出し
func ApiUserLikeList(api_url string, type_code string, id string) string {
	params := setDefaultParams()
	params["type"] = type_code
	params["id"] = id
	if api_url <= "" {
		api_url = api_user_like_list_url
	}
	return httpPost(api_url, params)
}

// フォロー一覧API呼び出し
func ApiUserFollowList(api_url string, type_code string) string {
	params := setDefaultParams()
	params["type"] = type_code
	if api_url <= "" {
		api_url = api_user_follow_list_url
	}
	return httpPost(api_url, params)
}

// ユーザー情報取得API呼び出し
func ApiUserProfile(id string) string {
	params := setDefaultParams()
	params["id"] = id
	return httpPost(api_user_profile_url, params)
}

// 独身証明書申込書取得API呼び出し
func ApiPublicSendDocument(token string, user_id string) string {
	params := setDefaultParams()
	params["token"] = token
	params["user"] = user_id
	return httpPost(api_public_send_document_url, params)
}

// 都道府県取得API呼び出し
func ApiPublicPrefectureList() string {
	params := setDefaultParams()
	return httpPost(api_public_prefecture_list_url, params)
}

// 駅検索（オートコンプリート）API呼び出し
func ApiPublicStations(name string) string {
	params := setDefaultParams()
	params["name"] = name
	return httpPost(api_public_stations_url, params)
}

// バージョンリリース状況取得API呼び出し
func ApiPublicRelease(version_no string) string {
	params := setDefaultParams()
	params["version_no"] = version_no
	return httpPost(api_public_release_url, params)
}

// APIキー取得API呼び出し（フェイスブック認証後に呼び出す）
func ApiPublicGenerateApiKey(device_token string, fb_id string, fb_token string, fb_token_expires string) string {
	params := setDefaultParams()
	params["device_token"] = device_token
	params["fb_id"] = fb_id
	params["fb_access_token"] = fb_token
	params["fb_token_expired_time"] = fb_token_expires
	return httpPost(api_public_generate_api_key_url, params)

}

// ユーザープロフィール更新準備(パラメーターの初期化)
func InitializeProfileUpdate() {
	api_profile_update_job_id = "";
	api_profile_update_home_prefecture_id = "";
	api_profile_update_prefecture_id = "";
	api_profile_update_avatar = "";
	api_profile_update_hobby = []string{};
	api_profile_update_income = "";
	api_profile_update_height = "";
	api_profile_update_smoking = "";
	api_profile_update_alcohol = "";
	api_profile_update_pr = "";
	api_profile_update_blood = "";
	api_profile_update_body_type = "";
	api_profile_update_educational = "";
	api_profile_update_housemate = "";
	api_profile_update_brothers = "";
	api_profile_update_country = "";
	api_profile_update_language = []string{};
	api_profile_update_children = "";
	api_profile_update_marriage_intention = "";
	api_profile_update_wants_children = "";
	api_profile_update_housework = "";
	api_profile_update_matching_type = "";
	api_profile_update_first_date_payment = "";
	api_profile_update_personality = []string{};
	api_profile_update_sociability = "";
	api_profile_update_holiday = "";
}

// ユーザープロフィールの職種を設定
func SetProfileUpdateJob(job_id string) {
	api_profile_update_job_id = job_id;
}

// ユーザープロフィールのm出身地を設定
func SetProfileUpdateHomePrefecture(home_prefecture_id string) {
	api_profile_update_home_prefecture_id = home_prefecture_id;
}

// ユーザープロフィールの居住地を設定
func SetProfileUpdatePrefecture(prefecture_id string) {
	api_profile_update_prefecture_id = prefecture_id;
}

// ユーザープロフィールの写真を設定
func SetProfileUpdateAvatar(avatar string) {
	api_profile_update_avatar = avatar;
}

// ユーザープロフィールの趣味を追加
func AddProfileUpdateHobby(hobby_id string) {
	api_profile_update_hobby = append(api_profile_update_hobby, hobby_id);
}

// ユーザープロフィールの収入を設定
func SetProfileUpdateIncome(income string) {
	api_profile_update_income = income;
}

// ユーザープロフィールの身長を設定
func SetProfileUpdateHeight(height string) {
	api_profile_update_height = height;
}

// ユーザープロフィールの喫煙事情を設定
func SetProfileUpdateSmoking(smoking string) {
	api_profile_update_smoking = smoking;
}

// ユーザープロフィールの飲酒事情を設定
func SetProfileUpdateAlcohol(alcohol string) {
	api_profile_update_alcohol = alcohol;
}

// ユーザープロフィールの自己紹介を設定
func SetProfileUpdatePr(pr string) {
	api_profile_update_pr = pr;
}

// ユーザープロフィールの血液型を設定
func SetProfileUpdateBlood(blood string) {
	api_profile_update_blood = blood;
}

// ユーザープロフィールの体型を設定
func SetProfileUpdateBodyType(body_type string) {
	api_profile_update_body_type = body_type;
}

// ユーザープロフィールの学歴を設定
func SetProfileUpdateEducational(educational string) {
	api_profile_update_educational = educational;
}

// ユーザープロフィールの同居人を設定
func SetProfileUpdateHousemate(housemate string) {
	api_profile_update_housemate = housemate;
}

// ユーザープロフィールの兄弟姉妹を設定
func SetProfileUpdateBrothers(brothers string) {
	api_profile_update_brothers = brothers;
}

// ユーザープロフィールの国籍を設定
func SetProfileUpdateCountry(country string) {
	api_profile_update_country = country;
}

// ユーザープロフィールの言語を追加
func AddProfileUpdateLanguage(language_id string) {
	api_profile_update_language = append(api_profile_update_language, language_id);
}

// ユーザープロフィールの子供の有無を設定
func SetProfileUpdateChildren(children string) {
	api_profile_update_children = children;
}

// ユーザープロフィールの結婚の意思を設定
func SetProfileUpdateMarriageIntention(marriage_intention string) {
	api_profile_update_marriage_intention = marriage_intention;
}

// ユーザープロフィールの子供が欲しいかを設定
func SetProfileUpdateWantsChildren(wants_children string) {
	api_profile_update_wants_children = wants_children;
}

// ユーザープロフィールの家事・育児を設定
func SetProfileUpdateHousework(housework string) {
	api_profile_update_housework = housework;
}

// ユーザープロフィールの出会うまでの希望を設定
func SetProfileUpdateMatchingType(matching_type string) {
	api_profile_update_matching_type = matching_type;
}

// ユーザープロフィールの初回デート費用を設定
func SetProfileUpdateFirstDatePayment(first_date_payment string) {
	api_profile_update_first_date_payment = first_date_payment;
}

// ユーザープロフィールの性格を追加
func AddProfileUpdatePersonality(personality_id string) {
	api_profile_update_personality = append(api_profile_update_personality, personality_id);
}

// ユーザープロフィールの社交性を設定
func SetProfileUpdateSociability(sociability string) {
	api_profile_update_sociability = sociability;
}

// ユーザープロフィールの社交性を設定
func SetProfileUpdateHoliday(holiday string) {
	api_profile_update_holiday = holiday;
}

// プロフィール編集API呼び出し
func ApiAccountProfileUpdate() string {
	params := setDefaultParams()
	// 職種が空文字でなければ設定
	if api_profile_update_job_id != "" {
		params["job_id"] = api_profile_update_job_id;
	}
	// 出身地が空文字でなければ設定
	if api_profile_update_home_prefecture_id != "" {
		params["home_prefecture_id"] = api_profile_update_home_prefecture_id;
	}
	// 居住地が空文字でなければ設定
	if api_profile_update_prefecture_id != "" {
		params["prefecture_id"] = api_profile_update_prefecture_id;
	}
	// プロフィール写真が空文字でなければ設定
	if api_profile_update_avatar != "" {
		params["avatar"] = api_profile_update_avatar;
	}
	// 趣味を登録された分だけ回す
	for key, val := range api_profile_update_hobby {
		params["hobby[" + strconv.Itoa((key+1)) + "]"] = val;
	}
	// 収入が空文字でなければ設定
	if api_profile_update_income != "" {
		params["income"] = api_profile_update_income;
	}
	// 収入が空文字でなければ設定
	if api_profile_update_income != "" {
		params["height"] = api_profile_update_height;
	}
	// 喫煙が空文字でなければ設定
	if api_profile_update_smoking != "" {
		params["smoking"] = api_profile_update_smoking;
	}
	// お酒が空文字でなければ設定
	if api_profile_update_alcohol != "" {
		params["alcohol"] = api_profile_update_alcohol;
	}
	// 自己紹介が空文字でなければ設定
	if api_profile_update_pr != "" {
		params["pr"] = api_profile_update_pr;
	}
	// 血液型が空文字でなければ設定
	if api_profile_update_blood != "" {
		params["blood"] = api_profile_update_blood;
	}
	// 体型が空文字でなければ設定
	if api_profile_update_body_type != "" {
		params["body_type"] = api_profile_update_body_type;
	}
	// 学歴が空文字でなければ設定
	if api_profile_update_educational != "" {
		params["educational"] = api_profile_update_educational;
	}
	// 同居人が空文字でなければ設定
	if api_profile_update_housemate != "" {
		params["housemate"] = api_profile_update_housemate;
	}
	// 兄弟・姉妹が空文字でなければ設定
	if api_profile_update_brothers != "" {
		params["brother"] = api_profile_update_brothers;
	}
	// 国籍が空文字でなければ設定
	if api_profile_update_country != "" {
		params["country"] = api_profile_update_country;
	}
	// 言語を登録された分だけ、回す
	for key, val := range api_profile_update_language {
		params["language[" + strconv.Itoa((key + 1)) + "]"] = val;
	}
	// 子供の有無が空文字でなければ設定
	if api_profile_update_children != "" {
		params["children"] = api_profile_update_children;
	}
	// 結婚の意思が空文字でなければ設定
	if api_profile_update_marriage_intention != "" {
		params["marriage_intention"] = api_profile_update_marriage_intention;
	}
	// 子供が欲しいかが空文字でなければ設定
	if api_profile_update_wants_children != "" {
		params["wants_children"] = api_profile_update_wants_children;
	}
	// 家事・育児が空文字でなければ設定
	if api_profile_update_housework != "" {
		params["housework"] = api_profile_update_housework;
	}
	// 出会うまでの希望が空文字でなければ設定
	if api_profile_update_matching_type != "" {
		params["matching_type"] = api_profile_update_matching_type;
	}
	// 初回デート費用が空文字でなければ設定
	if api_profile_update_first_date_payment != "" {
		params["first_date_payment"] = api_profile_update_first_date_payment;
	}
	// 性格を登録された分だけ回す
	for key, val := range api_profile_update_personality {
		params["personality[" + strconv.Itoa((key + 1)) + "]"] = val;
	}
	// 社交性が空文字でなければ設定
	if api_profile_update_sociability != "" {
		params["sociability"] = api_profile_update_sociability;
	}
	// 休日が空文字でなければ設定
	if api_profile_update_holiday != "" {
		params["holiday"] = api_profile_update_holiday;
	}
	return httpPost(api_account_profile_update_url, params)
}

// 新規登録API呼び出し
func ApiAccountInit(avatar_url string, prefecture string, token string) string {
//	return "{\"header\":{\"code\":200, \"message\":\"正常\", \"released\":0, \"isNew\":0}, \"body\":{\"api_key\":\"1234567890\", \"is_new_user\":0, \"users\":{\"id\":1286,\"type\":\"TYPE_APP\",\"avatar_url\":\"https://s3-ap-northeast-1.amazonaws.com/strage-face-scene-jp/face-scene.jp/prod/eb/a7/caebe25ed4e7fe2a7504ddbaa15a2637b817?X-Amz-Content-Sha256=e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855&X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIAI6U34YFM2IRMFPRA%2F20160113%2Fap-northeast-1%2Fs3%2Faws4_request&X-Amz-Date=20160113T133527Z&X-Amz-SignedHeaders=Host&X-Amz-Expires=240&X-Amz-Signature=44bebebf8fff37a9436b20a6fb6154d9c574ed31b0b9942f94210d550fe59103\",\"pr\":\"\",\"gender\":1,\"prefecture\":\"北海道\",\"age\":\"20\",\"income\":\"800万〜1000万\",\"hobbies\":[\"メジャーリーグ観戦\",\"ダーツ\",\"絵画\"],\"height\":131,\"smoking\":true,\"status\":{\"follow\":3,\"omiai\":1},\"talk_tickets\":{},\"certificate_status\":\"99\",\"follow_count\":6,\"follower_count\":4,\"liked_count\":4,\"post_count\":7}}}"
	params := setDefaultParams()
	params["image"] = avatar_url
	params["prefecture_id"] = prefecture
	if (token != "") {
		params["code"] = token
	}
	return httpPost(api_account_init_url, params)
}

// いいねAPI呼び出し
func ApiApproachLike(post_id string) string {
	params := setDefaultParams()
	params["id"] = post_id
	return httpPost(api_approach_like_url, params)
}

// フォローAPI呼び出し
func ApiApproachFollow(user_id string) string {
	params := setDefaultParams()
	params["id"] = user_id
	return httpPost(api_approach_follow_url, params)
}

// トーク承認API呼び出し
func ApiApproachAccept(approach_id string) string {
	params := setDefaultParams()
	params["id"] = approach_id
	return httpPost(api_approach_accept_url, params)
}

// トーク申請API呼び出し
func ApiApproachRequest(user_id string) string {
	params := setDefaultParams()
	params["id"] = user_id
	return httpPost(api_approach_request_url, params)
}

// ユーザーブロックAPI呼び出し
func ApiApproachBlock(user_id string) string {
	params := setDefaultParams()
	params["id"] = user_id
	return httpPost(api_approach_block_url, params)
}

// ユーザー違反報告API呼び出し
func ApiApproachPenalty(user_id string, type_id string, content string) string {
	params := setDefaultParams()
	params["id"] = user_id
	params["type"] = type_id
	params["content"] = content
	return httpPost(api_approach_penalty_url, params)
}

// アプローチのあったユーザーの一覧取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiApproachList(api_url string) string {
	params := setDefaultParams()
	if api_url <= "" {
		api_url = api_approach_list_url
	}
	return httpPost(api_url, params)
}

// トーク中の一覧取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiApproachTalkList(api_url string) string {
	params := setDefaultParams()
	if api_url <= "" {
		api_url = api_approach_talk_list_url
	}
	return httpPost(api_url, params)
}

// トーク中のAPI呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiMediaDownload(media_type string) string {
	params := setDefaultParams()
	params["media_type"] = media_type
	return httpPost(api_media_download_url, params)
}

// トーク中のAPI呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiMediaUpload(media_type string, media1 string, media2 string) string {
	params := setDefaultParams()
	params["media_type"] = media_type
	params["media1"] = media1
	_, err := os.Open(media2)
	// エラーじゃない場合ファイル
	if err == nil {
		params["media2"] = media2
	}
	return httpPost(api_media_upload_url, params)
}

// おすすめスポット一覧取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiSpotList(api_url string) string {
	params := setDefaultParams()
	if api_url <= "" {
		api_url = api_spot_list_url
	}
	return httpPost(api_url, params)
}

// おすすめスポットの検索結果一覧取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiSpotSearch(api_url string, genre string, area string) string {
	params := setDefaultParams()
	if genre > "" {
		params["genre"] = genre
	}
	if area > "" {
		params["area"] = area
	}
	if api_url <= "" {
		api_url = api_spot_search_url
	}
	return httpPost(api_url, params)
}

// おすすめスポットの詳細取得API呼び出し
func ApiSpotDetail(id string) string {
	params := setDefaultParams()
	params["id"] = id
	return httpPost(api_spot_detail_url, params)
}

// おすすめスポットの詳細取得API呼び出し
func ApiSpotGenres() string {
	params := setDefaultParams()
	return httpPost(api_spot_genres_url, params)
}

// 新着メッセージ取得API呼び出し（ここはページング有りのため、URLは渡してもらうようにしておく）
func ApiTalkPolling(api_url string, id string) string {
	params := setDefaultParams()
	params["id"] = id
	if api_url <= "" {
		api_url = api_talk_polling_url
	}
	return httpPost(api_url, params)
}

// メッセージ差分取得API呼び出し
func ApiTalkDiff(chat_id string, message_id string) string {
	params := setDefaultParams()
	params["chat_id"] = chat_id
	params["message_id"] = message_id
	return httpPost(api_talk_diff_url, params)
}

// メッセージ送信API呼び出し
func ApiTalkSendMessage(id string, type_code string, text string) string {
	params := setDefaultParams()
	params["id"] = id
	params["type"] = type_code
	params["text"] = text
	return httpPost(api_talk_send_url, params)
}

// おすすめスポット提案送信API呼び出し
func ApiTalkSendPlace(id string, type_code string, place_type string, place_id string, place_group_code string) string {
	params := setDefaultParams()
	params["id"] = id
	params["type"] = type_code
	params["place_type"] = place_type
	params["place_id"] = place_id
	params["place_group_code"] = place_group_code
	return httpPost(api_talk_send_url, params)
}

func setDefaultParams() map[string]string {
	return map[string]string{"device":api_device, "build_no":api_build_no, "apikey":api_key}
}

func httpPost(api_url string, params map[string]string) string {
	fmt.Println("request_url : " + getHttpUrl(api_url))

	// 送信用パラメーターを作成
	bparams := &bytes.Buffer{}
	writer := multipart.NewWriter(bparams)

	// パラメーターを書き込んでいく
	for key, val := range params {
		fmt.Println(key + ":" + val)
		file, err := os.Open(val)
		// エラーだった場合ファイルではない
		if err != nil {
			// 通常のパラメーターとして書き込み
			_ = writer.WriteField(key, val)
		} else {
			fmt.Println("is File")
			// ファイルとして書き込み
			part, err := writer.CreateFormFile(key, filepath.Base(val))
			if err != nil {
				return getErrorJSON(901, "ファイルの送信に失敗しました")
			}
			_, err = io.Copy(part, file)
		}
	}

	err := writer.Close()
	if err != nil {
		return getErrorJSON(902, "通信に失敗しました")
	}

	req, err := http.NewRequest(
        "POST",
        getHttpUrl(api_url),
        bparams)

	if err != nil {
		return getErrorJSON(903, "通信に失敗しました")
	}

	// 持続接続を設定
	req.Header.Set("Connection", "Keep-Alive")

	// Content-Type 設定
	req.Header.Add("Content-Type", writer.FormDataContentType())

    // Basic認証
//    req.Header.Add("Authorization", getBasicAuth(auth_id, auth_pass))

    // タイムアウトを15秒に設定
    client := &http.Client{ Timeout: time.Duration(15 * time.Second) }
    resp, err := client.Do(req)
    if err != nil {
    	fmt.Println(err.Error())
    	return getErrorJSON(500, "通信に失敗しました")
    }

    defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return getErrorJSON(904, "通信に失敗しました")
	}

	// ステータスによるエラーハンドリング(正常終了じゃなければ)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%d", resp.StatusCode)
		return getErrorJSON(resp.StatusCode, resp.Status)
    }

	return string(body)
}

// URLにプロトコルとドメインを追加する
func getHttpUrl(url string) string {
	// URLが空だった場合
	if len(url) <= 0 {
		return url;
	}

	HTTP := [7]string{"https://", "http://", "http://", "http://", "http://", "http://", "http://"}
	HOST_NAME := [7]string{"facescene.jp", "staging.facescene.jp", "test.facescene.jp", "demo.facescene.jp", "jack-russell-terrier.face-scene.okinawa", "siberian-husky.face-scene.okinawa/app_dev.php", "pallid-scops-owl.face-scene.okinawa"}

	// urlにhttp://もしくはhttps://が含まれていない場合
	if !strings.Contains(url, "http://") && !strings.Contains(url, "https://") {
		return HTTP[DEBUG_KBN] + HOST_NAME[DEBUG_KBN] + url
	}

	return url
}

// Base64に変換し、Basic認証用のパスを作る
func getBasicAuth(key string, pass string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(key) + ":" + url.QueryEscape(pass)))
}

// エラーデータの作成
func getErrorJSON(code int, message string) string {
	return "{\"header\":{\"code\":" + strconv.Itoa(code) + ", \"message\":\"" + message + "\", \"released\":0, \"isNew\":0}, \"body\":{}}"
}
