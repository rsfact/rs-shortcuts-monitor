# RS Shortcuts Monitor Client

Windows用のキーボードショートカット監視アプリケーションです。

## ビルド方法

```bash
cd client
./build.sh
```

## 設定

`client/settings.ini`でWebhook URLとユーザーIDを設定してください。

```ini
[default]
url = https://your-webhook-url.com
user_uid = your-user-id
```

## 実行

`client/dist/rs-shortcuts-monitor-client.exe`を実行してください。
