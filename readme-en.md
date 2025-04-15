## SyncDNS

A dns synchronization tool for AdGuardHome.

Query the latest resolution result of the domain name based on the DNS rewrite record in AdGuardHome, delete the invalid resolution record in AdGuardHome, and add the new resolution record to AdGuardHome.

## Usage

1. Modify the configuration file
    ```bash
    cp example.config.yaml config.yaml
    ```
2. Compile and run
    ```bash
    go build .
    ./github.com/qingchuwudi/syncdns -c ./config.yaml
    ```

## Description

1. If there is no 'domain' in the configuration file or 'domain' is empty, the DNS rewrite record will be queried in AdGuardHome and periodically synchronized on this basis.
2. If there is a 'domain' in the configuration file, it will be synchronized based on the 'domain' configuration.

The reason for the above strategy is that when we use AdGuardHome, in addition to the foreign domain name needs to be synchronized, there will also be some domain name resolution rewrites used by the LAN (such as rewriting baidu.com to 127.0.0.1).
To avoid these custom parsing records being overwritten, use the strategy above.

todo: In the future, we will consider adding an ignore parameter to ignore certain domain names and not resolve and synchronize them.
