# GoMobileSample
GoMobileで作ったModuleとその周辺のサンプルコード

# 構成
- GomobileAndroidSample　Androidのソースコード
- httpRequest Androidに組み込むGomobileで作ったmoduleとコード
- server リクエストを受け取るサーバーのソースコード

# 以下LTで使った資料

# GoMobileでModuleを作った話
株式会社アスクリード
水谷　健太

Twitter : @mizutani0222
Facebook : kenta.mizutani.56
Androidエンジニア
PHPとかちょっとかじってる

# なんで？
iOSとAndroidで共通化できたらいいなって思った
~~iOSのライブラリ作れたら、iOS覚える気になると思った~~

# やったこと
WebAPIを呼び出して、受け取ったJSONを返す

# これだけ！

# 良かったこと
サーバー側のAPIが出来てなくても、ライブラリで適当な戻り値を返してあげれば、Android本体のコードにまったく影響がなかった。

# ビビってること
まだ実験段階なので、やっぱりやめたとか言われたら泣ける

# 苦労したこと

# 検索しづらい 
 [郷ひろみ公式 GO- MOBILE](http://hiromi-go.net.fanmo.jp/pc/)とか出てくる
<img width="478" alt="スクリーンショット 2016-03-14 22.19.53.png (636.1 kB)" src="https://img.esa.io/uploads/production/attachments/2294/2016/03/14/6362/10930dbb-9027-4713-94cd-1d2e913eb8d6.png">

# これだけ！

# 本当に苦労したこと
publicなメソッドの引数や戻り値に配列が渡せないこと

# サンプルコード
## リクエストのパラメーターを指定するところ
``` postSample.go
func SetParams(name string, val string) {
	params[name] = val;
}
```
## リクエストパラメーターを設定
``` postSample.go
func HttpPost(api_url string) string {
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
```
## リクエストを作る
``` postSample.go
	req, err := http.NewRequest(
        "POST",
        getHttpUrl(api_url),
        bparams)

	if err != nil {
		return getErrorJSON(903, "通信に失敗しました")
	}
```
## ヘッダーの設定 
``` postSample.go
	// 持続接続を設定
	req.Header.Set("Connection", "Keep-Alive")

	// Content-Type 設定
	req.Header.Add("Content-Type", writer.FormDataContentType())

    // Basic認証がある場合はコメントを外す
//    req.Header.Add("Authorization", getBasicAuth(auth_id, auth_pass))
```
## 送信部分
``` postSample.go
    // タイムアウトを15秒に設定
    client := &http.Client{ Timeout: time.Duration(15 * time.Second) }
    resp, err := client.Do(req)
    if err != nil {
    	fmt.Println(err.Error())
    	return getErrorJSON(500, err.Error())
    }

    defer resp.Body.Close()
```
## 受け取った情報をチェック
``` postSample.go
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return getErrorJSON(904, err.Error())
	}

	// ステータスによるエラーハンドリング(正常終了じゃなければ)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%d", resp.StatusCode)
		return getErrorJSON(resp.StatusCode, resp.Status)
    }

	return string(body)
}
```

# Androidに組み込む
1. メニュー File > New > New Module... を選択、Import .JAR/.AAR Package でライブラリのモジュールを作る。
1. メニュー File > Project Structure... を選択、 ライブラリを使用するアプリケーションのモジュールを選択し、 「Dependencies」のタブを開き、 プラス記号「＋」アイコンを押し「Module Dependency」を選ぶ。 手順1で追加したモジュールを選択。

#  コード
``` MainActivity.java
switch (v.getId()) {
    // リクエスト１ボタンが押された
    case R.id.request1_btn:
        PostSample.SetParams("param1", param1.getText().toString());
        response.setText(PostSample.HttpPost("/apisample"));
         break;
    // リクエスト２ボタンが押された
    case R.id.request2_btn:
        PostSample.SetParams("param2", param2.getText().toString());
        response.setText(PostSample.HttpPost("/apisample"));
         break;
}
```

# こんな感じの画面
<img width="300" alt="Screenshot_2016-03-15-19-39-13.png (85.5 kB)" src="https://img.esa.io/uploads/production/attachments/2294/2016/03/15/6362/d380af33-a4a4-4131-b415-d419ba29f462.png">


# 動いた！
<img width="300" alt="Screenshot_2016-03-15-19-15-22.png (90.5 kB)" src="https://img.esa.io/uploads/production/attachments/2294/2016/03/15/6362/7fa71c45-6ed9-442d-b90b-e27d42251b30.png">



# なんかおかしい
<img width="300" alt="Screenshot_2016-03-15-19-15-22.png (120.2 kB)" src="https://img.esa.io/uploads/production/attachments/2294/2016/03/15/6362/8655257b-fecc-4ce2-bb90-dd6f6868fe8d.png">

# param1が消えてない！
これは、一回セットすると、アプリ落とすまで消えないんだね。
なるほど。

じゃあ

## API呼び出す前に必ず初期化だ
``` postSample.go
// とにかく最初に呼び出す
func Initialize() {
	params = map[string]string{}
}
```
## Androidも修正
``` MainActivity.java
switch (v.getId()) {
    // リクエスト１ボタンが押された
    case R.id.request1_btn:
        PostSample.Initialize();   // ここ追加
        PostSample.SetParams("param1", param1.getText().toString());
        response.setText(PostSample.HttpPost("/apisample"));
         break;
    // リクエスト２ボタンが押された
    case R.id.request2_btn:
        PostSample.Initialize();  // ここ追加
        PostSample.SetParams("param2", param2.getText().toString());
        response.setText(PostSample.HttpPost("/apisample"));
         break;
}
```

# できた
<img width="300" alt="Screenshot_2016-03-15-19-23-35.png (113.6 kB)" src="https://img.esa.io/uploads/production/attachments/2294/2016/03/15/6362/73be7563-eb3d-409b-be84-0e399a87ab0b.png">

# iOSに追加

# 担当者の方、頑張って下さい！
frameworkってファイルができるみたいです。

# まとめ
超簡単に導入できた！

# コードや資料は<br/>Githubに置いときましたー
Github : https://github.com/kmizutani0222/GoMobileSample


#  参考
- [C++によるiOSとAndroidでのクロスプラットフォーム開発：Dropboxの教訓](http://www.infoq.com/jp/news/2014/06/dropbox-cpp-crossplatform-mobile)
- [gomobileでiOS用のライブラリをビルドするまで](http://qiita.com/naoty_k/items/e2bec591737218da1819)
- [gvmでgoをバージョン指定で簡単インストール](http://qiita.com/isaoshimizu/items/1a5d51aed98a57a9bcd4)
- [gomobileでiOS用のライブラリをビルドするまで](http://qiita.com/naoty_k/items/e2bec591737218da1819)
