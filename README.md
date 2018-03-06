# unityHttpServerSample

・使い方

```shell command
cd project_dir
./build.sh
./run.sh [start/stop/restart]
```

port:9999でlistenします。
unityから`/sample1/request, /sample2/request`にPOSTしてください。

■ 疎通確認
`
http://localhost:9999/web/ping
`

