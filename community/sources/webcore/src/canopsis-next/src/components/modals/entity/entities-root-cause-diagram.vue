<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <state-settings-summary :entity="config.entity" />
      <entity-dependencies-by-state-settings :entity="config.entity" :color-indicator="config.colorIndicator" />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesByStateSettings from '@/components/other/entity/entity-dependencies-by-state-settings.vue';
import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entitiesRootCauseDiagram,
  components: {
    StateSettingsSummary,
    ModalWrapper,
    EntityDependenciesByStateSettings,
  },
  mixins: [modalInnerMixin],
  computed: {
    entity() {
      return this.config.entity ?? {};
    },

    title() {
      return this.config.title ?? this.$t('modals.entityDependenciesList.title', {
        name: this.entity.name,
      });
    },
  },
};
</script>
