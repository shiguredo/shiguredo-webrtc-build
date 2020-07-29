# 時雨堂のパッチ (iOS) について

## 対象のビルド設定

- ios-\*-shiguredo
- ios-\*-shiguredo-develop

## 目的

- 接続時、マイクの使用不使用に関わらずパーミッションが要求されてしまう事象を修正する。

## パッチ内容

- 接続時のマイクのパーミッション要求を抑制する。

- マイクの初期化を明示的に行う API を追加する。
  パッチ適用後はマイクは自動的に初期化されない。

- デフォルトで設定される ``AVAudioSession`` のカテゴリを ``AVAudioSessionCategoryPlayAndRecord`` から ``AVAudioSessionCategoryAmbient`` に変更する。

## パッチ適用後の使い方

- マイクを使う場合は ``RTCAudioSession.initializeInput(completionHandler:)`` を実行してマイクを初期化する。このメソッドはマイクの初期化の準備が整うまで非同期で待ち、初期化可能な状態になったらパーミッションを要求する。ユーザーがマイクの使用を許可したら初期化を行い、ハンドラを実行する。

- ``RTCAudioSession.initializeInput(completionHandler:)`` の前後で ``RTCAudioSession.lockForConfiguration`` と ``RTCAudiotSession.unlockForConfiguration`` を実行する必要はない。内部で必要に応じてロックしているので、 ``initializeInput(completionHandler:)`` の前後でロックすると、おそらく何かしらの競合を起こしてマイクが初期化されない。

- マイクを使わない場合は ``Info.plist`` にマイクの用途を記述する必要はない。
