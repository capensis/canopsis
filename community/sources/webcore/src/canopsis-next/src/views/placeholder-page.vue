<template>
  <c-layout gap="2" class="pa-4" column>
    <c-page-header :title="title" />

    <div class="text-body-2 text-medium-emphasis">
      {{ route.path }}
    </div>
  </c-layout>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();

const title = computed(() => route.meta.title || route.name);
</script>
