<template>
    <div class="video-wrapper">
        <n-form>
            <n-form-item label="视频索引">
                <n-input placeholder="请输入视频链接或视频ID(av号或bv号)" v-model:value="target"></n-input>
            </n-form-item>
            <n-form-item :showLabel="false">
                <n-space justify="space-evenly">
                    <n-button type="info" @click="getVideo" :loading="loading">搜索</n-button>
                    <n-button type="success" @click="submitVideoInfo" :loading="loading">提交</n-button>
                </n-space>
            </n-form-item>
            <n-progress 
                v-show="total_length != 0" 
                type="line" 
                :color="computed_color"
                :rail-color="changeColor(computed_color, { alpha: 0.2 })"
                :percentage="computed_progress" 
                :processing="loading"/>
        </n-form>
        <n-divider />
        <info-show-box :videoInfo="info" :index="completed_index"></info-show-box>
    </div>
</template>

<script setup lang="ts">
import { ref ,computed} from 'vue'
import { NForm, NFormItem, NInput, NButton, NDivider,NSpace,NProgress,useThemeVars } from 'naive-ui';
import { changeColor } from 'seemly'
import { GetVideoInfo, SubmitVideoInfo } from '../../../../wailsjs/go/main/App';
import { api } from '../../../../wailsjs/go/models';
import { InfoShowBox } from '@/components/common/index'
import {EventsOff, EventsOn} from '../../../../wailsjs/runtime'

let themeVars = useThemeVars()
let loading = ref(false)
let progress = ref(0)
let total_length = ref(0)
let target = ref('')
let info = ref<api.VideoInfo>()
let completed_index = ref<number[]>([])

let computed_progress = computed(()=>{
    return parseFloat((progress.value/total_length.value*100).toFixed(2))
})

let computed_color = computed(()=>{
    return computed_progress.value >= 100? themeVars.value.successColor:themeVars.value.infoColor
})
const getVideo = () => {
    loading.value = true
    GetVideoInfo(target.value).then((value) => {
        loading.value = false
        info.value = value
        total_length.value = value.titles.length
        progress.value = 0
    })
}
const updateProgressBar =()=>{
    progress.value += 1
    completed_index.value.push(progress.value)
}

const submitVideoInfo = () => {
    loading.value = true
    progress.value = 0
    EventsOff("postProgress")
    completed_index.value = []
    EventsOn("postProgress",updateProgressBar)
    SubmitVideoInfo(target.value).then(()=> {
        loading.value = false
        window.$message.success("数据上传完成")
    })
}
</script>

<style scoped>
.video-wrapper {
    text-align: left;
    .n-space {
        width: 100%;
        .n-button{
            width: 150%;
        }
    }
}
</style>