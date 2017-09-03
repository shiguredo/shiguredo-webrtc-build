# WebRTC ライブラリ用ビルドツール

iOS と Android 向けの WebRTC ライブラリをビルドします。WebRTC ライブラリのビルドは複雑でわかりにくいです。
また WebRTC ライブラリのバージョンが上がるごとにビルド方法が変わることも多く、追従するのは現実的ではありませんでした。

このツールはもともと[株式会社時雨堂](https://shiguredo.jp)の社内ツールでしたが、
少しでも WebRTC ライブラリに追従する負荷を削減できればと思い公開することにしました。

## 注意

このビルドツールの主な目的は[株式会社時雨堂](https://shiguredo.jp)の製品である [WebRTC SFU Sora](https://sora.shiguredo.jp) の SDK 利用のために開発されています。
そのため、目的を越えた汎用的なビルドツールにする予定はありません。リソースの都合上、バグ報告以外の PR に対応できるとは限りませんのでご了承ください。

## 方針

Chrome の安定版のバージョンに合わせて追従していきます。 master ブランチやリリースブランチの HEAD への追従は行いません。

## ダウンロード

実行可能なビルドツールは [Releases](https://github.com/shiguredo/sora-webrtc-build/releases) からダウンロードできます (ビルド済みの WebRTC ライブラリではありません) 。
リポジトリを clone する必要はありません。

## ビルドのシステム条件

Android のビルドは Linux のみサポートされているため、一台の Mac で iOS/Android の両方のビルドはできません。

- iOS
  - Mac OS X 10.12.5+
  - Xcode 8.3.3+
  - Python 2.7
- Android
  - Ubuntu Linux 16.04 64bit
  - Python 2.7

## バージョン表記について

バージョンは「メジャー.マイナー.メンテナンス」で表します。

- メジャーバージョン: WebRTC ライブラリのリリースブランチ番号を表します。
- マイナーバージョン: WebRTC ライブラリのコミットポジション番号を表します。 コミットポジション番号とは WebRTC ライブラリでリビジョンごとに割り振られる数値です。 リリースブランチ作成時のコミットポジションは 1 であるため、マイナーバージョンは常に 1 から始まります。
- メンテナンスバージョン: WebRTC ライブラリ以外の変更を加えたときに上がります。

[6d56d2e](https://chromium.googlesource.com/external/webrtc/+/6d56d2ebf0e40bc73dee99093ac7c223ddc7b6e5) を例にすると、リリースブランチ番号が 59 、コミットポジション番号が 4 なので、バージョンは 59.4.x になります。

## 仕様

- 対応する WebRTC のバージョン: M59 ([a100a39](https://chromium.googlesource.com/external/webrtc/+/a100a39fd25df18f51bf0144d1347fea2462a279))
- 対応するアーキテクチャ: arm64, armv7 (iOS), armarmeabi-v7a/arm64-v8a (Android)
- VP9 有効
- (iOS) Bitcode 対応

## 使い方

``webrtc-build`` コマンドを使います。
``all`` サブコマンドを指定すると一連の処理 (ソースの取得からビルドまで) を一括で行います。

```
$ ./webrtc-build all
```

### 主なサブコマンド

次の順序でサブコマンドを実行するとライブラリをビルドできます。
実行するプラットフォームによってビルド対象が異なります。
Mac OS X では iOS 向け、 Linux では Android 向けのライブラリがビルドされます。

1. ``./webrtc-build setup``

2. ``./webrtc-build fetch``

3. ``./webrtc-build build``

4. ``./webrtc-build dist`` (任意)

``./webrtc-build setup`` は WebRTC のビルドで使われるツール [depot_tools](https://www.chromium.org/developers/how-tos/depottools) をローカルに取得します。
マシンに depot_tools がインストールされていてもこちらが使われますので、必ず実行してください。

``./webrtc-build fetch`` は WebRTC のソースコードをダウンロードし、ビルドに必要となるファイルを生成します。途中で中断してしまっても、再度実行すれば再開できます。

``./webrtc-build build`` はライブラリをビルドします。ビルドの成果物は ``webrtc/build`` ディレクトリ以下にあります。

``./webrtc-build dist`` はビルド成果物の配布用アーカイブを生成します。このコマンドの実行は必須ではありません。

### iOS ライブラリのビルド用のコマンド

Mac OS X では次のコマンドで iOS 向けのライブラリを個別にビルドできます。

- ``build-framework-debug``: デバッグ設定のフレームワークをビルドします。
- ``build-framework-release``: リリース設定のフレームワークをビルドします。
- ``build-static-debug``: デバッグ設定の静的ライブラリをビルドします。
- ``build-static-release``: リリース設定の静的ライブラリをビルドします。

### Xcode について

iOS 向けライブラリのビルドは、 Xcode の代わりに WebRTC のソースコードに含まれる clang を使うようにしてあります。
最適化の効率は Xcode と比べて落ちるかもしれませんが、 Xcode のバージョンによるビルドエラーを回避できます。

### Android ライブラリのビルド用のコマンド

Linux では次のコマンドで Android 向けのライブラリを個別にビルドできます。

- ``build-debug``: デバッグ設定の AAR ファイルをビルドします。
- ``build-release``: リリース設定の AAR ファイルをビルドします。

### Android ライブラリのビルド手順

Android 向けのビルドはいくつか注意点があります。

- 実行時のシェルは Bash を推奨します。 Zsh だとビルドエラーになる場合がありました。

- 初回の ``./webrtc-build fetch`` は途中で失敗します。これは Google Play ライブラリのライセンスの同意を求める処理によるもので、 ``webrtc`` ディレクトリ下で ``./depot_tools/gclient sync`` してライセンスに同意してください。ライセンスに同意後にダウンロードが再開されます。ダウンロードに成功したら、再度 ``./webrtc-build fetch`` を実行してください。

以上を踏まえて、ビルドは次の手順で行います。

1. シェルを Bash に変更する

2. ``./webrtc-build setup`` を実行する

3. ``./webrtc-build fetch`` を実行する。初回はライセンスの同意を求める処理で止まります

4. ``webrtc`` ディレクトリ下で ``gclient sync`` を実行し、途中で表示されるライセンスに同意する。初回以降は不要です

   ```
   $ cd webrtc
   $ ./depot_tools/gclient sync
   ```

5. 再度 ``./webrtc-build fetch`` を実行する

6. ``./webrtc-build build`` を実行する

### その他のコマンド

その他のコマンドを次に示します。 iOS/Android で共通です。

- ``./webrtc-build update``: ``clean``, ``reset``, ``setup``, ``fetch`` を順に実行します。 WebRTC のリビジョンの変更時に行うと便利です。

- ``./webrtc-build clean``: ビルド過程で生成されたファイルをすべて削除します。

- ``./webrtc-build reset``: WebRTC のソースコードに加えられた変更をすべて破棄し、リビジョンの状態に戻します。

- ``./webrtc-build help``: ヘルプメッセージを表示します。

- ``./webrtc-build version``: ビルドツールのバージョンを表示します。

## ビルドの設定

``config.json`` でビルドの設定が可能です。
リリースブランチやリビジョンを変更したときは、 ``update`` または ``all`` で既存のソースコードへの変更を破棄してから更新しておくとビルドの失敗を防げます。

- ``webrtc_branch``: リリースブランチ番号。
- ``webrtc_commit``: コミットポジション番号。ソースコードのダウンロードに影響しません。
- ``webrtc_revision``: リビジョン番号。リリースブランチの取得後、指定したリビジョンをチェックアウトします。
- ``python``: Python のパス。 WebRTC のソースコードに含まれるビルドスクリプトで使われます。
- ``ios_arch``: iOS ライブラリでサポートするアーキテクチャ。 ['arm64', 'arm', 'x64', 'x86'] から複数選択可能です。
- ``android_arch``: Android ライブラリでサポートするアーキテクチャ。 ['armeabi-v7a', 'arm64-v8a', 'x86', 'x86_64'] から複数選択可能です。

## ビルド情報 (iOS)

ビルド時の情報は `build_info.json` に保存されます。フレームワークには ``build_info.json`` が含まれています。

- `webrtc_version` (string): WebRTC のリリースブランチ番号
- `webrtc_revision` (string): WebRTC のリビジョン番号

## トラブルシューティング

ビルドエラーの大半は depot_tools が原因である可能性が高いです。
エラーが出たら、まず次の処理を試してみてください。

- **(特に Android) シェルを Bash に変更する**

- ``./webrtc-build fetch`` を実行する

- ``webrtc`` ディレクトリで ``gclient runhooks`` を実行する

### depot_tools のコマンド (gclient など) を使いたい

``./webrtc-build setup`` でダウンロードした depot_tools は ``webrtc/depot_tools`` にダウンロードされます。コマンド検索パスにディレクトリを追加します。

例: ``webrtc/src`` ディレクトリで depot_tools のコマンドを使う場合

```
export PATH=../depot_tools:$PATH
```

### ``./webrtc-build fetch``: 実行が途中で止まってしまった

再度 ``./webrtc-build fetch`` を実行してください。
``./webrtc-build fetch`` は何度実行しても問題ありません。

### ``./webrtc-build fetch``: ``stderr:error: Your local changes to the following files would be overwritten by checkout``

パッチの適用後に ``./webrtc-build fetch`` を実行すると、リポジトリに変更があるために指定のリビジョンをチェックアウトできずにエラーになります。
``./webrtc-build reset`` を実行して、リポジトリの変更を戻してから再度 ``./webrtc-build fetch`` を実行してください。

### ``./webrtc-build build``: ``stderr:.gclient file in parent directory XXX might not be the file you want to use``

このエラーが出たら冒頭の方法を試してみてください。

```
stderr:.gclient file in parent directory /home/shiguredo/sora-webrtc-build/webrtc might not be the file you want to use
stderr:gn.py: Could not find gn executable at: /home/shiguredo/sora-webrtc-build/webrtc/src/buildtools/linux64/gn
```

### ``./webrtc-build build``: ``Error: Command 'XXX' returned non-zero exit status 1 in XXX``

このエラーでも冒頭の方法を試してみてください。

```
________ running '/usr/bin/python src/third_party/binutils/download.py' in '/home/shiguredo/sora-webrtc-build/webrtc'
Downloading /home/shiguredo/sora-webrtc-build/webrtc/src/third_party/binutils/Linux_x64/binutils.tar.bz2
Traceback (most recent call last):
  File "src/third_party/binutils/download.py", line 130, in <module>
    sys.exit(main(sys.argv[1:]))
  File "src/third_party/binutils/download.py", line 117, in main
    return FetchAndExtract(arch)
  File "src/third_party/binutils/download.py", line 82, in FetchAndExtract
    '-s', sha1file])
  File "/usr/lib/python2.7/subprocess.py", line 536, in check_call
    retcode = call(*popenargs, **kwargs)
  File "/usr/lib/python2.7/subprocess.py", line 523, in call
    return Popen(*popenargs, **kwargs).wait()
  File "/usr/lib/python2.7/subprocess.py", line 711, in __init__
    errread, errwrite)
  File "/usr/lib/python2.7/subprocess.py", line 1343, in _execute_child
    raise child_exception
OSError: [Errno 2] No such file or directory
Error: Command '/usr/bin/python src/third_party/binutils/download.py' returned non-zero exit status 1 in /home/shiguredo/sora-webrtc-build/webrtc
```

## 新しいバージョンへの対応 (開発者向け)

- コミットポジションごとに ``develop`` ブランチから新しいブランチを派生する。ブランチ名は "リリースブランチ.コミットポジション.x" とする。

- ``config.json`` に WebRTC のバージョンを記述する。

  ```
  "webrtc_branch": "60",
  "webrtc_commit": "9",
  "webrtc_revision": "9710de31ef15f42c86ccb0d69bd245da940b16fa",
  ```

- ``webrtc-build.go`` 内で定義されているバージョンを変更する。

  ```
  // 新しいバージョンに変更する
  var version = "60.9.1"
  ```

- CHANGES に追記する。

- タグを打つ。タグ名は "リリースブランチ.コミットポジション.メンテナンス" とする。

- ``develop`` ブランチにマージする。新しいブランチは削除せず、引き続き使用する。

- ``make`` を実行して ``webrtc-build.go`` をビルドする。

- ``make dist`` を実行する。
  実行したプラットフォーム向けの ``sora-webrtc-build-*.tar.gz`` が生成される。

- GitHub のリリースノートに ``sora-webrtc-build-*.tar.gz`` を添付する。
