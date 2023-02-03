## Requirements

- Go
  - 1.19
- PostgreSQL v15
  - (Docker Cli)
    - ローカル環境をDockerで立ち上げる場合
- [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [sqlc](https://github.com/kyleconroy/sqlc)
  - コードジェネレーター類
  - 生成物もコミットしてあるので、動かすだけなら不要です
  - `go install` にてcliツールが動く状態になっていれば大丈夫です
- GNU make

## Run

動作のための設定類は[Makefile](./Makefile)にマクロで記述されています。  
PostgreSQLをDockerで動かす場合は、リポジトリをクローンして特に変更なしに動作させることができるようになっているはずです。

```bash
$ make db/run

...
```

```bash
# in other terminal

$ make db/schema
psql -f sqlc/schema.sql postgresql://postgres:password@localhost:5432/g
DROP SCHEMA
CREATE SCHEMA
CREATE TABLE
CREATE TABLE
CREATE TABLE
CREATE TABLE
CREATE TABLE
CREATE INDEX
CREATE INDEX
CREATE TABLE

$ make g/run
2023/02/03 16:56:58 Listening on :3000

# openapiやsqlcの変更に対してコード生成しなおすmakeターゲット

$ make generate

```

## Architecture

```
g/
├─ main.go
├─ accounts/
│  ├─ controller.go
│  ├─ model.go
│  ├─ util.go
│  ├─ openapi.yml
│  ├─ openapi.gen.go     # Generated
├─ sqlc/
│  ├─ schema.sql
│  ├─ queries.sql
│  ├─ generated/         # Generated
│  ├─   ├─ ...
...
```


シンプルなModel-Viewベースの構成です。  
まず、送金システムという要件から、ReadしてCheckしてWriteのような処理のニーズがあり、トランザクショナルなデータベースを採用したいと考えました。
また、 Golang製のAPIサーバーのクライアントがGolangであることはなかなかないと思われるので、REST APIであればOpenAPIのスキーマを中心に開発をしていくと体験が良さそうだと考えました。

この2つの条件から、SQLのスキーマとOpenAPIのスキーマそれぞれから型安全なアダプターを生成できるように環境整備することで、独自通貨サービスの業務ロジックに集中できる基盤を作れるのではないかと考えました。
調査したところ、sqlcとoapi-codegenというツールが、他のライブラリへのロックインの少なさや、生成するコードの読みやすさ及び小ささから要件に合致していると考え、設計と実装を行いました。

設計については、コントローラー、モデルのような粒度よりも操作したいリソースでサブパッケージを切っていく方ような構造を採ってみました。  
名前空間を狭く近づけられたおかげで、それぞれのオブジェクト内で自然な命名とアクセスでメンバーを触れているのではないかなと思っています。  

細かいところとしては、起きうる例外のパターンに対してサブパッケージ内に独自errorを実装していたり、トランザクションを含む処理を楽に行うためのユーティリティの整備なども行っています。  
今後これらのコンポーネントにパラメーターやルーティングを追加しようと考えたときに、どこにどのような修正をすれば良いかの見通しがある程度効く構成を目標に実装を行いました。

## Demonstrations

#### Register

```bash
$ curl --data '{"name": "hoge"}' http://localhost:3000/accounts
{"accountId":1}
```

#### Register

```bash
$ curl http://localhost:3000/accounts/1/balance
{"balance":0}
```

#### Mint

```bash
$ curl http://localhost:3000/accounts/1/balance
{"balance":100}
$ curl --data '{"amount": 100}' http://localhost:3000/accounts/1/mint
{"transactionId":1}
$ curl http://localhost:3000/accounts/1/balance
{"balance":100}
```

#### Spend

```bash
$ curl http://localhost:3000/accounts/1/balance
{"balance":100}
$ curl --data '{"amount": 50}' http://localhost:3000/accounts/1/spend
{"transactionId":2}
$ curl http://localhost:3000/accounts/1/balance
{"balance":50}
```

#### Transfer

```bash
$ curl http://localhost:3000/accounts/1/balance
{"balance":50}
$ curl http://localhost:3000/accounts/2/balance
{"balance":0}

$ curl --data '{"amount": 20, "recipient": 2}' http://localhost:3000/accounts/1/transfer
{"transactionId":3}
$ curl http://localhost:3000/accounts/1/balance
{"balance":30}
$ curl http://localhost:3000/accounts/2/balance
{"balance":20}
```


#### Transactions

```bash
$ curl http://localhost:3000/accounts/1/transactions
[{"account":1,"amount":100,"id":1,"inserted_at":"2023-02-03T09:19:28.369151Z","type":"mint"},{"account":
1,"amount":50,"id":2,"inserted_at":"2023-02-03T09:19:38.552711Z","type":"spend"},{"account":1,"amount":2
0,"id":3,"inserted_at":"2023-02-03T09:20:12.201018Z","recipient":2,"type":"transfer"}]

$ curl http://localhost:3000/accounts/2/transactions
[{"account":1,"amount":20,"id":3,"inserted_at":"2023-02-03T09:20:12.201018Z","recipient":2,"type":"trans
```
