<script lang="ts" setup>
import { onMounted } from "vue";
import {NaiveProvider} from "@/components/common/index";
import { ChangeConfiguration } from "wailsjs/go/main/App";
import { config } from "wailsjs/go/models";


const initSetting = () => {
  
}

onMounted(()=>{
  
})
</script>

<template>
  <NaiveProvider>
    <RouterView></RouterView>
  </NaiveProvider>
  
</template>

<style>

</style>
