<template>
  <v-tabs
    background-color="secondary lighten-1"
    slider-color="primary"
    dark
    centered
  >
    <template v-if="isAvailabilityEnabled">
      <v-tab>{{ $t('common.availability') }}</v-tab>
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <entity-availability
                  :entity="item"
                  :default-time-range="availability?.default_time_range"
                  :default-show-type="availability?.default_show_type"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </template>

    <v-tab>{{ $tc('common.pbehavior', 2) }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <pbehaviors-simple-list
                :entity="item"
                :removable="hasDeleteAnyPbehaviorAccess"
                :updatable="hasUpdateAnyPbehaviorAccess"
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>

    <template v-if="item.type !== $constants.ENTITY_TYPES.service">
      <v-tab>{{ $t('context.impactDepends') }}</v-tab>
      <v-tab-item>
        <impact-depends-tab :entity="item" />
      </v-tab-item>
    </template>

    <v-tab>{{ $t('common.infos') }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <infos-tab
                :infos="item.infos"
                :columns-filters="columnsFilters"
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>

    <template v-if="hasWidgetCharts">
      <v-tab>{{ $t('context.charts') }}</v-tab>
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <entity-charts
                  :charts="charts"
                  :entity="item"
                  :available-metrics="item.filtered_perf_data"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </template>

    <template v-if="item.type === $constants.ENTITY_TYPES.service">
      <v-tab>{{ $t('context.treeOfDependencies') }}</v-tab>
      <v-tab-item>
        <v-layout class="pa-3">
          <v-flex xs12>
            <v-card>
              <v-card-text>
                <tree-of-dependencies-tab
                  :item="item"
                  :columns="serviceDependenciesColumns"
                  :type="treeOfDependenciesShowType"
                />
              </v-card-text>
            </v-card>
          </v-flex>
        </v-layout>
      </v-tab-item>
    </template>

    <v-tab>{{ $t('context.impactChain') }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text class="pa-0">
              <impact-chain-dependencies-tab
                :item="item"
                :columns="serviceDependenciesColumns"
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>

    <v-tab>{{ $t('context.activeAlarm') }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <entity-alarms-list-table
                :entity="item"
                :columns="activeAlarmsColumns"
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>

    <v-tab>{{ $t('context.resolvedAlarms') }}</v-tab>
    <v-tab-item>
      <v-layout class="pa-3">
        <v-flex xs12>
          <v-card>
            <v-card-text>
              <entity-alarms-list-table
                :entity="item"
                :columns="resolvedAlarmsColumns"
                resolved
              />
            </v-card-text>
          </v-card>
        </v-flex>
      </v-layout>
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import EntityCharts from '@/components/widgets/chart/entity-charts.vue';
import EntityAvailability from '@/components/other/entity/entity-availability.vue';

import ImpactDependsTab from './expand-panel-tabs/impact-depends-tab.vue';
import InfosTab from './expand-panel-tabs/infos-tab.vue';
import TreeOfDependenciesTab from './expand-panel-tabs/tree-of-dependencies-tab.vue';
import ImpactChainDependenciesTab from './expand-panel-tabs/impact-chain-dependencies-tab.vue';
import EntityAlarmsListTable from './expand-panel-tabs/entity-alarms-list-table.vue';

export default {
  components: {
    EntityAvailability,
    EntityCharts,
    PbehaviorsSimpleList,
    ImpactChainDependenciesTab,
    ImpactDependsTab,
    InfosTab,
    TreeOfDependenciesTab,
    EntityAlarmsListTable,
  },
  mixins: [permissionsTechnicalExploitationPbehaviorMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
    resolvedAlarmsColumns: {
      type: Array,
      required: true,
    },
    activeAlarmsColumns: {
      type: Array,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
    },
    serviceDependenciesColumns: {
      type: Array,
      default: () => [],
    },
    charts: {
      type: Array,
      default: () => [],
    },
    treeOfDependenciesShowType: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
    },
    availability: {
      type: Object,
      required: false,
    },
  },
  computed: {
    hasWidgetCharts() {
      return this.charts?.length && this.item.filtered_perf_data?.length;
    },

    isAvailabilityEnabled() {
      return this.availability?.enabled;
    },
  },
};
</script>
