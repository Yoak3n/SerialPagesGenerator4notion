import {RouteRecordRaw} from 'vue-router'


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
                component:()=>import('@/views/setting/index.vue'),
                name:'setting'
            }

        ]
    }
]

export default routes