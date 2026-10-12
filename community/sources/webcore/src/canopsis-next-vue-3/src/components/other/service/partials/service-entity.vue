<template>
  <div class="weather-service-entity-expansion-panel">
    <v-expansion-panels v-model="opened">
      <v-expansion-panel>
        <v-expansion-panel-title :class="entityColorClass">
          <service-entity-header
            :selected="selected"
            :selectable="!!availableActions.length"
            :entity="entity"
            :entity-name-field="entityNameField"
            :last-action-unavailable="lastActionUnavailable"
            @update:selected="$emit('update:selected', $event)"
            @remove:unavailable="$emit('remove:unavailable')"
          />
        </v-expansion-panel-title>
        <v-expansion-panel-text eager>
          <v-card>
            <v-card-text>
              <service-entity-info-tab
                v-if="!isService && !hasAccessToPbehaviors"
                :entity="entity"
                :template="template"
                :actions="availableActions"
                @apply="applyActionForEntity"
                @execute="executeAlarmInstruction"
              />
              <template v-else>
                <v-tabs
                  v-model="activeTab"
                  align-tabs="center"
                  color="primary"
                >
                  <v-tab :value="tabs.info">{{ $t('modals.service.entity.tabs.info') }}</v-tab>
                  <v-tab v-if="isService" :value="tabs.treeOfDependencies">
                    {{ $t('modals.service.entity.tabs.treeOfDependencies') }}
                  </v-tab>
                  <v-tab v-if="hasAccessToPbehaviors" :value="tabs.pbehavior">
                    {{ $t('common.pbehavior', 2) }}
                  </v-tab>
                </v-tabs>
                <v-window v-model="activeTab" class="mt-2">
                  <v-window-item :value="tabs.info">
                    <service-entity-info-tab
                      :entity="entity"
                      :template="template"
                      :actions="availableActions"
                      @apply="applyActionForEntity"
                      @execute="executeAlarmInstruction"
                    />
                  </v-window-item>
                  <v-window-item v-if="isService" :value="tabs.treeOfDependencies">
                    <service-entity-tree-of-dependencies-tab
                      :entity="entity"
                      :columns="serviceDependenciesColumns"
                      :type="treeOfDependenciesShowType"
                      class="pa-2"
                    />
                  </v-window-item>
                  <v-window-item v-if="hasAccessToPbehaviors" :value="tabs.pbehavior">
                    <pbehaviors-simple-list
                      :entity="entity"
                      removable
                      updatable
                      dense
                      with-active-status
                    />
                  </v-window-item>
                </v-window>
              </template>
            </v-card-text>
          </v-card>
        </v-expansion-panel-text>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import { ENTITY_TYPES, MODALS, TREE_OF_DEPENDENCIES_SHOW_TYPES, USER_PERMISSIONS } from '@/constants';

import {
  getAvailableActionsByEntity,
  isDisabledActionForEntityByActionsRequests,
} from '@/helpers/entities/entity/actions';
import { getEntityColorClass } from '@/helpers/entities/entity/color';
import { isInstructionTypeManual } from '@/helpers/entities/remediation/instruction/form';

import { authMixin } from '@/mixins/auth';
import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';

import ServiceEntityHeader from './service-entity-header.vue';
import ServiceEntityInfoTab from './service-entity-info-tab.vue';
import ServiceEntityTreeOfDependenciesTab from './service-entity-tree-of-dependencies-tab.vue';

const tabs = {
  info: 'info',
  treeOfDependencies: 'treeOfDependencies',
  pbehavior: 'pbehavior',
};

export default {
  components: {
    PbehaviorsSimpleList,
    ServiceEntityHeader,
    ServiceEntityInfoTab,
    ServiceEntityTreeOfDependenciesTab,
  },
  mixins: [
    authMixin,
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
  emits: ['add:action', 'refresh', 'update:selected', 'remove:unavailable'],
  data() {
    return {
      opened: undefined,
      activeTab: tabs.info,
      tabs,
    };
  },
  computed: {
    template() {
      return this.widgetParameters.entityTemplate || '';
    },

    serviceDependenciesColumns() {
      return this.widgetParameters.serviceDependenciesColumns;
    },

    treeOfDependenciesShowType() {
      return this.widgetParameters.treeOfDependenciesShowType ?? TREE_OF_DEPENDENCIES_SHOW_TYPES.custom;
    },

    colorIndicator() {
      return this.widgetParameters.colorIndicator;
    },

    isService() {
      return this.entity.source_type === ENTITY_TYPES.service;
    },

    entityColorClass() {
      return getEntityColorClass(this.entity, this.colorIndicator);
    },

    hasAccessToPbehaviors() {
      return this.entity.pbehaviors?.length
        && this.checkAccess(USER_PERMISSIONS.business.serviceWeather.actions.entityManagePbehaviors);
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
  methods: {
    applyActionForEntity({ type }) {
      this.applyEntityAction(type, [this.entity]);
    },

    executeAlarmInstruction(assignedInstruction) {
      const refreshEntities = () => this.$emit('refresh', true);

      this.$modals.show({
        name: isInstructionTypeManual(assignedInstruction?.type)
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
  :deep(.v-expansion-panel-title) {
    padding: 0 16px;
  }

  :deep(.v-expansion-panel-title__icon .v-icon),
  :deep(.v-expansion-panel-title .v-input .v-icon) {
    color: inherit !important;
  }
}
</style>
