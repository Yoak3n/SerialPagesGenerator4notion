import {RouteRecordRaw} from 'vue-router'
import Setting from '@/views/setting/index.vue'

const routes:RouteRecordRaw[] = [
    {
        path:'/',
        component:()=>import('@/views/homepage/index.vue'),
        name:'homepage',
        children:[
            {
                path:'/video',
                component:()=>import('@/views/homepage/video/index.vue'),
                name:'video'
            },{
                path:'/bangumi',
                component:()=>import('@/views/homepage/bangumi/index.vue'),
                name:'bangumi'
            },{
                path:'/setting',
                component:Setting,
                name:'setting'
            }

        ]
    }
]

export default routes