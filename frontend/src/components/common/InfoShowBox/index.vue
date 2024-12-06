<template>
    <div class="infoshow-wrapper">
        <div class="video-show" v-if="videoInfo">
            <div class="intro">
                <img :src="videoInfo?.cover" />
                <n-divider class="n-divider" vertical> </n-divider>
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
import { toRefs, ref, watch,h } from 'vue';
import {NDataTable,NDivider,NIcon } from 'naive-ui'
import type{DataTableColumns} from 'naive-ui';
import { api } from '../../../../wailsjs/go/models'
import {CheckmarkCircle,CloudUploadSharp} from "@vicons/ionicons5"
type Video = {
    no: number
    title: string
    status:boolean
}
const createColumns = (): DataTableColumns<Video> => {
    return [
        {
            title: 'No',
            key: 'no'
        }, {
            title: 'Title',
            key: 'title'
        },{
            title: 'Status',
            key: 'status',
            align:"center",
            render(row: Video){
                const index_temp = index?.value as Array<number>
                if (index_temp.includes(row.no)){
                    return h(NIcon,{
                        component: CheckmarkCircle,
                        color: 'green'
                    },{  })
                }else{
                    // 加载图标
                    return h(NIcon,{
                        component: CloudUploadSharp,
                        color:"grey"
                    },)
                }
            }
        }
    ]
}
let data = ref<Video[]>([])

const columns: DataTableColumns<Video> = createColumns()

const props = defineProps(["videoInfo", "bangumiInfo","index"])

let { videoInfo, bangumiInfo,index } = toRefs(props)

watch(() => videoInfo?.value, (c, o) => {
    if (c != undefined && c != o) {
        data.value = []
        let titles: string[] = JSON.parse(JSON.stringify(videoInfo?.value.titles))
        titles.forEach((value, index) => {
            let item: Video = { no: index + 1, title: value , status:false }
            data.value?.push(item)
        })
    }
})

watch(() => bangumiInfo?.value, (c, o) => {
    if (c != undefined && c != o) {
        data.value = []
        let main: api.Section[] = JSON.parse(JSON.stringify(bangumiInfo?.value.main))
        main.forEach((value, index) => {
            let item: Video = { no: index + 1, title: value.long_title ,status:false}
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
        .n-divider{
            height: 3rem;
        }
        .title{
            height: 3rem;
            font-size: 1.2rem; 
            font-weight: 600;
            line-height: 3rem;
        }
    }
    
}
</style>