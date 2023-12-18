<template lang="pug">
  div.weather-service-entity-expansion-panel
    v-expansion-panel(v-model="opened")
      v-expansion-panel-content(:style="{ backgroundColor: color }")
        template(#header="")
          service-entity-header(
            :selected="selected",
            :selectable="!!availableActions.length",
            :entity="entity",
            :entity-name-field="entityNameField",
            :last-action-unavailable="lastActionUnavailable",
            @update:selected="$listeners['update:selected']",
            @remove:unavailable="$listeners['remove:unavailable']"
          )
        v-card
          v-card-text
            service-entity-info-tab(
              v-if="!isService && !hasAccessToPbehaviors",
              :entity="entity",
              :template="template",
              :actions="availableActions",
              @apply="applyActionForEntity",
              @execute="executeAlarmInstruction"
            )
            v-tabs(
              v-else,
              ref="tabs",
              v-model="activeTab",
              slider-color="primary",
              fixed-tabs
            )
              v-tab {{ $t('modals.service.entity.tabs.info') }}
              v-tab-item
                service-entity-info-tab(
                  :entity="entity",
                  :template="template",
                  :actions="availableActions",
                  @apply="applyActionForEntity",
                  @execute="executeAlarmInstruction"
                )
              template(v-if="isService")
                v-tab {{ $t('modals.service.entity.tabs.treeOfDependencies') }}
                v-tab-item(lazy)
                  service-entity-tree-of-dependencies-tab.pa-2(
                    :entity="entity",
                    :columns="serviceDependenciesColumns"
                  )
              template(v-if="hasAccessToPbehaviors")
                v-tab {{ $tc('common.pbehavior', 2) }}
                v-tab-item(lazy)
                  pbehaviors-simple-list(
                    :entity="entity",
                    removable,
                    updatable,
                    dense,
                    with-active-status
                  )
</template>

<script>
import { isNull } from 'lodash';

import { ENTITY_TYPES, MODALS, USERS_PERMISSIONS } from '@/constants';

import { getEntityColor } from '@/helpers/color';
import { getAvailableActionsByEntity, isDisabledActionForEntityByActionsRequests } from '@/helpers/entities/entity';
import { isInstructionManual } from '@/helpers/forms/remediation-instruction';

import { authMixin } from '@/mixins/auth';
import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';
import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';

import ServiceEntityHeader from './service-entity-header.vue';
import ServiceEntityInfoTab from './service-entity-info-tab.vue';
import ServiceEntityTreeOfDependenciesTab from './service-entity-tree-of-dependencies-tab.vue';

export default {
  components: {
    PbehaviorsSimpleList,
    ServiceEntityHeader,
    ServiceEntityInfoTab,
    ServiceEntityTreeOfDependenciesTab,
  },
  mixins: [
    authMixin,
    vuetifyTabsMixin,
    widgetActionPanelServiceEntityMixin,
  ],
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
    actionsRequests: {
      type: Array,
      default: () => [],
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

    hasAccessToPbehaviors() {
      return this.entity.pbehaviors?.length
        && this.checkAccess(USERS_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors);
    },

    allActions() {
      return getAvailableActionsByEntity(this.entity);
    },

    availableActions() {
      return this.allActions
        .filter(this.actionsAccessFilterHandler)
        .map(action => ({
          ...action,
          loading: this.pendingByActionType[action.type],
          disabled: this.pendingByActionType[action.type]
              || isDisabledActionForEntityByActionsRequests(this.entity._id, action.type, this.actionsRequests),
        }));
    },
  },
  watch: {
    opened(opened) {
      if (!isNull(opened) && this.$refs.tabs) {
        this.$nextTick(this.callTabsUpdateTabsMethod);
      }
    },
  },
  methods: {
    applyActionForEntity({ type }) {
      this.applyEntityAction(type, [this.entity]);
    },

    executeAlarmInstruction(assignedInstruction) {
      const refreshEntities = () => this.$emit('refresh');

      this.$modals.show({
        name: isInstructionManual(assignedInstruction)
          ? MODALS.executeRemediationInstruction
          : MODALS.executeRemediationSimpleInstruction,
        config: {
          assignedInstruction,
          alarmId: this.entity.alarm_id,
          onClose: refreshEntities,
          onComplete: refreshEntities,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.weather-service-entity-expansion-panel {
  & ::v-deep .v-expansion-panel__header {
    padding: 0 16px;
    height: auto;
  }

  & ::v-deep .v-expansion-panel__header__icon .v-icon,
  & ::v-deep .v-expansion-panel__header .v-input .v-icon {
    color: white !important;
  }
}
</style>
