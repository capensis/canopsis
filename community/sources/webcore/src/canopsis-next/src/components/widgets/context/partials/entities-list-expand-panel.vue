<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.pbehavior', 2) }}
    v-tab-item
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card
            v-card-text
              pbehaviors-simple-list(
                :entity="item",
                :removable="hasDeleteAnyPbehaviorAccess",
                :updatable="hasUpdateAnyPbehaviorAccess"
              )

    template(v-if="item.type !== $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.impactDepends') }}
      v-tab-item(lazy)
        v-flex(:class="cardFlexClass")
          impact-depends-tab(:entity="item")

    v-tab {{ $t('common.infos') }}
    v-tab-item(lazy)
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card
            v-card-text
              infos-tab(:infos="item.infos", :columns-filters="columnsFilters")

    template(v-if="item.type === $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.treeOfDependencies') }}
      v-tab-item(lazy)
        v-layout.pa-3(row)
          v-flex(:class="cardFlexClass")
            v-card
              v-card-text
                tree-of-dependencies-tab(:item="item", :columns="serviceDependenciesColumns")

    v-tab {{ $t('context.impactChain') }}
    v-tab-item(lazy)
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card
            v-card-text.pa-0
              impact-chain-dependencies-tab(:item="item", :columns="serviceDependenciesColumns")

    v-tab {{ $t('context.activeAlarm') }}
    v-tab-item(lazy)
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card
            v-card-text
              entity-alarms-list-table(:entity="item", :columns="activeAlarmsColumns")

    v-tab {{ $t('context.resolvedAlarms') }}
    v-tab-item(lazy)
      v-layout.pa-3(row)
        v-flex(:class="cardFlexClass")
          v-card
            v-card-text
              entity-alarms-list-table(:entity="item", :columns="resolvedAlarmsColumns", resolved)
</template>

<script>
import { GRID_SIZES } from '@/constants';

import { getFlexClassesForGridRangeSize } from '@/helpers/grid';

import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import PbehaviorsSimpleList from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-simple-list.vue';

import ImpactDependsTab from './expand-panel-tabs/impact-depends-tab.vue';
import InfosTab from './expand-panel-tabs/infos-tab.vue';
import TreeOfDependenciesTab from './expand-panel-tabs/tree-of-dependencies-tab.vue';
import ImpactChainDependenciesTab from './expand-panel-tabs/impact-chain-dependencies-tab.vue';
import EntityAlarmsListTable from './expand-panel-tabs/entity-alarms-list-table.vue';

export default {
  components: {
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
    expandGridRangeSize: {
      type: Array,
      default: () => [GRID_SIZES.min, GRID_SIZES.max],
    },
  },
  computed: {
    cardFlexClass() {
      return getFlexClassesForGridRangeSize(this.expandGridRangeSize);
    },
  },
};
</script>
