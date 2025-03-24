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

复制b站上某个分p视频的链接到视频标签页的输入框内，先`搜索`确认是否是对应的视频及其分p标题**目前这一步不要跳过！**再点击`提交`，等待一段时间后完成导入。

由于notion设置了请求速率限制（3个/秒），当上传数量过于庞大时会增加上传的错误几率，属于正常现象。

# 介绍
B2N3即 `bilibili to notion version 3`，是`serial pages generator for notion`项目下`B2N`的迭代版本，基于[wails](https://github.com/wailsapp/wails)客户端GUI框架实现，仅支持Windows10即以上和Mac的PC端——使用系统自带的Webview2。


**注意：当前版本仅支持database对象，即只有且只能有一个数据库的页面**
# 使用步骤
## 第一步：获得目标database的id
### 1.获得database的url链接
网页端直接在地址栏获得database的url

桌面端和移动端需要通过界面右上角三个点下的Copy link复制得到url
### 2.提取url链接中页面的id
![image.png](https://yoaken-1316330335.cos.ap-chongqing.myqcloud.com/markdownPic/202301030532401.png)

把`?v=`之前的这一段从url链接截取下来即可

## 第二步：创建可以操作notion的integration并获得其token令牌
以下操作都在notion的intergration管理页面进行[点击前往该页面](https://www.notion.so/my-integrations)
### 1.创建integration
创建intergration，并确认要有插入新内容(Insert content)的权限

![image.png](https://yoaken-1316330335.cos.ap-chongqing.myqcloud.com/markdownPic/202301030632968.png)

### 2.获取intergration的token令牌
查看自己的intergration的相关信息，在密钥(Secrets)下查看(Show)并复制(Copy)该intergration的token令牌

![image.png](https://yoaken-1316330335.cos.ap-chongqing.myqcloud.com/markdownPic/202301030633303.png)

## 第三步：确认database的相关配置
### 1.指定控制当前页面的intergration
回到database页面，点击右上角的三个小点，Add connections让刚刚创建好的intergration和这个页面连接起来，使之可以操作这个页面
### 2.确认序号创建的属性properties
该属性的类型必须是number，属性名为Episode
### 3.如果需要分组，那么就确认分组的属性properties
该属性的类型必须是select，属性名为Name
