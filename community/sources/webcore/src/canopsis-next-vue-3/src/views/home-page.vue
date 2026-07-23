<script setup>
import { computed } from 'vue';

const template = computed(() => '<router-link :to="{ path: \'/login\' }">Login</router-link>');
</script>

<template>
  <div>
    <h1>Home</h1>
    <c-compiled-template :template="template" />
  </div>
</template>

<style scoped></style>
