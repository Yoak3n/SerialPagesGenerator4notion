# 描述
由于notion的database无法批量生成带序号的子页面，在某些严格的应用场景下属于痛点，因此用python写了一个做这件事的小工具，并做了简单的打包。

**注意：当前版本仅支持database对象，即只有且只能有一个数据库的页面**

# 使用步骤
## 第一步：获得目标database的id
### 1.获得database的url链接
网页端直接在地址栏获得database的url
桌面端和移动端需要通过界面右上角三个点下的Copy link复制得到url
### 2.提取url链接中页面的id
![image.png](https://yoaken-1316330335.cos.ap-chongqing.myqcloud.com/markdownPic/202301030532401.png)

把“?v=”之前的这一段从url链接截取下来即可

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
该属性的类型必须是number，记住属性名
### 3.如果需要分组，那么就确认分组的属性properties
该属性的类型必须是select，记住属性名

## 第四步：运行程序
### 1.首次运行
#### 初始化配置
首次运行程序会在程序所在的目录下创建一个clc目录来收集程序生成的文件(ps：当前版本只在clc目录下创建key目录用来存放target.json配置文件)，程序会要求根据提示依次输入：
* （1）目标database的id
* （2）intergration的token
* （3）序号所在的属性名
* （4）分组的属性名，不填就表示不分组
* （5）如果分组，还需要填组名
#### 输入要创建的页面数量
#### 完成创建后，在命令行按回车键(Enter)确认结束

### 2.后续运行

#### 是否修改配置文件
程序会询问是否修改配置，输入y就会重写配置文件，
不填就使用当前本地的配置
# 开发想法（并非计划）
## 多配置文件
后续可能会开发本地保存多个配置文件，当前版本要实现保存多个配置文件并选择使用其中某一个需要手动进行切换：把要使用的配置文件改名为target，其他配置文件改为另外的文件名。
## GUI程序
如果未来实现了多配置文件，那么开发GUI也应该提上日程。有了图形界面之后操作就能更加简洁，比如就不用每次运行程序都询问是否修改配置文件，再比如可以直接选择使用哪个配置文件。
问题在于开发GUI程序导入PyQt库之后这个python程序的大小不可避免的膨胀，也就不再轻量化。当然这仅仅只是当下的一个忧虑，具体还是要看结果到底如何。
