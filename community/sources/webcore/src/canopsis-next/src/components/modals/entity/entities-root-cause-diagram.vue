<template>
  <c-modal-wrapper close>
    <template #title>
      {{ title }}
    </template>
    <template #text>
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
  </c-modal-wrapper>
</template>

<script>
import { mapActions } from 'pinia';

import { MODALS } from '@/constants';

import { useEntityStore } from '@/stores/modules/entity';

import EntityDependenciesByStateSettings from '@/components/other/entity/entity-dependencies-by-state-settings.vue';
import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';

export default {
  name: MODALS.entitiesRootCauseDiagram,

  components: {
    EntityDependenciesByStateSettings,
    StateSettingsSummary,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      pending: true,
      stateSetting: undefined,
    };
  },

  computed: {
    config() {
      return this.modal.config ?? {};
    },

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
    ...mapActions(useEntityStore, {
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
