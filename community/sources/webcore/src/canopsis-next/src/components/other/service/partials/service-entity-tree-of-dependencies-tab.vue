<template>
  <service-dependencies
    :root="entity"
    :columns="columns"
    :type="type"
    openable-root
    show-state-setting
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
    type: {
      type: Number,
      required: false,
    },
  },
};
</script>
