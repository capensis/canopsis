<template lang="pug">
  service-dependencies.pa-2(
    :root="entity",
    :columns="columns",
    light,
    openable-root
  )
</template>

<script>
import ServiceDependencies from '@/components/other/service/table/service-dependencies.vue';

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
