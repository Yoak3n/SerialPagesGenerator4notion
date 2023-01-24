const fs = require('fs');
const path =require('path');
const axios =require('axios');
const util = require('util')

const ProgressBar = require('./progress-bar');
const prompt = require('prompt-sync')()


function check(nPath){
  nPath = path.resolve(__dirname, nPath)
  const flag = fs.existsSync(nPath)
  if (!flag){
    fs.mkdir(nPath,err=>{
    if (err) throw err;
  });
  }else return flag;
}

function entry() {
  check('./clc')
  check('./clc/key')
  let config ={}
  promise = new Promise((resolve,reject)=>{
    (async () => {
      try {
        const data = await util.promisify(fs.readFile)('./clc/key/mykey.json', { encoding: 'utf-8' });
        config = JSON.parse(data)
        resolve(config)
      } catch (err) {
        if (err.code === 'ENOENT') {
          console.error('配置文件不存在，请输入配置...')
          let config ={
            databaseId:"",
            token:""
          }
          config.databaseId = prompt("请输入数据库id：")
          config.token=prompt("请输入机器人token：")
          fs.writeFile('./clc/key/mykey.json',JSON.stringify(config),(err)=>{
          })
          resolve(config)
        } else {
          throw err;
        }
      }
    })();
  })
  return promise
}

function getHtml(targetUrl){
  let data = []
  var promise = new Promise((resolve,reject)=>{
  axios.get(targetUrl).then((res)=>{
    const body =res.data;

    const start_index = /<script>window.__INITIAL_STATE__=/.exec(body).index + 33;
    const end_index= /\;\(function\(\)\{var s\;/.exec(body.slice(start_index)).index+start_index  ;
    const jsonStr = body.slice(start_index,end_index);

    check('./.json')
    fs.writeFile('./\.json/videoData.json',jsonStr,err=>{
      if(err){
          console.log(err);
          }
      });
    let jsonObj = JSON.parse(jsonStr);
    let videoName = jsonObj['videoData']['title']
    if (videoName.includes('/')){
        videoName = videoName.replace('/','／')
    }else if (videoName.includes('\\')){
      videoName = videoName.replace('\\','＼')
    }else{
      videoName = videoName;}

    let parts = jsonObj.videoData.pages;
    let nameArr =[];

    parts.forEach(element => {
      const name = element.part;
      nameArr.push(name);
    });
    const names=nameArr.join('\n');
    check('./clc/list/')
    fs.writeFile('./clc/list/'+videoName+'.txt',names,err=>{
        if(err){
            console.log(err);
        }
        return names;
    });
    
    Promise.all([videoName,nameArr]).then( (dataArray) => {
      data = dataArray
      resolve(data)
    })
    
    
    }).catch(err=>{
    console.log(err);
    })
  })
  return promise
}

function post2notion(nameArr,title,database_id,token){
  const notionUrl='https://api.notion.com/v1/pages';
  headers = {
    "Accept": "application/json",
    "Notion-Version": "2022-06-28",
    "Content-Type": "application/json",
    "Authorization": "Bearer " + token,
    }
  let i =0
  const pb = new ProgressBar("正在上传至notion",50)
  var num = 0 ,total = nameArr.length
  nameArr.forEach(item =>{
    body = makeBodys(item,i,title,database_id)
    i++
    
    axios.defaults.retry = 5
    axios.defaults.retryDelay = 1000
 
    axios.interceptors.response.use(undefined, function axiosRetryInterceptor(err) {
        var config = err.config;
        // 如果配置不存在或未设置重试选项，则拒绝
        if (!config || !config.retry) return Promise.reject(err);
    
        // 设置变量以跟踪重试次数
        config.__retryCount = config.__retryCount || 0;
    
        // 判断是否超过总重试次数
        if (config.__retryCount >= config.retry) {
            // 返回错误并退出自动重试
            return Promise.reject(err);
        }
    
        // 增加重试次数
        config.__retryCount += 1;
    
        //打印当前重试次数
        console.log(config.url +' 自动重试第' + config.__retryCount + '次');
    
        // 创建新的Promise
        var backoff = new Promise(function (resolve) {
            setTimeout(function () {
                resolve();
            }, config.retryDelay || 1);
        });
    
        // 返回重试请求
        return backoff.then(function () {
            return axios(config);
        });
    });

    axios.post(
      notionUrl,
      body,
      {headers:headers}
    ).then(function downloading(res){
      let code = res.status
        if (num<=total && code==200){
        // 更新进度条
        num++;
        pb.render({ completed: num, total: total });
        
        return
      }
      
    })
  })

}

function makeBodys(part,i,title,database_id){
  body = {
    "parent": {
        "type": "database_id",
        "database_id": database_id
    },
    "properties": {
        "Episode": {"number": i+1 },
        "EpisodeName": {"title": [{"type": "text", "text": {"content": part}}]},
        "Name": {"select": {"name": title}}
    }
  }
  body =JSON.stringify(body)
  return body
}





entry().then(value=>{
  const databaseId = value.databaseId
  const myToken = value.token
  let targetUrl =prompt('请输入输入视频地址：')
  const p2= getHtml(targetUrl)

  p2.then((value=>{
      const [title,nameArr] = value
      post2notion(nameArr,title,databaseId,myToken)
  }))

})
  

