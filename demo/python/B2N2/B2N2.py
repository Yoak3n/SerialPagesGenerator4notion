# encoding=utf-8
# Time：2023/1/2 10:35
# by Yoake
import json
import os
import time
import re

from tqdm import tqdm
import requests, lxml.etree, jsonpath


def check_config():
    if not os.path.exists('../clc'):
        os.mkdir('../clc')
        os.mkdir('../clc/key')
        os.mkdir('../clc/list')
    else:
        if not os.path.exists('../clc/key'):
            os.mkdir('../clc/key')
        if not os.path.exists('../clc/list'):
            os.mkdir('../clc/list')
    exist = os.path.exists('../clc/key/mykey.json')
    if exist:
        check = input("是否使用新配置文件？y/n")
        if check == 'y' or check == 'Y' or check == 'yes' or check == "YES" or check == "Yes":
            no_edit = False
        else:
            no_edit = True
    else:
        no_edit = False

    if no_edit:
        with open('../clc/key/mykey.json', 'r', encoding='utf-8') as fp:
            config_data = json.load(fp)
            page_id = config_data["pageid"]
            token = config_data["token"]
            episode = config_data["episode"]
            episode_name = config_data["episodename"]
            title = config_data["title"]

    else:
        page_id = input('请输入数据库id：')
        token = input('请输入机器人token：')
        episode = input('请输入集数的序号所在的属性名：')
        episode_name = input('请输入每一p的标题所在的属性名：')
        title = input('请输入该视频的标题所在的属性名：')
        with open('../clc/key/mykey.json', 'w', encoding='utf-8') as fp:
            config_data_dict = {"pageid": page_id, "token": token, "episode": episode, "episodename": episode_name,
                                "title": title}
            config_data = json.dumps(config_data_dict)
            fp.write(config_data)

    return page_id, token, episode, episode_name,title


def remove_illegal(video_name: str) ->str:
    illegal = ['?','/','\\',':','<','>','\"','*','|']
    for i in illegal:
        if i in video_name:
            video_name = video_name.replace(i,'')
    return video_name

def get_partName(bid):
    result = re.match('https://www.bilibili.com/video/', bid)
    if result is None:
        bvid = bid.lstrip("BV")
        url =  f'https://api.bilibili.com/x/web-interface/view?bvid={bvid}'
    else:
        bvid = bid[result.end():].lstrip('BV')
        url = f"https://api.bilibili.com/x/web-interface/view?bvid={bvid}"
    while True:
        try:
            response = requests.get(url)

            if response.status_code == 200:
                js = json.loads(response.text)

                videoName = jsonpath.jsonpath(js, '$.data.title')[0]
                
                partName_list = jsonpath.jsonpath(js, '$..pages..part')
                
                coverUrl = jsonpath.jsonpath(js, '$.data.pic')[0]
                

                partName_str = '\n'.join(partName_list)
                videoName = remove_illegal(videoName).strip()
                with open(f"../clc/list/{videoName}分集标题.txt", 'w', encoding='utf-8') as fp:
                    fp.writelines(partName_str)

                print(f'已生成{videoName}的分集标题文件\n')
                
                return videoName, partName_list, coverUrl
            else:
                time.sleep(1)
                continue
        except requests.exceptions.ConnectionError:
            time.sleep(1)
            continue


def post_notion(total, part, database_id, token, episode, episode_name,title):
    count = len(part)
    bar = tqdm(total=count)
    number = 0
    if episode == '' or episode is None:
        episode = 'Episode'
    if episode_name == '' or episode_name is None:
        episode_name = 'EpisodeName'
    if title == '' or title is None:
        title = 'Name'
    bar.set_description_str('正在上传notion...')

    url = "https://api.notion.com/v1/pages"

    headers = {
        "Accept": "application/json",
        "Notion-Version": "2022-06-28",
        "Content-Type": "application/json",
        "Authorization": "Bearer " + token,
    }

    for i in range(count):

        while True:
            body = {
                "parent": {
                    "type": "database_id",
                    "database_id": database_id
                },
                "properties": {
                    episode: {"number": i + 1},
                    episode_name: {"title": [{"type": "text", "text": {"content": part[i]}}]},
                    title: {"select": {"name": total}},
                }
            }
            try:
                r = requests.post(url, json=body, headers=headers)
                if r.status_code == 200:
                    bar.update(1)
                    break
                else:
                    number += 1
                    time.sleep(2)
                    continue
            except requests.exceptions.ConnectionError:
                if number <= 10:
                    number += 1
                    print(f"\033[0;91;40m上传notion过于频繁的请求，正在等待重试...第{number}次\033[0m")
                    time.sleep(1)

                    continue
                else:
                    print(f"\033[0;91;40m上传notion请求重试超过10次，程序终止，请检查配置和网络！\033[0m")
                    break

    bar.set_description_str('上传notion已完成')
    bar.close()
    if number > 0:
        print(f'总计上传错误{number}次')
    else:
        print('\n上传零错误')
    return

if __name__ == '__main__':
    try:
        page_id, token, episode, episode_name, title = check_config()
        bid = input("请输入分集视频的BV号：")
        if bid != '':
            videoname, partname_list, coverUrl = get_partName(bid)
            post_notion(videoname, partname_list, page_id, token, episode, episode_name, title)
            print(f'{videoname}的封面链接是：{coverUrl}')
        else:
            print('请检查是否输入正确的BV号')

    except Exception as ex:
        print(ex)

    finally:
        check = input('请按回车确认并结束程序')
        if check != '谁都不会打出来的一行字':
            print('程序结束，感谢使用！')
