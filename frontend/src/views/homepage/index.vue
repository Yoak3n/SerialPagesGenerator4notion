<template>
    <div class="homepage-wrapper">

        <n-layout has-sider style="height: 100vh">
            <n-layout-sider 
                show-trigger 
                collapse-mode="width" 
                :collapsed-width="64" 
                :collapsed="collapsed" 
                bordered
                :width="180"
                @collapse="collapsed = true"
                @expand="collapsed = false">
                <!-- <div class="sider-bar"> -->
                <n-menu :options="menuOptions" v-model:value="activeKey" :collapsed-width="64" :collapsed="collapsed"
                    :collapsed-icon-size="24.5">
                </n-menu>
                <n-button text style="font-size: 36px" size="large" class="setting" @click="openSetting">
                    <template #icon>
                        <n-icon>
                            <settings />
                        </n-icon>
                    </template>
                </n-button>


                <!-- </div> -->

            </n-layout-sider>
            <n-layout-content content-style="padding: 24px;">
                <RouterView></RouterView>
            </n-layout-content>

        </n-layout>


        <!-- <n-space justify="space-around" align="center" vertical>
                <n-button size="large" class="view-btn">
                    视频
                </n-button>
                <n-button size="large" class="view-btn">
                    番剧
                </n-button>
                <n-button text style="font-size: 24px" size="large">
                    <template #icon>
                        <n-icon>
                            <settings />
                        </n-icon>
                    </template>
                </n-button>
            </n-space> -->
        <n-divider vertical />
        <RouterView></RouterView>
    </div>
</template>

<script setup lang="ts">
import { h, ref, Component, onMounted } from 'vue'
import { RouterLink,useRouter } from 'vue-router'
import { NButton, NIcon, NDivider, NMenu, NLayout, NLayoutSider, NLayoutContent } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { Settings, Videocam, PlayCircle } from '@vicons/ionicons5'


const $router = useRouter()
let activeKey = ref('')
let collapsed = ref(false)
// 渲染图标函数
const renderIcon = (icon: Component) => {
    return () => h(NIcon, null, { default: () => h(icon) })
}
// 菜单选项
const menuOptions: MenuOption[] = [
    {
        label: () =>
            h(
                RouterLink,
                {
                    to: '/video',

                },
                { default: () => '视频' }),
        key: 'video',
        icon: renderIcon(Videocam),
    },
    {
        label: () =>
            h(
                RouterLink,
                {
                    to: '/bangumi'
                },
                { default: () => '番剧' }),
        key: 'bangumi',
        icon: renderIcon(PlayCircle),
    }

]


const openSetting = () => {
    $router.push('/setting')
    activeKey.value = 'setting'
}

// 初始化
onMounted(()=>{
    activeKey.value = 'video'
    $router.push('/video')
})





</script>

<style scoped lang="less">
.homepage-wrapper {
    overflow: hidden;

    .sider-bar {
        overflow: hidden;
        background-color: antiquewhite;
        align-items: center;
        height: 100vh;
        display: flex;
        flex-direction: column;
        justify-content: space-between;


    }

    .n-layout-sider {
        text-align: start;
        line-height: 2.625rem;
        background-color: #eee;

        a {
            font-size: 24.5px;
        }

        // .setting-wrapper{
        //     text-align: center;
        // }
        .n-button {
            width: 100%;
            margin: 0 auto;
        }
    }

    .setting {
        text-align: center;
        position: absolute;
        bottom: 1rem;
    }

    .n-menu {
        width: 100%;
    }



}
</style>