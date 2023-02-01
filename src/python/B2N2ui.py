# -*- coding: utf-8 -*-
# Time：2023/1/28 12:25
# by Yoake

import json
import sys
import os
import time

from PyQt5 import QtCore, QtGui, QtWidgets

import B2N2


class MyThread(QtCore.QThread):
    finishedSignal = QtCore.pyqtSignal(str)

    def __init__(self, url, pageID, token):
        super(MyThread, self).__init__()
        self.url = url
        self.pageID = pageID
        self.token = token

    def run(self):
        videoName, partName_list, coverUrl = B2N2.get_partName(self.url)
        B2N2.post_notion(videoName, partName_list, self.pageID, self.token)
        self.finishedSignal.emit(f"视频《{videoName}》的封面链接为：{coverUrl}")


class myUI(QtWidgets.QWidget):
    def __init__(self):
        super(myUI, self).__init__()
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

        self.child = subUI()

    def setUpUI(self):
        self.resize(684, 425)
        self.setWindowTitle("B2N2")

        horizontalLayout_1 = QtWidgets.QHBoxLayout()
        verticalLayout_2_1 = QtWidgets.QVBoxLayout()
        verticalLayout_2_3 = QtWidgets.QVBoxLayout()
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
        self.label_4.setFont(QtGui.QFont("Arial", 10))
        self.label_4.setObjectName("label_4")
        self.label_4.setText("视频链接：")
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

        self.setLayout(horizontalLayout_1)

    def loadConfig(self):
        if not os.path.exists('./clc'):
            os.mkdir('./clc')
            os.mkdir('./clc/key')
            os.mkdir('./clc/list')
        else:
            if not os.path.exists('./clc/key'):
                os.mkdir('./clc/key')
            if not os.path.exists('./clc/list'):
                os.mkdir('./clc/list')
        configs = os.listdir('./clc/key/')
        if len(configs) > 0:
            self.jsonBox.addItems(configs)
        else:
            self.runShower.append("本地未找到配置文件，请创建配置")
        self.jsonBox.addItem("创建新配置")

    def startApp(self):
        choice = self.jsonBox.currentText()
        if choice == "创建新配置":
            self.createConfig()
        else:
            self.runShower.append(f"选择配置文件{choice}")
            with open('./clc/key/'+choice,'r') as fp:
                configData = json.load(fp)

            url = self.urlEdit.text()
            pageID = configData["pageid"]
            token = configData["token"]
            self.thread=MyThread(url, pageID, token)
            self.runShower.append("任务进行中...")
            self.thread.start()
            self.thread.finishedSignal.connect(self.finishedFunc)

    def createConfig(self):
        self.child.show()
        self.jsonBox.addItem(self.child.line3.text())

    def stopApp(self):
        self.thread.terminate()
        self.runShower.append("停止任务!")

    def finishedFunc(self, val):
        self.runShower.append(val)
        self.runShower.append("任务完成!!!")


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
        self.setFixedSize(500, 300)
        old_x, old_y, width, heigh = self.frameGeometry().getRect()
        self.move(x - int(width / 2), y - int(heigh / 2))

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
        self.line1 = QtWidgets.QLineEdit()
        self.line1.setFont(QtGui.QFont("Arial", 15))
        self.line2 = QtWidgets.QLineEdit()
        self.line2.setFont(QtGui.QFont("Arial", 15))
        self.line3 = QtWidgets.QLineEdit()
        self.line3.setFont(QtGui.QFont("Arial", 15))
        self.commitBtn = QtWidgets.QPushButton()
        self.commitBtn.setText("确认")
        self.commitBtn.clicked.connect(self.commitConfig)

        self.commitBtn.setFont(QtGui.QFont("Arial", 15))
        fBox.addRow(self.label1, self.line1)
        fBox.addRow(self.label2, self.line2)
        fBox.addRow(self.label3, self.line3)
        fBox.addRow(self.commitBtn)
        fBox.setRowWrapPolicy(QtWidgets.QFormLayout.RowWrapPolicy(2))
        fBox.setVerticalSpacing(15)
        fBox.setContentsMargins(10, 20, 10, 10)
        self.setLayout(fBox)

    def commitConfig(self):
        pageID = self.line1.text()
        token = self.line2.text()
        with open(f"./clc/key/{self.line3.text()}.json",'w',encoding='utf-8') as fp:
            config_data_dict = {"pageid": pageID, "token": token}
            config_data = json.dumps(config_data_dict)
            fp.write(config_data)
        self.commitSignal.emit("创建配置文件成功")
        self.close()


if __name__ == '__main__':
    app = QtWidgets.QApplication(sys.argv)
    window = myUI()
    window.show()
    sys.exit(app.exec_())
