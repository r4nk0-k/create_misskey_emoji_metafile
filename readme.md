# 概要
画像が入ったディレクトリにmisskeyインポート用のjsonファイルを作る

# 使い方
`./cfg/config.yaml`を編集
`go run ./pkg/create_metafile.go 画像のあるディレクトリ`
もしくは
`./bin/create_metafile(.exe) 画像のあるディレクトリ`

完了すると画像のあるディレクトリにjsonファイルが出来るので、それと画像をzipに固めてインポートしてください