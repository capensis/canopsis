<template>
  <v-layout class="gap-3" column>
    <div v-if="title" class="text-body-2">
      {{ title }}
    </div>
    <v-btn
      v-for="link in links"
      :key="link.title"
      :href="link.href"
      color="primary"
    >
      {{ link.title }}
    </v-btn>
  </v-layout>
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
    links: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
