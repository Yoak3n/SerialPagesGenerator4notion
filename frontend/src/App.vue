<script lang="ts" setup>
import { onMounted } from "vue";
import {NaiveProvider} from "@/components/common/index";
</script>

<template>
  <NaiveProvider>
    <RouterView></RouterView>
  </NaiveProvider>
  
</template>

<style>

</style>
