<template>
  <service-dependencies
    :root="item"
    :columns="columns"
    :type="type"
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
    type: {
      type: Number,
      required: false,
    },
  },
};
</script>
