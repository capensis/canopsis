<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      service-dependencies.pa-3(
        :root="item",
        :columns="columns",
        impact,
        openable-root
      )
</template>

<script>
import ServiceDependencies from '@/components/other/service/table/service-dependencies.vue';

export default {
  components: {
    ServiceDependencies,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
