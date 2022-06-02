<template lang="pug">
  v-card.secondary.lighten-2(flat)
    v-card-text
      service-dependencies.pa-3(
        :root="rootDependency",
        :columns="widget.parameters.serviceDependenciesColumns",
        impact,
        openable
      )
</template>

<script>
import ServiceDependencies from '@/components/other/service/table/service-dependencies.vue';

import { serviceToServiceDependency } from '@/helpers/treeview/service-dependencies';

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
  computed: {
    rootDependency() {
      return serviceToServiceDependency(this.item);
    },
  },
};
</script>
