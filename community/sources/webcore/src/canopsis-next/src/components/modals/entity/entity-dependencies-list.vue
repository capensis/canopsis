<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ title }}</span>
    </template>
    <template #text="">
      <entity-dependencies-list-component
        :widget="widget"
        :columns="widget.parameters.widgetColumns"
        :entity-id="config.entityId"
        :impact="config.impact"
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { generatePreparedDefaultContextWidget } from '@/helpers/entities/widget/form';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesListComponent from '@/components/other/entity/entity-dependencies-list.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entityDependenciesList,
  components: {
    ModalWrapper,
    EntityDependenciesListComponent,
  },
  mixins: [modalInnerMixin],
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.entityDependenciesList.title');
    },

    widget() {
      return this.config.widget ?? generatePreparedDefaultContextWidget();
    },
  },
};
</script>
