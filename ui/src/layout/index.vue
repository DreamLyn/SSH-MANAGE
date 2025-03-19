<template>
  <div :class="classObj" class="app-wrapper">
    <navbar  />
    <app-main />
  </div>
</template>

<script setup>
import { AppMain, Navbar } from './components'
</script>

<style lang="scss" scoped>
  @use "@/assets/styles/mixin.scss" as *;
  @use "@/assets/styles/variables.module.scss" as *;

.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;

}
</style>