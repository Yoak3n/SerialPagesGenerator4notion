B2N3即 `bilibili to notion version 3`，是`serial pages generator for notion`项目下`B2N`的迭代版本，基于[wails](https://github.com/wailsapp/wails)客户端GUI框架实现，仅支持Windows10即以上和Mac的PC端——使用系统自带的Webview2。

**目前B2N3仅初步完成导入分p视频到notion数据库的功能，且运行尚不稳定，请酌情使用。**

## 编译

```bash
git clone https://github.com/Yoak3n/SerialPagesGenerator4notion.git
cd SerialPagesGenerator4notion
wails build
```

wails交叉编译暂未成功，如有需求可配置好相应环境（`go`+`pnpm`+`wails`）后自行编译

## 运行

运行程序后进入设置界面完成notion自动化相关配置，即可开始使用，此后每次上传前需再次确认配置内容。

复制b站上某个分p视频的链接到视频标签页的输入框内，先`搜索`确认是否是对应的视频及其分p标题**目前这一步不要跳过！**，再点击`提交`，等待一段时间后完成导入。

由于notion设置了请求速率限制（3个/秒），当上传数量过于庞大时会增加上传的错误几率，属于正常现象。