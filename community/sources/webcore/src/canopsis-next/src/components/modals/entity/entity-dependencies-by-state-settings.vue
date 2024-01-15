<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      {{ $t('stateSetting.computeMethod') }}: {{ stateMethodName }}
      <entity-dependencies-by-state-settings-component :entity="config.entity" />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesByStateSettingsComponent from '@/components/other/entity/entity-dependencies-by-state-settings.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.entityDependenciesByStateSettings,
  components: {
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

    stateMethodName() {
      /**
       * TODO: Should be changed on real data
       */
      return '%rule name%';
    },
  },
};
</script>
