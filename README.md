# mackerel-plugin-linux-usage

LinuxのCPU使用率、CPUあたりのロードアベレージ、プロセス数、TCP接続に関するメトリクスを取得するMackerel用プラグインです。Linuxの `/proc` を参照し、mackerel-agentが読み取れる形式で標準出力に出力します。

## 動作環境

- Linux（`/proc` がマウントされ、実行ユーザーが必要な情報を読み取れること）
- リリースバイナリの対応アーキテクチャ: amd64、arm64
- ソースからビルドする場合: Go 1.26.0以上
- 状態ファイルの保存先に対する読み書き権限

## インストール

### mkrを利用する場合

```sh
mkr plugin install monitoring-forge/mackerel-plugin-linux-usage
```

### リリースバイナリを利用する場合

[GitHub Releases](https://github.com/monitoring-forge/mackerel-plugin-linux-usage/releases)から、環境に合ったZIPファイルをダウンロードして展開してください。実行権限を付与し、任意のディレクトリに配置します。

## 使い方

### 手動で実行する

追加の引数は不要です。

```sh
/path/to/mackerel-plugin-linux-usage
```

出力は「メトリクス名・値・UNIXタイムスタンプ」のタブ区切りです。以下は出力の一部です。

```text
linux-usage.cpu.user	2.5	1750000000
linux-usage.cpu.idle	95	1750000000
linux-usage.loadavg.loadavg1	0.25	1750000000
linux-usage.process.all	120	1750000000
linux-usage.process.running	2	1750000000
linux-usage.tcp-opens.active	12	1750000000
```

| オプション | 説明 |
|---|---|
| `-h`, `--help` | ヘルプを表示します。 |
| `-v`, `--version` | バージョン、OS、アーキテクチャ、Goバージョン、コミット情報を表示します。 |

CPU使用率とTCP関連のメトリクスは前回実行との差分を使うため、状態ファイルがない初回実行では出力されません。通常は約1分間隔で実行してください。

### mackerel-agentに設定する

`mackerel-agent.conf` に次の設定を追加し、mackerel-agentを再起動してください。`command` は実際のインストール先に合わせて変更します。

```toml
[plugin.metrics.linux-usage]
command = "/opt/mackerel-agent/plugins/bin/mackerel-plugin-linux-usage"
```

グラフ定義のみを確認する場合は、次のように実行します。

```sh
MACKEREL_AGENT_PLUGIN_META=1 /path/to/mackerel-plugin-linux-usage
```

## メトリクスの意味

### CPU使用率: `linux-usage.cpu.*`

`/proc/stat` のCPU時間の前回値との差分から算出します。単位は `%` で、すべての論理CPUを合計した時間を100%として正規化します。CPU数によって上限が増えることはありません。`idle` を含む各項目の合計は約100%になります。

| メトリクス名 | 意味 |
|---|---|
| `linux-usage.cpu.user` | 通常のユーザープロセスの実行時間。`guest` 分を除きます。 |
| `linux-usage.cpu.nice` | nice値で優先度を調整したユーザープロセスの実行時間。`guest_nice` 分を除きます。 |
| `linux-usage.cpu.system` | カーネルでの実行時間。 |
| `linux-usage.cpu.idle` | CPUがアイドル状態だった時間。 |
| `linux-usage.cpu.iowait` | I/O完了待ちとして計上された時間。ディスク使用率そのものではありません。 |
| `linux-usage.cpu.irq` | ハードウェア割り込み処理の時間。 |
| `linux-usage.cpu.softirq` | ソフトウェア割り込み処理の時間。 |
| `linux-usage.cpu.steal` | 仮想CPUがハイパーバイザーにより実行を待たされた時間。 |
| `linux-usage.cpu.guest` | ゲストの仮想CPUの実行時間。 |
| `linux-usage.cpu.guest_nice` | nice値で優先度を調整したゲストの仮想CPUの実行時間。 |

### ロードアベレージ: `linux-usage.loadavg.*`

`/proc/loadavg` のロードアベレージを、`/proc/stat` から取得した論理CPU数で割った値です。単位は無次元で、パーセントではありません。実行可能なタスクに加えて、割り込み不能な待機状態（I/O待ちなど）のタスクも負荷に含まれます。

| メトリクス名 | 意味 |
|---|---|
| `linux-usage.loadavg.loadavg1` | 過去1分のロードアベレージ ÷ 論理CPU数。 |
| `linux-usage.loadavg.loadavg5` | 過去5分のロードアベレージ ÷ 論理CPU数。 |
| `linux-usage.loadavg.loadavg15` | 過去15分のロードアベレージ ÷ 論理CPU数。 |

例: 4論理CPUでロードアベレージが2の場合、出力は `0.5` です。`1.0` は論理CPU数と同じ数のタスクが負荷に含まれる状態を表しますが、CPU使用率100%を意味するものではありません。

### プロセス数: `linux-usage.process.*`

`/proc` にあるPIDごとの `stat` を読み取り、取得できたプロセスを数えます。単位はプロセス数です。

| メトリクス名 | 意味 |
|---|---|
| `linux-usage.process.all` | 状態を読み取れたプロセスの総数。 |
| `linux-usage.process.running` | 状態が `R`（実行中または実行待ち）のプロセス数。 |

スレッド単位の集計ではありません。読み取り中に終了したプロセスや、権限不足などで状態を取得できなかったプロセスは除外されます。

### TCP接続: `linux-usage.tcp-opens.*`

`/proc/self/net/snmp` の累積カウンターを使います。出力値は **1分あたりの増加数（件/分）** です。現在の接続数ではありません。

| メトリクス名 | 元のカウンター | 意味 |
|---|---|---|
| `linux-usage.tcp-opens.active` | `Tcp.ActiveOpens` | 能動的な接続開始（CLOSEDからSYN-SENTへの遷移）。接続成功数とは異なります。 |
| `linux-usage.tcp-opens.passive` | `Tcp.PassiveOpens` | 受動的な接続開始（LISTENからSYN-RCVDへの遷移）。 |

### TCP待ち受け: `linux-usage.tcp-listen.*`

`/proc/self/net/netstat` の累積カウンターを使います。出力値は **1分あたりの増加数（件/分）** です。

| メトリクス名 | 元のカウンター | 意味 |
|---|---|---|
| `linux-usage.tcp-listen.overflows` | `TcpExt.ListenOverflows` | 待ち受けソケットのacceptキューが満杯になったことによるオーバーフロー。 |
| `linux-usage.tcp-listen.drops` | `TcpExt.ListenDrops` | 待ち受けソケットで発生したパケットのドロップ。キューあふれ以外の原因も含みます。 |

TCP関連の値は `go-mackerel-plugin` が `(現在値 − 前回値) × 60 ÷ 経過秒数` で計算するため、小数になる場合があります。対応するカウンターが取得できない環境では、その項目は出力されません。

## 状態ファイルと注意事項

- 保存先は環境変数 `MACKEREL_PLUGIN_WORKDIR` です。未設定の場合はOSの一時ディレクトリ（通常 `/tmp`）を使います。指定するディレクトリは事前に作成し、実行ユーザーに読み書き権限を付与してください。
- CPU状態は `mackerel-plugin-linux-usage-<実効UID>` に保存します。JSON形式とファイル名は従来と互換性があり、`saferio` で読み取りとアトミックな更新を行います。TCP差分用の状態は `go-mackerel-plugin` が別ファイルに保存します。
- 同じ実効UID・保存先で複数の実行を行うとCPU状態を共有します。手動確認で定期実行の差分計算に影響を与えたくない場合は、別の `MACKEREL_PLUGIN_WORKDIR` を使ってください。
- 前回のCPU状態が600秒より古い場合、その回はメトリクス取得がエラーになります。CPU状態は更新されるため、次回実行から再計算できます。TCP差分にも600秒の上限があり、超過した回はログにエラーが記録され、TCP関連の値は0になります。
- CPUカウンターに変化がない場合は、CPU使用率を出力しません。再起動などでカウンターが減少した場合、CPUの負の差分は0として扱い、TCPの負の差分はログにエラーが記録され、0として出力されます。
- コンテナ内では、見えている `/proc` とネットワーク名前空間の情報を取得します。CPU使用率・ロードアベレージはcgroupのCPU制限に合わせた正規化ではありません。


## ライセンス

[MIT License](LICENSE)
