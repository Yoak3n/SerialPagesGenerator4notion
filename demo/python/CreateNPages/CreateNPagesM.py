# encoding=utf-8
# Time：2023/2/5 23:52
# by Yoake

import time
import json
import os
from threading import Thread

import requests
from tqdm import tqdm


class Single_threads(Thread):
    def __init__(self, func, args):
        Thread.__init__(self)
        self.func = func
        self.args = args

    def run(self) -> None:
        self.func(*self.args)


def post_notion(times, index_name, database_id, token, group_name, group):
    print('')
    # 创建进度条对象
    bar = tqdm(total=times)
    bar.set_description_str('正在上传notion...')
    headers = {
        "Accept": "application/json",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json",
        "Authorization": "Bearer " + token,
    }
    # 批量创建页面的循环
    thread_list = []
    for i in range(times):
        if group_name != '':
            body = {
                "parent": {
                    "type": "database_id",
                    "database_id": database_id
                },
                "properties": {
                    index_name: {"number": i + 1},
                    group_name: {"select": {"name": group}},
                }
            }
        else:
            body = {
                "parent": {
                    "type": "database_id",
                    "database_id": database_id
                },
                "properties": {
                    index_name: {"number": i + 1}
                }
            }
        # 构建死循环来捕获请求过于频繁的错误，防止因这个错误而终止程序
        my_thread = Single_threads(func=single_post, args=(body, headers, bar))
        thread_list.append(my_thread)
        my_thread.start()

    for t in thread_list:
        t.join()
    bar.set_description_str('\033[0;92;40m上传notion已完成\033[0m')
    bar.close()


def single_post(body, headers, bar):
    url = "https://api.notion.com/v1/pages"
    while True:
        try:
            r = requests.post(url, json=body, headers=headers)
            if r.status_code == 200:
                bar.update(1)
                return
            else:
                time.sleep(1)
                continue
        except requests.exceptions.ConnectionError:
            print("\033[1;91;40m过于频繁的请求，正在等待重试...\033[0m")
            time.sleep(1)
            continue


def get_config():
    if not os.path.exists('../clc'):
        os.mkdir('../clc')
        os.mkdir('../clc/key')
    else:
        if not os.path.exists('../clc/key'):
            os.mkdir('../clc/key')

    # 如果没有本地配置文件
    configs = os.listdir("../clc/key")
    length = len(configs)
    show_list = [f"{i + 1}.{configs[i]} " for i in range(length)]
    print(''.join(show_list) + f'{length + 1}.创建配置')
    choice = input("选择配置文件？")
    if choice == '':
        choice = '1'
    if int(choice) == length+1:
        noedit = False
    else:
        noedit = True
    # 如果不需要重新配置
    if noedit:
        with open(f'../clc/key/{configs[int(choice)-1]}', 'r', encoding='utf-8') as fp:
            js = json.load(fp)
            pageid = js["pageid"]
            token = js["token"]
            index_name = js["index_name"]
            group_name = js["group_name"]
            group = js["group"]
        print('已加载配置文件')
    else:
        print("正在创建配置文件，请根据提示依次输入")
        json_name = input('请输入将创建的配置文件名（无拓展名）：')
        pageid = input('请输入数据库id：')
        token = input('请输入机器人token：')
        index_name = input("请输入生成序号所在的属性名：")
        group_name = input("请输入分组的属性名（为空即不分组）：")
        if group_name != '':
            group = input("请输入分组的组名：")
        else:
            group = ''
        dict = {"pageid": pageid,"token": token, "index_name": index_name, "group_name": group_name, "group": group}
        with open(f'../clc/key/{json_name}.json', 'w', encoding='utf-8') as fp:
            fp.write(json.dumps(dict))

    return index_name, pageid, token, group_name, group


if __name__ == '__main__':
    try:
        config = get_config()
        count = int(input("请输入要创建的数据个数："))
        post_notion(count, config[0], config[1], config[2], config[3], config[4])
    except Exception as e:
        print(e)
    finally:
        check = input('请按回车确认并结束程序')
        if check != '谁都不会打出来，打出来也没用的一行字':
            print('\n程序结束，感谢使用！')