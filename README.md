# unityHttpServerSample

・使い方

```sh
cd project_dir
sh build.sh
./bin/unitySample
```

port:9999でlistenします。
unityから`/sample1/request, /sample2/request`にPOSTしてください。

■ 疎通確認
`
http://localhost:9999/web/ping
`

