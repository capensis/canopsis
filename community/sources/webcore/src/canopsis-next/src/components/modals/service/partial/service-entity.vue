<template lang="pug">
  div.weather-service-entity-expansion-panel
    v-expansion-panel(v-model="opened", dark)
      v-expansion-panel-content(:style="{ backgroundColor: color }")
        template(#header="")
          entity-header(
            :selected="selected",
            :entity="entity",
            :entity-name-field="entityNameField",
            :last-action-unavailable="lastActionUnavailable",
            :active="hasActivePbehavior",
            :paused="hasPausedPbehavior",
            @select="$listeners.select",
            @remove-unavailable="$listeners['remove-unavailable']"
          )
        v-card(color="white black--text")
          v-card-text
            entity-info-tab(
              v-if="!isService",
              :entity="entity",
              :template="template",
              :active="hasActivePbehavior",
              :paused="hasPausedPbehavior",
              @add:action="$listeners['add:action']"
            )
            v-tabs(
              v-else,
              ref="tabs",
              v-model="activeTab",
              slider-color="primary",
              fixed-tabs,
              light
            )
              v-tab {{ $t('modals.service.entity.tabs.info') }}
              v-tab-item
                entity-info-tab(
                  :entity="entity",
                  :template="template",
                  :active="hasActivePbehavior",
                  :paused="hasPausedPbehavior",
                  @add:action="$listeners['add:action']"
                )
              v-tab {{ $t('modals.service.entity.tabs.treeOfDependencies') }}
              v-tab-item(lazy)
                entity-tree-of-dependencies-tab(
                  :entity="entity",
                  :columns="serviceDependenciesColumns"
                )
</template>

<script>
import { isNull } from 'lodash';

import { ENTITY_TYPES } from '@/constants';

import { getEntityColor } from '@/helpers/color';
import { hasActivePbehavior, hasPausedPbehavior } from '@/helpers/entities/pbehavior';

import vuetifyTabsMixin from '@/mixins/vuetify/tabs';

import EntityHeader from './service-entity-header.vue';
import EntityInfoTab from './service-entity-info-tab.vue';
import EntityTreeOfDependenciesTab from './service-entity-tree-of-dependencies-tab.vue';

export default {
  inject: ['$actionsQueue'],
  components: {
    EntityHeader,
    EntityInfoTab,
    EntityTreeOfDependenciesTab,
  },
  mixins: [vuetifyTabsMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    lastActionUnavailable: {
      type: Boolean,
      default: false,
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    widgetParameters: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      opened: [false],
      activeTab: 0,
    };
  },
  computed: {
    template() {
      return this.widgetParameters.entityTemplate || '';
    },

    serviceDependenciesColumns() {
      return this.widgetParameters.serviceDependenciesColumns;
    },

    colorIndicator() {
      return this.widgetParameters.colorIndicator;
    },

    hasActivePbehavior() {
      return hasActivePbehavior(this.entity.pbehaviors);
    },

    hasPausedPbehavior() {
      return hasPausedPbehavior(this.entity.pbehaviors);
    },

    isService() {
      return this.entity.source_type === ENTITY_TYPES.service;
    },

    color() {
      return getEntityColor(this.entity, this.colorIndicator);
    },
  },
  watch: {
    opened(opened) {
      if (!isNull(opened) && this.$refs.tabs) {
        this.$nextTick(this.callTabsUpdateTabsMethod);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .weather-service-entity-expansion-panel /deep/ .v-expansion-panel__header {
    padding: 0 12px;
    height: auto;
  }
</style>
