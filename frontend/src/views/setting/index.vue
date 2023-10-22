<template>
    <div class="setting-wrapper">
        <!-- 设置 -->
        <n-form>
            <n-form-item label="配置文件" size="large">
                <n-space justify="space-between" class="space">
                    <n-select class="select" :options="options" :loading="selectLoading" v-model:value="optionValue"
                        @update:value="readSetting" @click="loadSetting">
                    </n-select>

                    <n-button-group>
                        <n-button type="success" @click="Confirm">选择/创建</n-button>
                        <n-button type="error" @click="deleteConfiguration">删除</n-button>
                    </n-button-group>
                </n-space>

            </n-form-item>
            <n-form-item label="数据库ID">
                <n-input class="input" v-model:value="database_id" placeholder="请输入数据库ID" />
            </n-form-item>
            <n-form-item label="机器人token">
                <n-input class="input" v-model:value="token" placeholder="请输入机器人令牌" />
            </n-form-item>
        </n-form>
        <input-config-name v-model:showModal="showModal" :input="inputName" :close="closeModal" />

    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NForm, NFormItem, NInput, NSelect, NButton, SelectOption, NSpace, NButtonGroup } from 'naive-ui'

import { InputConfigName } from "@/components/modal/index";

import { Configuration } from './types'
import { ScanConfiguraitonFiles, CreateConfiguration, GetCurrentConfiguration, ChangeConfiguration, DeleteConfigurationFile } from '../../../wailsjs/go/main/App'
import { config } from '../../../wailsjs/go/models'


let options = ref<SelectOption[]>([])
let optionValue = ref("")

let showModal = ref(false)


let database_id = ref("")
let token = ref("")

let selectLoading = ref(false)

onMounted(async () => {
    await loadSetting()
})

const loadSetting = async () => {
    // 读取配置
    selectLoading.value = true
    // 读取文件配置缓冲
    let optionTemp: SelectOption[] = []
    let result = await ScanConfiguraitonFiles()

    for (let i = 0; i < result.length; i++) {
        let reg = result[i].split('.json')[0]
        optionTemp.push({
            value: reg,
            label: result[i]
        })
        options.value = optionTemp
        optionValue.value = options.value[0].value as string
    }
    selectLoading.value = false
    readSetting(optionValue.value)
}


const Confirm = async () => {
    let command: string = optionValue.value
    if (command == "创建新配置") {
        showModal.value = true
    } else {
        let conf = config.Config.createFrom({
            database_id: database_id.value,
            token: token.value,
        })
        let res = await ChangeConfiguration(conf)
        if (res === 0) {
            window.$message.success("配置修改成功！", { duration: 2000 })

        }
    }

}

const deleteConfiguration = () => {
    DeleteConfigurationFile(optionValue.value).then((res)=>{
        if(res===0){
            window.$message.success("配置删除成功！", { duration: 2000 })
            loadSetting()
        }else{
            window.$message.error("配置删除失败！", { duration: 2000 })
            return
        }
    })
}

const inputName = async (name: string) => {
    if (name === "") {
        window.$message.error("配置名称不能为空")
        return
    }
    if (name.length > 100) {
        window.$message.error("配置名称不能超过100个字符")
        return
    }
    if (name.length < 1) {
        window.$message.error("配置名称不能少于1个字符")
        return
    }
    showModal.value = false
    let conf = config.Config.createFrom({
        database_id: database_id.value,
        token: token.value
    })
    let result = await CreateConfiguration(conf, name)
    if (result == "配置文件创建成功") {
        window.$message.success(result, { duration: 2000 })
    } else {
        window.$message.warn(result, { duration: 2000 })
    }
}

const closeModal = () => {
    showModal.value = false
}

const readSetting = (value: string) => {
    GetCurrentConfiguration(value).then((value) => {
        database_id.value = value.database_id
        token.value = value.token
    })
}

</script>

<style scoped lang="less">
.setting-wrapper {
    .space {
        width: 100%;
    }
    text-align: left;
}
</style>