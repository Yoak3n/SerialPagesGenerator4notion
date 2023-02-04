# encoding=utf-8
# Time：2023/2/2 2:39
# by Yoake

import json
import sys
import os

from PyQt5 import QtCore, QtGui, QtWidgets


class myUI(QtWidgets.QWidget):
    def __init__(self):
        super(myUI, self).__init__()
        self.setUpUI()

    def setUpUI(self):
        mainLayout = QtWidgets.QGridLayout()
        self.lable1 = QtWidgets.QLabel()
        self.lable2 = QtWidgets.QLabel()
        self.jsonBox = QtWidgets.QComboBox()




        self.setLayout(mainLayout)

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

