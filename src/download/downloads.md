Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.2.1 (08-03-2018)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.2.1.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.darwin-amd64.tgz) | darwin | amd64 | 3.9M | 856823b99bb5fb1c6fece9c90d048a3f55444561fd2d020cade328f9c4851000 |
[dapperdox-1.2.1.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.linux-amd64.tgz) | linux | amd64 | 3.9M | de536eeff36cefb26d2d61d636c50edc00574183490b6379554a9b62f49af646 |
[dapperdox-1.2.1.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.linux-arm.tgz) | linux | arm | 3.6M | 38224995b2160f2ce4ca33dad9a458148619a4f3bf200f08eb760161cbc7f0f8 |
[dapperdox-1.2.1.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.linux-arm64.tgz) | linux | arm64 | 3.7M | 81101825691f1ba086809c9326070c695fc954323f0ed16439be4e233382936b |
[dapperdox-1.2.1.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.linux-x86.tgz) | linux | x86 | 3.7M | 7b9424eaf6b29f22b5709e860ea740f2ec11ae59fde05a6b98c6f94c0d2a095c |
[dapperdox-1.2.1.windows-amd64.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.windows-amd64.zip) | windows | amd64 | 3.9M | 10b2e7095a8665e61c769fd270f4afd40f6bf288cae3f52a22292306ae80c0e0 |
[dapperdox-1.2.1.windows-x86.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.2.1/dapperdox-1.2.1.windows-x86.zip) | windows | x86 | 3.7M | a0e23e524b836fe0a7bfe8a8c4f9551d5ba3c1927b45d83e35e84fd96459b584 |

**1.1.1 (06-04-2017)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.1.1.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.darwin-amd64.tgz) | darwin | amd64 | 4.0M | f80297c68efa43502c1e98e6eef508ffe9df91854ce93127e62781d9b0617919 |
[dapperdox-1.1.1.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-amd64.tgz) | linux | amd64 | 4.0M | 3e959b0d972bd4035a46a45810b3afc3faf64d6c7a8173aed4ba7dd1a7a1e846 |
[dapperdox-1.1.1.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-arm.tgz) | linux | arm | 3.6M | 8ceaba3296a865c2b734534ffa19403c07b6b0548f8b1cf5d6f25c32ede522d5 |
[dapperdox-1.1.1.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-arm64.tgz) | linux | arm64 | 3.7M | 7bfd70731fb1ad250872be854b99330f93a66c6b38b69b6b235a5dde5f341240 |
[dapperdox-1.1.1.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.linux-x86.tgz) | linux | x86 | 3.8M | 8576e869d66ffb7ce1bdfaf3e267fcc42471c2ace2095bcf63ae9b31c99b7b46 |
[dapperdox-1.1.1.windows-amd64.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.windows-amd64.zip) | windows | amd64 | 4.0M | 56326726780700a8b48504e6366a123893bbaf818957b7ac9d7680e4dcf2eea4 |
[dapperdox-1.1.1.windows-x86.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.1.1/dapperdox-1.1.1.windows-x86.zip) | windows | x86 | 3.8M | 8052904d8eb209ae28e03b41b6cbdeb0bff0326351820dcc5a97d597efa2726f |

## Building from source

To build from source, clone the [GitHub repo](https://github.com/dapperdox/dapperdox):

```bash
> git clone https://github.com/dapperdox/dapperdox
```

Now build dapperdox (this assumes that you have your [golang](https://golang.org/doc/install) environment configured correctly):

```
go get && go build
```

Alternatively, just type `make`
