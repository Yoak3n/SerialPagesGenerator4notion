<template>
    <div class="input-file-name">
        <n-modal 
            v-model:show="showModal"
            preset="dialog"
            positive-text="确认"
            @positive-click="changeState"
            @close="props.close"
            :show-icon="true"
            type="info"
            style="text-align: left;"
         >
            <n-input v-model:value="inputConfigName" placeholder="请输入配置文件名称"></n-input>
            
        </n-modal>

    </div>


</template>


<script setup lang="ts">
import {ref ,toRefs} from 'vue'
import {NModal,NInput} from 'naive-ui'


let inputConfigName = ref('')
const props = defineProps(['showModal','input','close'])
let {showModal} = toRefs(props)

const changeState = () => {
    props.input(inputConfigName.value)
}

</script>

<style scoped>
</style>