<template lang="pug">
  v-card
    v-card-title
      span.headline Watcher name
    v-divider
    v-card-text
      div Name: Watcher name
      div Org: Some org
      div.entities-list
        entity-data
</template>

<script>
import EntityData from '@/components/modals/watcher-data/entity-data.vue';

export default {
  components: {
    EntityData,
  },
};
</script>

<style scoped>
  .entities-list {
    margin-top: 20px;
  }
</style>
