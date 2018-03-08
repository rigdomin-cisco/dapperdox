Title: Download DapperDox
Description: Download the latest pre-built release of DapperDox
Keywords: Download, linux, mac, osx, windows, binary, tar, zip

# Download DapperDox

We provide the lastest release builds for the most common operating systems and architectures.
If yours is not listed here, then clone the source from [GitHub](http://github.com/dapperdox/dapperdox) and follow the [build instructions](#building-from-source).

## Precompiled releases

**1.2.2 (08-03-2018)**

| Filename | OS   | Arch | Size | Checksum |
| -------- | ---- | ---- | ---- | -------- |
[dapperdox-1.2.2.darwin-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.darwin-amd64.tgz) | darwin | amd64 | 3.9M | eafa9db96213cea5d6fa75ad0f6e59883ed353724d17596545ae04d1774c0531 |
[dapperdox-1.2.2.linux-amd64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.linux-amd64.tgz) | linux | amd64 | 3.9M | d339757c75d392eb4c3ee7e30ff80a9540fcb104224db29f5c8c79d87c05cfe4 |
[dapperdox-1.2.2.linux-arm.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.linux-arm.tgz) | linux | arm | 3.6M | e752099b249405a0b229ed98f789d0588305c71df8cae1c92e19dd784717b911 |
[dapperdox-1.2.2.linux-arm64.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.linux-arm64.tgz) | linux | arm64 | 3.7M | 02accebf4af676ed7425558b1c9ea03cedf2d95f6817ae2604fef7088a4ac3d0 |
[dapperdox-1.2.2.linux-x86.tgz](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.linux-x86.tgz) | linux | x86 | 3.7M | f953191256a81804843c04ce1d77157d33ab489b12c75869a25b27389293d924 |
[dapperdox-1.2.2.windows-amd64.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.windows-amd64.zip) | windows | amd64 | 3.9M | 4b2a6a61f788d5313d1da26b4b4dfa8580f2cc4f435a693aa6c0afa76c95876a |
[dapperdox-1.2.2.windows-x86.zip](https://github.com/DapperDox/dapperdox/releases/download/v1.2.2/dapperdox-1.2.2.windows-x86.zip) | windows | x86 | 3.7M | d7d75a12bf90aad6553517d3218f53f861004c8bdf552db1dac090c7bfe9fcd6 |

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
