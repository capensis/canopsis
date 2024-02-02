<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <state-settings-summary :pending="stateSettingPending" :entity="config.entity" :state-setting="stateSetting" />

      <v-layout
        v-if="stateSettingPending"
        class="my-10"
        justify-center
        align-center
      >
        <v-progress-circular color="primary" indeterminate />
      </v-layout>
      <entity-dependencies-by-state-settings v-else :entity="config.entity" />
    </template>
  </modal-wrapper>
</template>

<script>
import { MODALS } from '@/constants';

import { infosToArray } from '@/helpers/entities/shared/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { checkStateSettingMixin } from '@/mixins/entities/check-state-setting';

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
  mixins: [modalInnerMixin, checkStateSettingMixin],
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.entityDependenciesList.title');
    },

    entity() {
      return this.config.entity;
    },
  },
  mounted() {
    this.checkStateSetting({
      name: this.entity.name,
      type: this.entity.type,
      infos: infosToArray(this.entity.infos),
      impact_level: this.entity.impact_level,
    });
  },
};
</script>
