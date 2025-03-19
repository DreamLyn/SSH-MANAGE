<h1 align="center">🔑 ssh-manage</h1>
<div align="center">

[![Stars](https://img.shields.io/github/stars/dreamlyn/ssh-manage?style=flat)](https://github.com/dreamlyn/ssh-manage)
[![Forks](https://img.shields.io/github/forks/dreamlyn/ssh-manage?style=flat)](https://github.com/dreamlyn/ssh-manage)
[![Docker Pulls](https://img.shields.io/docker/pulls/dreamlyn/ssh-manage?style=flat)](https://hub.docker.com/r/dreamlyn/ssh-manage)
[![Release](https://img.shields.io/github/v/release/dreamlyn/ssh-manage?sort=semver)](https://github.com/dreamlyn/ssh-manage/releases)
[![License](https://img.shields.io/github/license/dreamlyn/ssh-manage)](https://mit-license.org/)

</div>
<div align="center">

中文 ｜ [English](README_EN.md)

</div>

---

## 🚩 项目简介

我们经常遇到需要在陌生的 PC 环境下对公司以及某些个人服务器进行运维，这时候重新查找服务器的账号密码让人痛苦不堪，SSH- MANAGE 就是为了解决这个问题，在 NAS 或者服务器上安装后，我们通过浏览器就可以对服务器进行管理。

ssh-manage 是一个仿 Termius 的 Web 版 SSH 终端，支持 **SSH 连接** 和 **端口转发**，它具有以下优势：

本地部署：一键安装，只需要下载二进制文件，然后直接运行即可。同时也支持 Docker 部署、源代码部署等方式。​
数据安全：由于是私有部署，所有数据均存储在自己的服务器上，不会经过第三方，确保数据的隐私和安全。​
操作简单：简单配置即可轻松对服务器进行管理。

**功能特点：**

- 🌐 **Web 端 SSH 连接**：随时随地访问你的服务器
- 🔄 **端口转发**：支持远程版（Docker 安装）端口转发
- 🔒 **安全管理**：支持密钥登录，保障连接安全
- 📂 **多会话管理**：同时管理多个服务器连接
- 📦 **Docker 部署**：一键启动，轻量级安装

🚀 [**Github 项目地址**](https://github.com/dreamlyn/ssh-manage)

---

## **快速启动**

**5 分钟部署 SSH-MANAGE！**

**二进制部署参考**
从 [GitHub Releases](https://github.com/dreamlyn/ssh-manage/releases) 页面下载预先编译好的二进制可执行文件压缩包，解压缩后在终端中执行：

```bash
./ssh-manage serve
```

浏览器中访问 `http://127.0.0.1:8090`。

初始的管理员账号及密码：

- 账号：`lyn@dreamlyn.cn`
- 密码：`1234567890`

即刻使用 ssh-manage

**docker-compose 部署参考**

创建 docker 网络
`docker network create -d bridge --attachable=true docker_default`

docker-compose.yml文件参考
```
version: "3"
services:
  ssh-manage:
    image: registry.cn-hangzhou.aliyuncs.com/dreamlyn_i/ssh-manage:latest
    container_name: ssh-manage
    volumes:
      - ./data:/app/pb_data
    environment:
      - PUID=1000
      - PGID=1000
      - TZ=Asia/Shanghai
    ports:
      - "8090:8090"
      - "3306:3306" #请修改为具体需要端口映射的值（不用端口映射可删除此行）
    restart: unless-stopped
    networks:
      - docker_default
networks:
  docker_default:
    external: true
```

---

## ⭐ 运行界面

![Screenshot](https://img.dreamlyn.cn:8443/i/2025/03/19/111534.webp)

![Screenshot](https://img.dreamlyn.cn:8443/i/2025/03/19/111801.webp)

## 🤝 参与贡献

ssh-manage 是一个免费且开源的项目，采用 [MIT License](./LICENSE.md)。你可以使用它做任何你想做的事，甚至把它当作一个付费服务提供给用户。

你可以通过以下方式来支持 ssh-manage 的开发：

- 提交代码：如果你发现了 Bug 或有新的功能需求，而你又有相关经验，可以提交代码给我们。
- 提交 Issue：功能建议或者 Bug 可以[提交 Issue](https://github.com/dreamlyn/ssh-manage/issues) 给我们。

支持更多UI 的优化改进、Bug 修复、文档完善等，欢迎参与贡献。

## ⛔ 免责声明

ssh-manage 基于 MIT License 发布，完全免费提供，旨在“按现状”供用户使用。作者及贡献者不对使用本软件所产生的任何直接或间接后果承担责任，包括但不限于性能下降、数据丢失、服务中断、或任何其他类型的损害。

**无任何保证**：本软件不提供任何明示或暗示的保证，包括但不限于对特定用途的适用性、无侵权性、商用性及可靠性的保证。

**用户责任**：使用本软件即表示您理解并同意承担由此产生的一切风险及责任。

## 🌐 加入社群

- 微信群聊（超 200 人需邀请入群，可先加作者好友）

  <img src="https://img.dreamlyn.cn:8443/i/2025/03/19/113416.webp" width="240"/>

<!-- ## 🚀 Star 趋势图

[![Stargazers over time](https://starchart.cc/dreamlyn/ssh-manage.svg?variant=adaptive)](https://starchart.cc/dreamlyn/ssh-manage) -->

