<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <state-settings-summary />
      <entity-dependencies-by-state-settings-component :entity="config.entity" />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesByStateSettingsComponent from '@/components/other/entity/entity-dependencies-by-state-settings.vue';
import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entityDependenciesByStateSettings,
  components: {
    StateSettingsSummary,
    ModalWrapper,
    EntityDependenciesByStateSettingsComponent,
  },
  mixins: [modalInnerMixin],
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.entityDependenciesList.title');
    },

    entity() {
      return this.config.entity;
    },
  },
};
</script>
