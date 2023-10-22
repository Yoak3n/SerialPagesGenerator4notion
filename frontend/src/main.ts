import {createApp} from 'vue'
import App from './App.vue'

const app = createApp(App)

import router from './routers/index'
app.use(router)

import pinia from './store/index'
app.use(pinia)

app.mount('#app')
