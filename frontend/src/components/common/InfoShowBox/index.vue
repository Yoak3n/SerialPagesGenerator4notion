<template>
    <div class="infoshow-wrapper">
        <div class="video-show" v-if="videoInfo">
            <div class="intro">
                <img :src="videoInfo?.cover" />
                <div class="title">{{ videoInfo.name }}</div>
            </div>
            
            <div class="cover-url">
                封面链接：
                <a :href="videoInfo?.cover" target="_blank" rel="noopener noreferrer">{{ videoInfo?.cover }}</a>
            </div>
            <div class="data-table-wrapper">
                <n-data-table :data="data" :columns="columns">

                </n-data-table>
            </div>
        </div>
        <div class="bangumi-show" v-if="bangumiInfo">
            <div class="intro">
                <div class="title">{{ bangumiInfo?.name }}</div>
                <img :src="bangumiInfo?.cover" />
            </div>

            <div class="cover-url">
                封面链接：
                <a :href="bangumiInfo.bangumi_cover" target="_blank" rel="noopener noreferrer">{{ bangumiInfo.bangumi_cover }}</a>
            </div>
            <div class="data-table-wrapper">
                <n-data-table :data="data" :columns="columns">

                </n-data-table>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { toRefs, ref, watch } from 'vue';
import {NDataTable } from 'naive-ui'
import { api } from 'wailsjs/go/models'
type Video = {
    no: number
    title: string
}
const createColumns = () => {
    return [
        {
            title: 'No',
            key: 'no'
        }, {
            title: 'Title',
            key: 'title'
        }
    ]
}
let data = ref<Video[]>([])

const columns = createColumns()

const props = defineProps(["videoInfo", "bangumiInfo"])

let { videoInfo, bangumiInfo } = toRefs(props)

watch(() => videoInfo?.value, (c, o) => {
    
    if (c != undefined && c != o) {
        data.value = []
        let titles: string[] = JSON.parse(JSON.stringify(videoInfo?.value.titles))
        titles.forEach((value, index) => {
            let item: Video = { no: index + 1, title: value }
            data.value?.push(item)
        })
    }

})
watch(() => bangumiInfo?.value, (c, o) => {
    if (c != undefined && c != o) {
        data.value = []
        let main: api.Section[] = JSON.parse(JSON.stringify(bangumiInfo?.value.main))
        console.log(main);
        main.forEach((value, index) => {
            let item: Video = { no: index + 1, title: value.long_title }
            data.value?.push(item)
        })

    }
})
</script>

<style scoped>
.infoshow-wrapper {
    .intro{
        height: 3rem;
        position: relative;
        display: flex;
        justify-items: center;
        .title{
            height: 3rem;
            font-size: 1.2rem; 
            font-weight: 600;
            line-height: 3rem;
        }
    }
    
}
</style>