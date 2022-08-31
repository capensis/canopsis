<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      service-dependencies.pa-3(
        :root="item",
        :columns="widget.parameters.serviceDependenciesColumns",
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
    widget: {
      type: Object,
      required: true,
    },
  },
};
</script>
