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
        </n-form>
        <n-divider />
        <info-show-box :videoInfo="info"></info-show-box>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NForm, NFormItem, NInput, NButton, NDivider,NSpace } from 'naive-ui';
import { GetVideoInfo, SubmitVideoInfo } from '../../../../wailsjs/go/main/App';
import { api } from 'wailsjs/go/models';
import { InfoShowBox } from '@/components/common/index'

let loading = ref(false)
let target = ref('')
let info = ref<api.VideoInfo>()

const getVideo = () => {
    loading.value = true
    GetVideoInfo(target.value).then((value) => {
        loading.value = false
        info.value = value
    })
}


const submitVideoInfo = () => {
    loading.value = true
    SubmitVideoInfo().then((datas)=> {
        console.log(datas)
        loading.value = false
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