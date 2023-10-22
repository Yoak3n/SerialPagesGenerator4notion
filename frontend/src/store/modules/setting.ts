import {defineStore} from 'pinia'


let useSettingStore = defineStore('settingStore',{
    state:()=>{
        return {
            token:''
        }
    }
})



export default useSettingStore