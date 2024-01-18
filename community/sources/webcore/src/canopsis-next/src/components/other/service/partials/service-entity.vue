<template>
  <div class="weather-service-entity-expansion-panel">
    <v-expansion-panels v-model="opened">
      <v-expansion-panel>
        <v-expansion-panel-header :color="color">
          <template #actions="">
            <v-icon color="white">
              keyboard_arrow_down
            </v-icon>
          </template>
          <service-entity-header
            :selected="selected"
            :selectable="!!availableActions.length"
            :entity="entity"
            :entity-name-field="entityNameField"
            :last-action-unavailable="lastActionUnavailable"
            @update:selected="$listeners['update:selected']"
            @remove:unavailable="$listeners['remove:unavailable']"
          />
        </v-expansion-panel-header>
        <v-expansion-panel-content>
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
              <v-tabs
                v-else
                ref="tabs"
                v-model="activeTab"
                slider-color="primary"
                centered
              >
                <v-tab>{{ $t('modals.service.entity.tabs.info') }}</v-tab>
                <v-tab-item>
                  <service-entity-info-tab
                    :entity="entity"
                    :template="template"
                    :actions="availableActions"
                    @apply="applyActionForEntity"
                    @execute="executeAlarmInstruction"
                  />
                </v-tab-item>
                <template v-if="isService">
                  <v-tab>{{ $t('modals.service.entity.tabs.treeOfDependencies') }}</v-tab>
                  <v-tab-item>
                    <service-entity-tree-of-dependencies-tab
                      class="pa-2"
                      :entity="entity"
                      :columns="serviceDependenciesColumns"
                      :type="treeOfDependenciesShowType"
                    />
                  </v-tab-item>
                </template>
                <template v-if="hasAccessToPbehaviors">
                  <v-tab>{{ $tc('common.pbehavior', 2) }}</v-tab>
                  <v-tab-item>
                    <pbehaviors-simple-list
                      :entity="entity"
                      removable
                      updatable
                      dense
                      with-active-status
                    />
                  </v-tab-item>
                </template>
              </v-tabs>
            </v-card-text>
          </v-card>
        </v-expansion-panel-content>
      </v-expansion-panel>
    </v-expansion-panels>
  </div>
</template>

<script>
import { isNull } from 'lodash';

import { ENTITY_TYPES, MODALS, TREE_OF_DEPENDENCIES_SHOW_TYPES, USERS_PERMISSIONS } from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import { getAvailableActionsByEntity } from '@/helpers/entities/entity/actions';
import { isInstructionManual } from '@/helpers/entities/remediation/instruction/form';

import { authMixin } from '@/mixins/auth';
import { vuetifyTabsMixin } from '@/mixins/vuetify/tabs';
import { widgetActionPanelServiceEntityMixin } from '@/mixins/widget/actions-panel/service-entity';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';

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
  },
  data() {
    return {
      opened: [0],
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

    treeOfDependenciesShowType() {
      return this.widgetParameters.treeOfDependenciesShowType ?? TREE_OF_DEPENDENCIES_SHOW_TYPES.custom;
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
  & ::v-deep .v-expansion-panel-header {
    padding: 0 16px;
  }

  & ::v-deep .v-expansion-panel-header__icon .v-icon,
  & ::v-deep .v-expansion-panel__header .v-input .v-icon {
    color: white !important;
  }
}
</style>
