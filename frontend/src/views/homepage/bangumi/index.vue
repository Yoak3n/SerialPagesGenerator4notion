<template>
    <div class="bangumi-wrapper">
        <n-form>
            <n-form-item label="视频索引">
                <n-input placeholder="请输入视频链接或视频ID(ss)" v-model:value="target"></n-input>
            </n-form-item>
            <n-form-item :showLabel="false">
                <n-space justify="space-evenly">
                    <n-button type="info" @click="getBangumi" :loading="loading">搜索</n-button>
                    <n-button type="success" @click="getBangumi" :loading="loading">提交</n-button>
                </n-space>
                
            </n-form-item>
        </n-form>
        <n-divider />
        <info-show-box :bangumiInfo="info"></info-show-box>
    </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {NForm,NFormItem,NSpace,NInput,NButton,NDivider} from 'naive-ui'
import { api } from '../../../../wailsjs/go/models';
import { InfoShowBox } from '@/components/common/index'
import { GetBangumiInfo } from '../../../../wailsjs/go/main/App';

let loading = ref(false)
let target = ref('')
let info = ref<api.Bangumi>()

const getBangumi = ()=>{
    loading.value = true
    GetBangumiInfo(target.value).then((value) => {
        loading.value = false
        info.value = value
    })
}

</script>

<style scoped>
.bangumi-wrapper {
    text-align: left;

    .n-space {
        width: 100%;
        .n-button{
            width: 150%;
        }
    }
}

</style>