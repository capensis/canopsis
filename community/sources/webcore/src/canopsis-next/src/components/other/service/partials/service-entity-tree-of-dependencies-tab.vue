<template>
  <service-dependencies
    :root="entity"
    :columns="columns"
    openable-root
  />
</template>

<script>
import ServiceDependencies from '@/components/other/service/partials/service-dependencies.vue';

export default {
  components: { ServiceDependencies },
  props: {
    entity: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      required: false,
    },
  },
};
</script>
