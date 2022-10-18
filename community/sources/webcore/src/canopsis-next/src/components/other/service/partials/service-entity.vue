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
            @select="$listeners.select",
            @remove-unavailable="$listeners['remove-unavailable']"
          )
        v-card(color="white black--text")
          v-card-text
            entity-info-tab(
              v-if="!isService && !hasAccessToPbehaviors",
              :entity="entity",
              :template="template",
              @add:action="$listeners['add:action']",
              @refresh="$listeners.refresh"
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
                  @add:action="$listeners['add:action']",
                  @refresh="$listeners.refresh"
                )
              template(v-if="isService")
                v-tab {{ $t('modals.service.entity.tabs.treeOfDependencies') }}
                v-tab-item(lazy)
                  entity-tree-of-dependencies-tab(
                    :entity="entity",
                    :columns="serviceDependenciesColumns"
                  )
              template(v-if="hasAccessToPbehaviors")
                v-tab {{ $tc('common.activePbehavior') }}
                v-tab-item(lazy)
                  pbehaviors-simple-list(:entity="entity", editable, deletable, dense, with-active-status)
</template>

<script>
import { isNull } from 'lodash';

import { ENTITY_TYPES, USERS_PERMISSIONS } from '@/constants';

import { getEntityColor } from '@/helpers/color';

import { authMixin } from '@/mixins/auth';
import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';

import PbehaviorsSimpleList from '@/components/other/pbehavior/partials/pbehaviors-simple-list.vue';

import EntityHeader from './service-entity-header.vue';
import EntityInfoTab from './service-entity-info-tab.vue';
import EntityTreeOfDependenciesTab from './service-entity-tree-of-dependencies-tab.vue';

export default {
  inject: ['$actionsQueue'],
  components: {
    PbehaviorsSimpleList,
    EntityHeader,
    EntityInfoTab,
    EntityTreeOfDependenciesTab,
  },
  mixins: [authMixin, vuetifyTabsMixin],
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

    isService() {
      return this.entity.source_type === ENTITY_TYPES.service;
    },

    color() {
      return getEntityColor(this.entity, this.colorIndicator);
    },

    hasPbehaviors() {
      return !!this.entity.pbehaviors.length;
    },

    hasAccessToPbehaviors() {
      return this.hasPbehaviors
        && this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors);
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
    padding: 0 16px;
    height: auto;
  }
</style>
