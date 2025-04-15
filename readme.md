## SyncDNS

[English Readme](readme-en.md)

一款用于AdGuardHome的dns同步工具。

根据AdGuardHome中的DNS重写记录，查询其中域名的最新解析结果，并把失效的解析记录在AdGuardHome中删除、新增的解析记录添加到AdGuardHome中。

## 用法

1. 修改配置文件
    ```bash
    cp example.config.yaml config.yaml
    ```
2. 编译并运行
    ```bash
    go build .
    ./github.com/qingchuwudi/syncdns -c ./config.yaml
    ```

## 说明

1. 如果配置文件中没有 `domain` 或者 `domain` 为空，就会在AdGuardHome中查询DNS重写记录并以此为基础进行定期同步。
2. 如果配置文件中有 `domain`，就会以 `domain` 配置为基准进行同步。

上述策略的原因是我们在使用AdGuardHome的时候，除了国外域名需要同步以外，还会有一些局域网使用的域名解析重写（比如把baidu.com 改写为 127.0.0.1），
为了避免这些自定义解析记录被覆盖，采用上面的策略。

todo: 后续会考虑增加忽略参数，用以忽略某些域名，不对其做解析和同步。
