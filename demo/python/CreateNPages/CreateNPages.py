# encoding=utf-8
# Time：2023/1/2 12:25
# by Yoake
import time, os, json

from tqdm import tqdm
import requests


def post_notion(count, index_name, database_id, token, group_name, group):
    print('')
    # 创建进度条对象
    bar = tqdm(total=count)
    # 上传过程错误记数
    number = 0

    bar.set_description_str('正在上传notion...')
    url="https://api.notion.com/v1/pages"
    headers = {
        "Accept": "application/json",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json",
        "Authorization": "Bearer " + token,
    }
    # 批量创建页面的循环
    for i in range(count):
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
        while True:
            try:
                r = requests.post(url, json=body, headers=headers)
                if r.status_code == 200:
                    bar.update(1)
                    break
                else:
                    number += 1
                    time.sleep(2)
                    continue
            except requests.exceptions.ConnectionError :

                print("\n\033[1;91;40m过于频繁的请求，正在等待重试...\033[0m")
                number +=1
                time.sleep(2)
                continue

    bar.set_description_str('\033[0;92;40m上传notion已完成\033[0m')
    bar.close()
    if number>0:
        print(f'\n总计上传错误{number}次')
    else:
        print('\n上传过程无错误')


def get_config():
    if not os.path.exists('../clc'):
        os.mkdir('../clc')
        os.mkdir('../clc/key')
    else:
        if not os.path.exists('../clc/key'):
            os.mkdir('../clc/key')

    exist = os.path.exists('../clc/key/target.json')
    # 如果没有本地配置文件
    if exist:
        check = input("是否重新配置文件？y/n")
        if check =='y' or check=='Y' or check == 'yes' or check=="YES" or check=="Yes":
            noedit =False
        else:
            noedit = True
    else:
        noedit = False
    # 如果不需要重新配置
    if noedit:
        with open('../clc/key/target.json', 'r', encoding='utf-8') as fp:
            js = json.load(fp)
            page_id = js["page_id"]
            token = js["token"]
            index_name = js["index_name"]
            group_name = js["group_name"]
            group = js["group"]
        print('已加载配置文件')
    elif exist :
        print("修改配置文件，请根据提示依次输入")
        page_id = input('请输入数据库id：')
        token = input('请输入机器人token：')
        index_name = input("请输入生成序号所在的属性名：")
        group_name = input("请输入分组的属性名：") # 为空即不分组
        if group_name !='':
            group =input("请输入分组的组名：")
        else :
            group = ''
        dict = {"page_id": page_id, "token": token, "index_name": index_name, "group_name": group_name, "group": group}
        with open('../clc/key/target.json', 'w', encoding='utf-8') as fp:
            fp.write(json.dumps(dict))
    else:
        print("未找到配置文件，请根据提示依次输入")
        page_id = input('请输入数据库id：')
        token = input('请输入机器人token：')
        index_name = input("请输入生成序号所在的属性名：")
        group_name = input("请输入分组的属性名：")  # 为空即不分组
        if group_name != '':
            group = input("请输入分组的组名：")
        else:
            group = ''

        dict = {"page_id": page_id,"token": token,"index_name": index_name,"group_name":group_name,"group":group}
        with open('../clc/key/target.json', 'w', encoding='utf-8') as fp:
            fp.write(json.dumps(dict))

    return index_name, page_id, token, group_name, group


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
