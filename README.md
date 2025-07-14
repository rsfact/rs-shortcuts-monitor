# RS Shortcuts Monitor Client

Windows用のキーボードショートカット監視アプリケーションです。

## ビルド方法

```bash
cd client
./build.sh
```

バージョンを付与してZip配布

## 署名例

```bat
cd C:\Program Files (x86)\Windows Kits\10\bin\10.0.26100.0\x64

signtool.exe sign /fd SHA256 /f "D:\Desktop\Systems\certificate.pfx" /p "password" /t "http://timestamp.digicert.com" /d "RS-Shortcuts-Monitor" "D:\Desktop\Systems\Github\rsfact\rs-shortcuts-monitor\client\dist\rs-shortcuts-monitor-client.exe"
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
