# 変更履歴

- UPDATE
    - 下位互換がある変更
- ADD
    - 下位互換がある追加
- CHANGE
    - 下位互換のない変更
- FIX
    - バグ修正

# 63.13.x

M63

## 63.13.0

- [CHANGE] M63 のビルドに対応した

- [CHANGE] webrtc-build.go のバージョンを WebRTC のバージョンと分けた

- [CHANGE] 次のサブコマンドを削除した

  - all

  - setup

  - update

  - reset

  - debug

  - release

  - framework-debug (iOS)

  - framework-release (iOS)

  - static-debug (iOS)

  - static-release (iOS)

- [CHANGE] "dist" サブコマンドを "archive" に変更した

- [CHANGE] "fetch" サブコマンドで depot_tools をダウンロードするように変更した

- [CHANGE] "build" サブコマンドは config.json で指定された設定のみをビルドするように変更した

- [CHANGE] "-no-patch" コマンドラインオプションを削除した

- [CHANGE] config.json に次の項目を追加した

  - ios_targets

  - ios_bitcode

  - build_config

  - vp9

  - apply_patch

  - patches

- [CHANGE] パッチファイルを config.json で指定できるようにした

# 62.12.x

M62

## 60.12.0

- [CHANGE] コミットポジション 12 のビルドに対応した

- [CHANGE] コマンドラインオプション "-no-patch" を追加した

- [CHANGE] gclient sync の実行をソースコードのチェックアウト後に変更した

- [CHANGE] iOS: アーキテクチャに armv7 を追加した

# 60.x

## 60.9.1

- [CHANGE] コミットポジション 9 のビルドに対応した

- [CHANGE] iOS: コンパイラに Xcode を使うようにした

- [CHANGE] iOS: Xcode 8.3.3 に対応した

- [CHANGE] iOS: Bitcode を有効にした

## 60.4.1

- [CHANGE] アーキテクチャに x86_64 を追加した

## 60.4.0

- [CHANGE] コミットポジション 4 のビルドに対応した

- [CHANGE] サブコマンド "all", "update" を追加した

## 60.1.0

- [CHANGE] M60 のビルドに対応した

- [CHANGE] config.json で指定できた webrtc-build のバージョンをソースコードに直接記述した

- [CHANGE] webrtc-build のサブコマンド名を短縮した

# 59.x

## 59.1.4

60.4.0 以降に 59.1.3 からブランチを切ったため、以下の変更は 60.x に影響しない。

- [CHANGE] アーキテクチャに x86_64 を追加した

## 59.1.3

- [CHANGE] Makefile を追加した

- [CHANGE] iOS: フレームワーク: WebRTC.h にインポートするヘッダーファイルを追加した

  - RTCCameraVideoCapturer.h

  - RTCMTLVideoView.h

  - RTCVideoCapturer.h

- [FIX] アーカイブの拡張子が .zip になっていたのを .tar.gz に修正した

## 59.1.2

- [CHANGE] パッチファイルをリリースバイナリに含める

## 59.1.1

- [CHANGE] iOS: フレームワーク: ヘッダーファイルに RTCCameraVideoCapturer.h を追加する

## 59.1.0

最初のリリース。
