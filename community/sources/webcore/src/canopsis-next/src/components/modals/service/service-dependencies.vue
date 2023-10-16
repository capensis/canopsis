<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <service-dependencies-table
        :root="config.root"
        :columns="config.columns"
        :impact="config.impact"
        :openable-root="config.openableRoot"
        include-root
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import ServiceDependenciesTable from '@/components/other/service/partials/service-dependencies.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.serviceDependencies,
  components: { ServiceDependenciesTable, ModalWrapper },
  mixins: [modalInnerMixin],
  computed: {
    title() {
      const type = this.config.impact
        ? 'impacts'
        : 'dependencies';

      return this.$t(`modals.serviceDependencies.${type}.title`, { name: this.config.root?.name });
    },
  },
};
</script>
