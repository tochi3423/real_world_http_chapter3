package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	// getRequest()
	// getRequestWithQueryParam()
	// headRequest()
	// postFormData()
	// postFile()
	// postText()
	// postMultipartData()
	// getRequestWithCookies()
	// getThroughProxy()
	getFileThroughProtocol()
	// deleteRequest()
}
func getRequest() {
	// この部分でhttp.Get関数を使用してHTTPリクエストを送っています。
	// これは指定されたURLへGETリクエストを送り、そのレスポンスを*http.Responseのrespに格納します。
	// *http.Response型の値は、HTTPレスポンスのステータスコード、ヘッダー、ボディ等の情報を含みます。
	// Goでは、関数に変数を渡した際に値がコピーされます。しかし、大きなデータ構造体をコピーするとパフォーマンス上不利な場合があるため、
	// その場合はポインタを利用してデータ構造体のメモリ上の位置（アドレス）だけを渡すことが一般的です。
	// エラーが発生した場合はerrに格納します。errがnil（エラーが発生しない）でない場合、panic関数を使ってプログラムを中断します。
	// Docker内部での通信では、serverではなく、サービス名やコンテナ名を使用します。
	resp, err := http.Get("http://server:18888")
	if err != nil {
		panic(err)
	}
	// deferを使って、main関数が終了する前にresp.Body.Close()が実行されるようにしています。
	// これはHTTPレスポンスのbodyを閉じるためのもので、bodyを閉じないとネットワークリソースが無駄に消費される可能性があります。
	defer resp.Body.Close()
	// ここでioutil.ReadAll関数を使ってレスポンスのbodyを読み込みます。
	// 全てのデータを読み取った後、それをbyteスライスとしてbodyに格納します。
	// エラーが発生した場合はerrに格納します。ここでも、エラーの発生をチェックし、エラーがあった場合はpanicでプログラムを中断します。
	body, err := ioutil.ReadAll(resp.Body)
	// body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// log.Println(body)
	// 以下のようなbyteスライスが出力されます。
	// [60 104 116 109 108 62 60 98 111 100 121 62 104 101 108 108 111 60 47 98 111 100 121 62 60 47 104 116 109 108 62 194 165 110]
	// 最後に、log.Println関数を使ってレスポンスbody（byteスライス）をstringに変換し、それをコンソールに出力します。そして、main関数は終わります。
	log.Println(string(body))
	log.Println("Status", resp.Status)
	log.Println("Status", resp.StatusCode)
	log.Println("Status", resp.Header)
	log.Println("Status", resp.Header.Get("Content-Length"))
}
func getRequestWithQueryParam() {
	values := url.Values{
		"query": {"hello world"},
	}
	// http.Get関数を用いて"http://server:18888"にリクエストを送信しています。
	// その際、URLの末尾にクエリ文字列（この場合は "?"に続く "query=hello+world"）を追加しています。
	// ここでクエリ文字列を作成するために values.Encode() メソッドを使用しています。
	resp, _ := http.Get("http://server:18888" + "?" + values.Encode())
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(body))
}
func headRequest() {
	resp, err := http.Head("http://server:18888")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// HTTPのステータスコード（例：200 OKなど）とヘッダー情報（例：Content-Type, Content-Lengthなど）を出力します。
	log.Println("Status", resp.Status)
	log.Println("Status", resp.Header)
}
func postFormData() {
	values := url.Values{
		"test": {"value"},
	}
	// http.PostForm関数を利用して作成した values を含む POST リクエストを "http://server:18888" へ送信します。
	// この関数は application/x-www-form-urlencoded 形式でデータをエンコードし、Content-Type ヘッダーを設定します。
	resp, err := http.PostForm("http://server:18888", values)
	if err != nil {
		// 送信失敗
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
func postFile() {
	// os.Open 関数を使用して、 Dockerfile というファイルを開いています。
	// この関数はファイルを読み取り専用として開きます。
	// エラーが発生した場合（例えば、ファイルが存在しない）、 err にエラー情報が格納され、そのエラーは panic 関数を介して処理されます。
	file, err := os.Open("./Dockerfile")
	if err != nil {
		panic(err)
	}
	// http.Post 関数を使い、先ほど開いたDockerfileの内容を http://server:18888 へPOSTリクエスト送信します。
	// ここでは、MIME typeとして text/plain を指定しています。再度エラーチェックを行い、何か問題があればエラー情報が格納されます。
	resp, err := http.Post("http://server:18888", "text/plain", file)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
func postText() {
	// strings.NewReader関数を使用して、"テキスト"という文字列から新しいstrings.Readerを作成します。
	// これは、文字列がソースとなるIOリーダーを生成します。
	// このstrings.Readerは、後続のHTTP POSTリクエストでデータソースとして使用されます。
	reader := strings.NewReader("テキスト")
	// http.Post関数を使って、先ほど作成したstrings.Reader（reader）から取得したデータをhttp://server:18888へPOSTリクエスト送信します。
	// ここでは、MIMEタイプとしてtext/plainを指定しています。エラーチェックを行い、何か問題があればエラー情報が格納され、エラーはpanic関数によって処理されます。
	resp, err := http.Post("http://server:18888", "text/plain", reader)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
func postMultipartData() {
	//  bytes.Bufferを作成し、これを用いてmultipart.NewWriterを作ります。
	// これにより、マルチパートエンコーディングのHTTPリクエストボディを作成することが可能になります。
	// そして、そのwriterを用いて"name"というキーに"Michael Jackson"という値を書き込みます。
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "Michael Jackson")
	// 同様のwriterを使用して"thumbnail"というキーに"photo.jpg"という値を書き込みます。
	// ここでは、"photo.jpg"というファイルを開き、その内容を新たに作成したフィールドであるfilerWriterにコピーします。
	filerWriter, err := writer.CreateFormFile("thumbnail", "photo.jpg")
	if err != nil {
		panic(err)
	}
	readFile, err := os.Open("photo.jpg")
	io.Copy(filerWriter, readFile)
	writer.Close()
	// 作成したマルチパート形式のHTTPリクエストボディを含むHTTP POSTリクエストをhttp://server:18888に対して送信します。
	// ここでは、writer.FormDataContentType()を用いてContent-Typeヘッダーを設定しています。
	resp, err := http.Post("http://server:18888", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
func getRequestWithCookies() {
	// cookiejar.New 関数を使用して新たなCookieジャー(jar)を作成します。
	// CookieジャーはHTTPクライアントによって送信されるHTTPリクエストに関連付けられたCookieを保存します。
	// これによりHTTPクライアントはリクエスト間でCookieを維持し、セッション管理やページ間の情報の永続化などに利用することができます。
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	// 新しく作成したCookieジャーを持つHTTPクライアント(client)を作成します。
	client := http.Client{
		Jar: jar,
	}
	// このHTTPクライアントを使ってhttp://server:18888/cookieへ2回GETリクエストを送信します。
	// 2回目以降のアクセスでクッキーをサーバーに送信する仕組みなので、2回アクセスします。
	// 各レスポンスについて、httputil.DumpResponse関数を使用してHTTPレスポンスをダンプし、その結果をログに出力します。
	// レスポンスダンプは、HTTPレスポンスヘッダーとHTTPレスポンスボディの両方を含みます。
	for i := 0; i < 2; i++ {
		resp, err := client.Get("http://server:18888/cookie")
		if err != nil {
			panic(err)
		}
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			panic(err)
		}
		log.Println(string(dump))
	}
}
func getThroughProxy() {
	// url.Parse 関数を用いて文字列からURLオブジェクトを作成します。
	// ここで作成されたURLオブジェクト(proxyUrl)は、後続のHTTPクライアントのプロキシ設定で使用されます。
	proxyUrl, err := url.Parse("http://server:18888")
	if err != nil {
		panic(err)
	}
	// HTTPクライアントを作成します。
	// ここでは、HTTPトランスポートのProxy属性にhttp.ProxyURL(proxyUrl)を指定して、
	// 全てのHTTPリクエストが指定したプロキシURL(http://server:18888)を経由するように設定しています。
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	// 作成したHTTPクライアントを用いてhttp://github.xn--comhttp-7g4f GETリクエストを送信します。
	resp, err := client.Get("http://github.com")
	if err != nil {
		panic(err)
	}
	// レスポンスを受け取ったら、httputil.DumpResponse 関数を使ってHTTPレスポンスをダンプし、その結果をログに出力します。
	// レスポンスダンプは、HTTPレスポンスヘッダーとHTTPレスポンスボディの両方を含みます。
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
func getFileThroughProtocol() {
	// 新しい http.Transport を作成し、RegisterProtocol メソッドを呼び出して file プロトコルに対するハンドラを設定しています。
	// このハンドラとして http.NewFileTransport(http.Dir(".")) を指定しています。
	// これは、ローカルのファイルシステム上のファイルを参照するためのハンドラで、http.Dir(".") は現在のディレクトリを表しています。
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	// 新しく作成したトランスポートを使用するHTTPクライアントを作成します。
	client := http.Client{
		Transport: transport,
	}
	// 作成したHTTPクライアントを用いて file://./Dockerfile へHTTP GETリクエストを送信します。
	// ここでの file://./Dockerfile は、現在のディレクトリ下にある Dockerfile ファイルを指しています。
	resp, err := client.Get("file://./Dockerfile")
	if err != nil {
		panic(err)
	}
	// レスポンスを受け取ったら、httputil.DumpResponse 関数を使ってHTTPレスポンスをダンプし、その結果をログに出力します。
	// レスポンスダンプは、HTTPレスポンスヘッダーとHTTPレスポンスボディの両方を含みます。
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}

// DELETEリクエストは該当URLのリソースの削除を意味します。
// 例えば、http://server:18888/users/1234 に対するDELETEリクエストは、
// IDが1234のユーザを削除することを示すといった具体的な意味を持つことになります。
func deleteRequest() {
	// 新しいHTTPクライアントを作成します。このHTTPクライアントは、後続のHTTPリクエストを処理するためのものです。
	client := &http.Client{}
	// http.NewRequest関数を使って新たなHTTPリクエストを作成します。
	// ここでは、DELETEメソッドを使用してhttp://server:18888を対象とするリクエストを作成しています。
	request, err := http.NewRequest("DELETE", "http://server:18888", nil)
	if err != nil {
		panic(err)
	}
	// 作成したHTTPクライアントを利用して、先程作成したHTTPリクエストを送信します。その結果として得られたHTTPレスポンスはrespに格納されます。
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	// httputil.DumpResponse関数を使ってHTTPレスポンスをダンプし、その結果をログに出力します。
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
