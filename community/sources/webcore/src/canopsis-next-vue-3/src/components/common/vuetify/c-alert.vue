<script setup>
defineOptions({
  inheritAttrs: false,
});
</script>

<template>
  <v-alert v-bind="$attrs" class="c-alert">
    <slot />
  </v-alert>
</template>

<style>
.c-alert {
  width: 100%;
}
</style>
