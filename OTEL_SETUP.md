# OpenTelemetry Collector Setup

このドキュメントでは、Docker ComposeでOpenTelemetry Collectorを実行し、tracerからgRPCでOTLPエンドポイントにトレースデータを送信する方法を説明します。

## セットアップ

### 1. OpenTelemetry Collectorの起動

Docker Composeを使用してOpenTelemetry Collectorを起動します：

```bash
docker-compose up -d otel-collector
```

これにより、以下のポートでOTel Collectorが起動します：
- `4317`: OTLP gRPCレシーバー
- `4318`: OTLP HTTPレシーバー

### 2. 環境変数の設定

アプリケーションでOTel Collectorに接続するために、以下の環境変数を設定できます：

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
export OTEL_SERVICE_NAME=user-service
```

環境変数を設定しない場合、デフォルト値が使用されます：
- `OTEL_EXPORTER_OTLP_ENDPOINT`: `localhost:4317`
- `OTEL_SERVICE_NAME`: `user-service`

### 3. アプリケーションのビルドと実行

```bash
cd apps/user
go build -o user-service
./user-service
```

## 動作確認

### 1. ヘルスチェックエンドポイント

```bash
curl http://localhost:8080/health
```

レスポンス：
```json
{"status":"ok"}
```

### 2. サンプルトレースの生成

```bash
curl http://localhost:8080/example
```

レスポンス：
```json
{"message":"example operation completed"}
```

このエンドポイントは、ネストされたスパン（子スパン）を含むトレースを生成します：
- `example_operation` (親スパン)
  - `process_data` (子スパン)
  - `validate_data` (子スパン)

### 3. トレースの確認

OpenTelemetry Collectorのログを確認して、トレースが正常に送信されているかを確認できます：

```bash
docker-compose logs otel-collector
```

## アーキテクチャ

```
[Application] --gRPC OTLP--> [OTel Collector] --> [Debug Exporter (stdout)]
     :4317
```

現在の設定では、OTel Collectorは受信したトレースをdebugエクスポーターを使用して標準出力に出力します。

## OTel Collector設定

OTel Collectorの設定は `otel-collector-config.yaml` に定義されています：

- **Receivers**: OTLP gRPC (port 4317) および HTTP (port 4318)
- **Processors**: バッチプロセッサー（1秒のタイムアウト、1024のバッチサイズ）
- **Exporters**: デバッグエクスポーター（詳細ログ出力）

必要に応じて、設定ファイルを編集して他のエクスポーター（Jaeger、Zipkin、Prometheusなど）を追加できます。

## トラブルシューティング

### トレースが送信されない場合

1. OTel Collectorが起動していることを確認：
   ```bash
   docker-compose ps otel-collector
   ```

2. OTel Collectorのログを確認：
   ```bash
   docker-compose logs otel-collector
   ```

3. ネットワーク接続を確認：
   ```bash
   nc -zv localhost 4317
   ```

4. アプリケーションのログを確認して、エラーがないかチェック

### 接続エラーが発生する場合

- `OTEL_EXPORTER_OTLP_ENDPOINT` 環境変数が正しく設定されているか確認
- ファイアウォールやネットワーク設定でポート4317がブロックされていないか確認
