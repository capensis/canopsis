<template>
  <service-dependencies
    class="pa-3"
    :root="item"
    :columns="columns"
    impact
    openable-root
  />
</template>

<script>
import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';

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
