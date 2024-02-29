<template>
  <modal-wrapper close>
    <template #title="">
      {{ title }}
    </template>
    <template #text="">
      <state-settings-summary
        :entity="config.entity"
        :pending="pending"
        :state-setting="stateSetting"
      />
      <entity-dependencies-by-state-settings
        :entity="config.entity"
        :pending="pending"
        :state-setting="stateSetting"
        :color-indicator="config.colorIndicator"
      />
    </template>
  </modal-wrapper>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';

import EntityDependenciesByStateSettings from '@/components/other/entity/entity-dependencies-by-state-settings.vue';
import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions: mapEntityActions } = createNamespacedHelpers('entity');

export default {
  name: MODALS.entitiesRootCauseDiagram,
  components: {
    StateSettingsSummary,
    ModalWrapper,
    EntityDependenciesByStateSettings,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      pending: true,
      stateSetting: undefined,
    };
  },
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
  mounted() {
    this.fetchEntityStateSetting();
  },
  methods: {
    ...mapEntityActions({
      fetchEntityStateSettingWithoutStore: 'fetchStateSettingWithoutStore',
    }),

    async fetchEntityStateSetting() {
      this.pending = true;

      try {
        this.stateSetting = await this.fetchEntityStateSettingWithoutStore({ params: { _id: this.entity._id } });
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
