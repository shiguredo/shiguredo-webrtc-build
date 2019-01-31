# 変更履歴

- UPDATE
    - 下位互換がある変更
- ADD
    - 下位互換がある追加
- CHANGE
    - 下位互換のない変更
- FIX
    - バグ修正

# 72.26.0

M72 コミットポジション 26

# 71.16.0

M71 コミットポジション 16

# 70.17.0

M70 コミットポジション 17

# 70.14.0

M70 コミットポジション 14

# 68.10.1

- [ADD] Docker での AAR ビルドに config.json フィールドいくつかを Java 定数フィールドとして
  持つクラスを追加した。

# 68.10.x

M68 コミットポジション 10

# 67.28.x

M67 コミットポジション 28

- [ADD] AAR ビルド用の Dockerfile, make ターゲットを追加した
- [ADD] M68 お試しビルドをしたので、68.6.0 でビルド、AAR を
  sora-webrtc-android にアップロードした
- [FIX] ios build から build_type が消えたため対応した
  - ios 用 patch の微修正
  - webrtc-build から build type 関連(static)を削除
- [FIX] ios archive コマンドのパスを修正した
- [ADD] build-clean コマンドを追加した
  - fetch, build でエラーになった場合に build-clean すると fetch の状態でビルドが可能
- [ADD] build_info.json を gitignore に追加した
- [FIX] ios ビルド後に dSYM diretory が存在しないため、archive からコピーを消した

# 66.8.x

M66 コミットポジション 8

## 66.8.2

- [CHANGE] 設定ファイルの内容を全体的に変更した

- [UPDATE] Google Paly services のライセンスに y で答えるよう変更した

## 66.8.1

- [CHANGE] Android: gclient の設定ファイルを 66.8 に対応した

## 66.8.0

- [CHANGE] コミットポジション 8 に対応した

# 66.6.x

M66 コミットポジション 6

# 64.6.x

M64

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
