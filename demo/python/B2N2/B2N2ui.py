# -*- coding: utf-8 -*-
# Time：2023/1/28 12:25
# by Yoake

import json
import sys
import os
import time
import threading
import re

from PyQt5 import QtCore, QtGui, QtWidgets
import jsonpath
import requests
from lxml import etree


# post上传多线程
class PostThreads(threading.Thread):
    def __init__(self, func, args):
        threading.Thread.__init__(self)
        self.func = func
        self.args = args
        self.__flag = threading.Event()  # 用于暂停线程的标识
        self.__flag.set()  # 设置为True
        self.__running = threading.Event()  # 用于停止线程的标识
        self.__running.set()  # 将running设置为True

    def run(self) -> None:
        self.func(*self.args)


# GUI多线程
class MyThread(QtCore.QThread):
    postSignal = QtCore.pyqtSignal(str)
    finishedSignal = QtCore.pyqtSignal(str)

    def __init__(self, url, pageID, token, episode, episode_name):
        super(MyThread, self).__init__()
        self.url = url
        self.pageID = pageID
        self.token = token
        self.episode = episode
        self.episode_name = episode_name

    def get_partName(self):
        result = re.match('https://www.bilibili.com/', self.url)
        if result is None:
            url = f'https://www.bilibili.com/video/{self.url}/?p=1'
        else:
            url = self.url
        while True:
            try:
                response = requests.get(url)

                if response.status_code == 200:
                    content = response.text
                    tree = etree.HTML(content)
                    result = tree.xpath('//script[4]/text()')[0]

                    result = result.split('=', 1)[1]
                    result = result.split(';(function(){var s;')[0]

                    js = json.loads(result)

                    videoName = jsonpath.jsonpath(js, '$.videoData.title')[0]
                    if '/' in videoName:
                        videoName = videoName.replace('/', '／')
                    elif '\\' in videoName:
                        videoName = videoName.replace('\\', '＼')
                    else:
                        videoName = videoName
                    partName_list = jsonpath.jsonpath(js, '$..pages..part')
                    coverUrl = jsonpath.jsonpath(js, '$.videoData.pic')[0]

                    partName_str = '\n'.join(partName_list)

                    with open(f"./clc/list/{videoName}分集标题.txt", 'w', encoding='utf-8') as fp:
                        fp.writelines(partName_str)

                    print(f'已生成{videoName}的分集标题文件\n')

                    return videoName, partName_list, coverUrl
                else:
                    time.sleep(1)
                    continue
            except requests.exceptions.ConnectionError:
                time.sleep(1)
                continue

    def postNotion(self, url, body, headers, partName):
        while True:
            self.postSignal.emit(f"<u>《{partName}》</u>正在上传...")
            try:
                r = requests.post(url, json=body, headers=headers)
                if r.status_code == 200:
                    # if r.status_code == 200 or r.status_code == 201 or r.status_code == 202 or r.status_code == 204:
                    self.postSignal.emit(f"<u>《{partName}》</u>上传成功")
                    return
                else:
                    self.postSignal.emit(
                        f"<font color=red><u>《{partName}》</u>上传失败正在重试:<br>状态码<b>{r.status_code}</b></font>")
            except Exception as e:
                self.postSignal.emit(f"<font color=red><u>《{partName}》</u>上传失败正在重试:<br><b>{e}</b></font>")
                time.sleep(1)
                continue

    def run(self):
        videoName, partName_list, coverUrl = self.get_partName()
        token = self.token
        database_id = self.pageID
        episode = self.episode
        episode_name = self.episode_name
        count = len(partName_list)
        url = "https://api.notion.com/v1/pages"
        headers = {
            "Accept": "application/json",
            "Notion-Version": "2022-06-28",
            "Content-Type": "application/json",
            "Authorization": "Bearer " + token,
        }
        totalThreads = []
        for i in range(count):
            partName = partName_list[i]
            body = {
                "parent": {
                    "type": "database_id",
                    "database_id": database_id
                },
                "properties": {
                    episode: {"number": i + 1},
                    episode_name: {"title": [{"type": "text", "text": {"content": partName}}]},
                    "Name": {"select": {"name": videoName}},
                }
            }
            single_post = PostThreads(self.postNotion, (url, body, headers, partName))
            totalThreads.append(single_post)
            single_post.start()
        for t in totalThreads:
            t.join()
        self.finishedSignal.emit(f"<hr>视频《{videoName}》的封面链接为：<font color=blue>{coverUrl}</font>")


# 主界面
class MyUI(QtWidgets.QWidget):
    def __init__(self):
        super(MyUI, self).__init__()
        self.configButton = None
        self.thread = None
        self.stopButton = None
        self.startButton = None
        self.urlEdit = None
        self.label_4 = None
        self.runShower = None
        self.runLabel = None
        self.jsonBox = None
        self.label_3 = None
        self.setUpUI()
        self.loadConfig()

        self.child = None

    def setUpUI(self):
        self.resize(800, 425)
        self.setWindowTitle("B2N2")

        horizontalLayout_1 = QtWidgets.QHBoxLayout()
        verticalLayout_2_1 = QtWidgets.QVBoxLayout()
        horizontalLayout_3_1 = QtWidgets.QHBoxLayout()
        horizontalLayout_3_2 = QtWidgets.QHBoxLayout()
        horizontalLayout_1.setContentsMargins(0, 0, 0, 0)

        verticalLayout_2_1.setContentsMargins(10, 10, 10, 10)
        self.label_3 = QtWidgets.QLabel()
        self.label_3.setObjectName("label_3")
        self.label_3.setFont(QtGui.QFont("Arial", 10))
        self.label_3.setText("选择配置：")
        self.jsonBox = QtWidgets.QComboBox()
        self.jsonBox.setFont(QtGui.QFont("Arial", 15))
        verticalLayout_2_1.addLayout(horizontalLayout_3_1, 1)
        horizontalLayout_3_1.addWidget(self.label_3, 0)
        horizontalLayout_3_1.addWidget(self.jsonBox, 2)
        self.label_4 = QtWidgets.QLabel()
        self.label_4.setFont(QtGui.QFont("Arial", 11))
        self.label_4.setObjectName("label_4")
        self.label_4.setText("视频链接/<br>bv号：")
        self.urlEdit = QtWidgets.QLineEdit()
        self.urlEdit.setPlaceholderText("输入目标视频地址")
        self.urlEdit.setFont(QtGui.QFont("Arial", 15))
        horizontalLayout_3_2.addWidget(self.label_4)
        horizontalLayout_3_2.addWidget(self.urlEdit)
        self.startButton = QtWidgets.QPushButton()
        self.startButton.setText("开始")
        self.startButton.setFont(QtGui.QFont("Arial", 15))
        self.startButton.clicked.connect(self.startApp)
        verticalLayout_2_1.addLayout(horizontalLayout_3_1)
        verticalLayout_2_1.addLayout(horizontalLayout_3_2)
        verticalLayout_2_1.addWidget(self.startButton)
        self.stopButton = QtWidgets.QPushButton()
        self.stopButton.setFont(QtGui.QFont("Arial", 15))
        self.stopButton.setText("停止")

        self.stopButton.clicked.connect(self.stopApp)
        verticalLayout_2_1.addWidget(self.stopButton)
        self.configButton = QtWidgets.QPushButton()
        self.configButton.setText("创建配置文件")
        self.configButton.setFont(QtGui.QFont("Arial", 15))
        self.configButton.clicked.connect(self.createConfig)
        verticalLayout_2_1.addWidget(self.configButton)

        self.runLabel = QtWidgets.QLabel()
        self.runLabel.setText("运行记录")
        self.runShower = QtWidgets.QTextBrowser()
        self.runShower.setText("程序开始运行...")
        verticalLayout_2_3 = QtWidgets.QVBoxLayout()
        verticalLayout_2_3.setContentsMargins(10, 10, 10, 10)
        verticalLayout_2_3.addWidget(self.runLabel)
        verticalLayout_2_3.addWidget(self.runShower)

        horizontalLayout_1.addLayout(verticalLayout_2_1)
        horizontalLayout_1.addLayout(verticalLayout_2_3)
        horizontalLayout_1.setStretch(1, 1)
        self.setLayout(horizontalLayout_1)

    def loadConfig(self):
        if not os.path.exists('../clc'):
            os.mkdir('../clc')
            os.mkdir('../clc/key')
            os.mkdir('../clc/list')
        else:
            if not os.path.exists('../clc/key'):
                os.mkdir('../clc/key')
            if not os.path.exists('../clc/list'):
                os.mkdir('../clc/list')
        configs = os.listdir('../clc/key/')
        if len(configs) > 0:
            self.jsonBox.addItems(configs)
            self.runShower.append("<font color=green>本地存在配置，请选择配置</font>")
        else:
            self.runShower.append("<font color=red>本地未找到配置文件，请创建配置</font>")
        self.jsonBox.addItem("创建新配置")

    def startApp(self):
        choice = self.jsonBox.currentText()
        if choice == "创建新配置":
            self.createConfig()
        else:
            self.runShower.append(f"选择配置文件{choice}")
            with open('./clc/key/' + choice, 'r') as fp:
                configData = json.load(fp)

            url = self.urlEdit.text()
            pageID = configData["pageid"]
            token = configData["token"]
            episode = configData["episode"]
            episode_name = configData["episodename"]
            self.thread = MyThread(url, pageID, token, episode, episode_name)
            self.runShower.append("任务进行中...<hr>")
            self.thread.start()
            self.thread.finishedSignal.connect(self.finishedFunc)
            self.thread.postSignal.connect(self.showDetails)

    def showDetails(self, val):
        self.runShower.append(val)

    def createConfig(self):
        self.child = subUI()
        self.child.show()
        self.child.commitSignal.connect(self.loadAgain)

    def loadAgain(self, val):
        self.runShower.append(val)
        self.jsonBox.addItem(self.child.line3.text() + ".json")

    def stopApp(self):
        self.thread.terminate()
        self.runShower.append("停止任务!")

    def finishedFunc(self, val):
        self.runShower.append(val)
        self.runShower.append("<b>任务完成!!!</b>")


# 子界面
class subUI(QtWidgets.QWidget):
    commitSignal = QtCore.pyqtSignal(str)

    def __init__(self):
        super(subUI, self).__init__()
        self.setUpUI()

    def setUpUI(self):
        center_pointer = QtWidgets.QDesktopWidget().availableGeometry().center()
        x = center_pointer.x()
        y = center_pointer.y()
        self.setWindowTitle("创建配置文件")
        self.setFixedSize(600, 500)
        old_x, old_y, width, height = self.frameGeometry().getRect()
        self.move(x - int(width / 2), y - int(height / 2))

        fBox = QtWidgets.QFormLayout()

        self.label1 = QtWidgets.QLabel()
        self.label1.setFont(QtGui.QFont("Arial", 10))
        self.label1.setText("数据库ID：")

        self.label2 = QtWidgets.QLabel()
        self.label2.setFont(QtGui.QFont("Arial", 10))
        self.label2.setText("Token：")

        self.label3 = QtWidgets.QLabel()
        self.label3.setFont(QtGui.QFont("Arial", 10))
        self.label3.setText("配置文件名：")

        self.label4 = QtWidgets.QLabel()
        self.label4.setText("序号属性：")
        self.label4.setFont(QtGui.QFont("Arial", 10))

        self.label5 = QtWidgets.QLabel()
        self.label5.setText("单集标题：")
        self.label5.setFont(QtGui.QFont("Arial", 10))

        self.line1 = QtWidgets.QLineEdit()
        self.line1.setFont(QtGui.QFont("Arial", 15))
        self.line2 = QtWidgets.QLineEdit()
        self.line2.setFont(QtGui.QFont("Arial", 15))
        self.line3 = QtWidgets.QLineEdit()
        self.line3.setFont(QtGui.QFont("Arial", 15))
        self.line4 = QtWidgets.QLineEdit()
        self.line4.setFont(QtGui.QFont("Arial", 15))
        self.line5 = QtWidgets.QLineEdit()
        self.line5.setFont(QtGui.QFont("Arial", 15))

        self.commitBtn = QtWidgets.QPushButton()
        self.commitBtn.setText("确认")
        self.commitBtn.clicked.connect(self.commitConfig)

        self.commitBtn.setFont(QtGui.QFont("Arial", 15))
        fBox.addRow(self.label1, self.line1)
        fBox.addRow(self.label2, self.line2)
        fBox.addRow(self.label4, self.line4)
        fBox.addRow(self.label5, self.line5)
        fBox.addRow(self.label3, self.line3)
        fBox.addRow(self.commitBtn)
        fBox.setRowWrapPolicy(QtWidgets.QFormLayout.RowWrapPolicy(2))
        fBox.setVerticalSpacing(15)
        fBox.setContentsMargins(10, 20, 10, 10)
        self.setLayout(fBox)

    # 提交配置
    def commitConfig(self):
        pageID = self.line1.text()
        token = self.line2.text()
        episode = self.line4.text()
        episode_name = self.line5.text()
        with open(f"./clc/key/{self.line3.text()}.json", 'w', encoding='utf-8') as fp:
            config_data_dict = {"pageid": pageID, "token": token, "episode": episode, "episodename": episode_name}
            config_data = json.dumps(config_data_dict)
            fp.write(config_data)
        self.commitSignal.emit("创建配置文件成功")
        self.close()


if __name__ == '__main__':
    app = QtWidgets.QApplication(sys.argv)
    window = MyUI()
    window.show()
    sys.exit(app.exec_())
