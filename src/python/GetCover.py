# encoding=utf-8

import json
import time
import re
import requests,lxml.etree,jsonpath

def GetCover(url):
    while True:
        try:
            response = requests.get(url)

            if response.status_code == 200:
                content = response.text

                tree = lxml.etree.HTML(content)
                result = tree.xpath('//script[4]/text()')[0]

                result = result.split('=',1)[1]
                result = result.split(';(function(){var s;')[0]

                js = json.loads(result)
                
                videoName = jsonpath.jsonpath(js, '$.videoData.title')[0]
                coverUrl = jsonpath.jsonpath(js, '$.videoData.pic')[0]
                
                
                return videoName, coverUrl
            else:
                time.sleep(1)
                continue
        except requests.exceptions.ConnectionError:
            time.sleep(1)
            continue
        
if __name__ == '__main__':

    try:
        bid = input("请输入分集视频的BV号：")
        if bid != '':
            result = re.match('https://www.bilibili.com/', bid)
            if result is None:
                url =  f'https://www.bilibili.com/video/{bid}/?p=1'
            else:
                url = bid
            videoName, coverUrl = GetCover(url)
            print(f'{videoName}的封面链接是：{coverUrl}')
        else:
            print('请检查是否输入正确的BV号')

    except Exception as ex:
        print(ex)

    finally:
        check = input('请按回车确认并结束程序')
        if check != '谁都不会打出来的一行字':

            print('程序结束，感谢使用！')