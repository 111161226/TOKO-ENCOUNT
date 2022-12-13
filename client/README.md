# System Design Todo List Client
## 起動方法
ルートディレクトリで`docker-compose up`すればクライアントサイドもビルドされ、[localhost:3050](http://localhost:3050)で配信される。

## クライアントサイドのみでデバックする場合
### 初回のみ必要
```bash
$ npm i
$ npm run gen-api
```

### 毎回必要
まずターミナルで以下を実行し、モックサーバーを起動する
```bash
$ npm run mock
```

その後別のターミナルを起動し、以下のコマンドを実行後、[localhost:3000](http://localhost:3000)にアクセス
```bash
$ npm run dev
```

## ディレクトリ構成
```txt
client
├── public: ロゴ置き場
├── scripts: コードを自動生成する為のコード
├── src: 実装部分のコード置き場
│   ├── layout: レイアウトの実装
│   ├── lib/api: apiを叩くためのコード。generatedに自動生成コードが入る
│   ├── pages: 各ページの実装
│   ├── router: ルーティングの設定
│   ├── store: storeの設定
│   ├── types: 型の定義
│   └── util: その他必要なロジックの実装
└── README.md: このファイル
```
