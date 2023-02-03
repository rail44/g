## Requirements

- Go
  - 1.19
- PostgreSQL v15
  - (Docker Cli)
    - ローカル環境をDockerで立ち上げる場合
- [https://github.com/deepmap/oapi-codegen](oapi-codegen)
- [https://github.com/kyleconroy/sqlc](sqlc)
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
Golang製のAPIサーバーのクライアントがGolangであることはなかなかないと思われるので、OpenAPIのスキーマを中心に記述量を少なく実装したいと考えました。  
一方で、送金システムという要件からトランザクションを手で管理しながら一貫性を保つ必要がありそうなので、begin-commit内で上から順に愚直にコマンドを実行するスタイルが採りやすいことも大事そうだという印象を持ちました。  
このような条件から、SQLスキーマから型安全なアダプターを生成できるsqlcと、OpenAPIのスキーマから生成したやはり型安全なアダプターとの橋渡しに集中できるようなモデル層を、興味ことにサブパッケージとして切っていくようなディレクトリ構造を試しに構築した感じです。

また、コントローラー層とモデル層で発生しうる例外を分類してパッケージ内の独自エラーとして整備しました。  
今後もし別の例外が起きうるようなコードを追加する場合でも、どのような修正をすれば良いかの見通しがある程度効くかなと思います。
