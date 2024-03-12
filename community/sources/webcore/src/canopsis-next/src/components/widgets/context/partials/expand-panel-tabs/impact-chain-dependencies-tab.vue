<template>
  <service-dependencies
    :root="item"
    :columns="columns"
    class="pa-3"
    impact
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
