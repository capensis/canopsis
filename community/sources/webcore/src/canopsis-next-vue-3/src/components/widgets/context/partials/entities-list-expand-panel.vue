<template>
  <div class="entities-list-expand-panel">
    <v-tabs
      v-model="activeTab"
      class="entities-list-expand-panel__tabs"
      bg-color="secondary"
      color="primary"
      align-tabs="center"
    >
      <v-tab v-if="isAvailabilityEnabled" :value="tabs.availability">
        {{ $t('common.availability') }}
      </v-tab>
      <v-tab v-if="hasAccessToPbehavior" :value="tabs.pbehavior">
        {{ $t('common.pbehavior', 2) }}
      </v-tab>
      <v-tab v-if="isImpactDependsShown" :value="tabs.impactDepends">
        {{ $t('context.impactDepends') }}
      </v-tab>
      <v-tab :value="tabs.infos">
        {{ $t('common.infos') }}
      </v-tab>
      <v-tab v-if="hasWidgetCharts" :value="tabs.charts">
        {{ $t('context.charts') }}
      </v-tab>
      <v-tab v-if="isService" :value="tabs.treeOfDependencies">
        {{ $t('context.treeOfDependencies') }}
      </v-tab>
      <v-tab :value="tabs.impactChain">
        {{ $t('context.impactChain') }}
      </v-tab>
      <v-tab :value="tabs.activeAlarm">
        {{ $t('context.activeAlarm') }}
      </v-tab>
      <v-tab :value="tabs.resolvedAlarms">
        {{ $t('context.resolvedAlarms') }}
      </v-tab>
      <v-tab v-if="commentsShown" :value="tabs.comments">
        {{ $t('common.comment', 2) }}
      </v-tab>
    </v-tabs>

    <v-window
      v-model="activeTab"
      :transition="false"
      :reverse-transition="false"
      class="entities-list-expand-panel__content"
    >
      <c-window-expand-item
        v-if="isAvailabilityEnabled"
        :value="tabs.availability"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <entity-availability
          :entity="item"
          :default-time-range="availability?.default_time_range"
          :default-show-type="availability?.default_show_type"
        />
      </c-window-expand-item>

      <c-window-expand-item
        v-if="hasAccessToPbehavior"
        :value="tabs.pbehavior"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <pbehaviors-simple-list
          :entity="item"
          :removable="hasDeleteAnyPbehaviorAccess"
          :updatable="hasUpdateAnyPbehaviorAccess"
        />
      </c-window-expand-item>

      <v-window-item v-if="isImpactDependsShown" :value="tabs.impactDepends">
        <v-row>
          <v-col :cols="expandPanelCols" :offset="expandPanelOffset" class="pa-0">
            <impact-depends-tab :entity="item" />
          </v-col>
        </v-row>
      </v-window-item>

      <c-window-expand-item
        :value="tabs.infos"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <infos-tab :infos="item.infos" :columns-filters="columnsFilters" />
      </c-window-expand-item>

      <c-window-expand-item
        v-if="hasWidgetCharts"
        :value="tabs.charts"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <entity-charts
          :charts="charts"
          :entity="item"
          :available-metrics="item.filtered_perf_data"
        />
      </c-window-expand-item>

      <c-window-expand-item
        v-if="isService"
        :value="tabs.treeOfDependencies"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <tree-of-dependencies-tab
          :item="item"
          :columns="serviceDependenciesColumns"
          :type="treeOfDependenciesShowType"
        />
      </c-window-expand-item>

      <c-window-expand-item
        :value="tabs.impactChain"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <impact-chain-dependencies-tab
          :item="item"
          :columns="serviceDependenciesColumns"
        />
      </c-window-expand-item>

      <c-window-expand-item
        :value="tabs.activeAlarm"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <entity-alarms-list-table :entity="item" :columns="activeAlarmsColumns" />
      </c-window-expand-item>

      <c-window-expand-item
        :value="tabs.resolvedAlarms"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <entity-alarms-list-table
          :entity="item"
          :columns="resolvedAlarmsColumns"
          resolved
        />
      </c-window-expand-item>

      <c-window-expand-item
        v-if="commentsShown"
        :value="tabs.comments"
        :cols="expandPanelCols"
        :offset="expandPanelOffset"
      >
        <entity-comments :entity="item" addable editable />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { ENTITY_TYPES, GRID_SIZES, TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';
import { permissionsWidgetsContextPbehavior } from '@/mixins/permissions/widgets/context/pbehavior';
import { permissionsWidgetsEventComment } from '@/mixins/permissions/widgets/entity-comment';

import EntityAvailability from '@/components/other/entity/entity-availability.vue';
import EntityComments from '@/components/other/entity/entity-comments.vue';
import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/pbehaviors-simple-list.vue';
import EntityCharts from '@/components/widgets/chart/entity-charts.vue';

import EntityAlarmsListTable from './expand-panel-tabs/entity-alarms-list-table.vue';
import ImpactChainDependenciesTab from './expand-panel-tabs/impact-chain-dependencies-tab.vue';
import ImpactDependsTab from './expand-panel-tabs/impact-depends-tab.vue';
import InfosTab from './expand-panel-tabs/infos-tab.vue';
import TreeOfDependenciesTab from './expand-panel-tabs/tree-of-dependencies-tab.vue';

const CONTEXT_EXPAND_PANEL_TABS = {
  availability: 'availability',
  pbehavior: 'pbehavior',
  impactDepends: 'impactDepends',
  infos: 'infos',
  charts: 'charts',
  treeOfDependencies: 'treeOfDependencies',
  impactChain: 'impactChain',
  activeAlarm: 'activeAlarm',
  resolvedAlarms: 'resolvedAlarms',
  comments: 'comments',
};

export default {
  components: {
    EntityAlarmsListTable,
    EntityAvailability,
    EntityCharts,
    EntityComments,
    ImpactChainDependenciesTab,
    ImpactDependsTab,
    InfosTab,
    PbehaviorsSimpleList,
    TreeOfDependenciesTab,
  },

  mixins: [
    authMixin,
    permissionsWidgetsEventComment,
    permissionsWidgetsContextPbehavior,
    permissionsTechnicalExploitationPbehaviorMixin,
  ],

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
    expandGridRangeSize: {
      type: Array,
      default: () => [GRID_SIZES.min, GRID_SIZES.max],
    },
  },

  data() {
    return {
      activeTab: undefined,
      tabs: CONTEXT_EXPAND_PANEL_TABS,
    };
  },

  computed: {
    isService() {
      return this.item.type === ENTITY_TYPES.service;
    },

    isImpactDependsShown() {
      return this.item.type !== ENTITY_TYPES.service;
    },

    expandPanelOffset() {
      return this.expandGridRangeSize[0];
    },

    expandPanelCols() {
      const [start, end] = this.expandGridRangeSize;

      return end - start;
    },

    hasWidgetCharts() {
      return this.charts?.length && this.item.filtered_perf_data?.length;
    },

    isAvailabilityEnabled() {
      return this.availability?.enabled;
    },

    commentsShown() {
      return this.hasAccessToComment && [ENTITY_TYPES.service, ENTITY_TYPES.resource].includes(this.item.type);
    },

    availableTabs() {
      return [
        this.isAvailabilityEnabled && CONTEXT_EXPAND_PANEL_TABS.availability,
        this.hasAccessToPbehavior && CONTEXT_EXPAND_PANEL_TABS.pbehavior,
        this.isImpactDependsShown && CONTEXT_EXPAND_PANEL_TABS.impactDepends,
        CONTEXT_EXPAND_PANEL_TABS.infos,
        this.hasWidgetCharts && CONTEXT_EXPAND_PANEL_TABS.charts,
        this.isService && CONTEXT_EXPAND_PANEL_TABS.treeOfDependencies,
        CONTEXT_EXPAND_PANEL_TABS.impactChain,
        CONTEXT_EXPAND_PANEL_TABS.activeAlarm,
        CONTEXT_EXPAND_PANEL_TABS.resolvedAlarms,
        this.commentsShown && CONTEXT_EXPAND_PANEL_TABS.comments,
      ].filter(Boolean);
    },
  },

  watch: {
    availableTabs: {
      immediate: true,
      handler(availableTabs) {
        if (!availableTabs.includes(this.activeTab)) {
          [this.activeTab] = availableTabs;
        }
      },
    },
  },
};
</script>

<style lang="scss" scoped>
.entities-list-expand-panel {
  background: rgb(var(--v-theme-secondary));
}

.entities-list-expand-panel__content {
  background: rgb(var(--v-theme-secondary-lighten-1));
}
</style>
