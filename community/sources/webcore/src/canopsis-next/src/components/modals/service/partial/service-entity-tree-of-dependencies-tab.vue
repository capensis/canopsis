<template lang="pug">
  service-dependencies.pa-2(
    :root="dependency",
    :columns="columns",
    light,
    openable-root
  )
</template>

<script>
import { serviceToServiceDependency } from '@/helpers/treeview/service-dependencies';

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
  computed: {
    dependency() {
      return serviceToServiceDependency(this.entity);
    },
  },
};
</script>
