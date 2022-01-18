<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t('context.pbehaviors') }}
    v-tab-item
      pbehaviors-list-tab(:item-id="item._id", :tab-id="tabId")

    template(v-if="item.type !== $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.impactDepends') }}
      v-tab-item(lazy)
        impact-depends-tab(:impact="item.impact", :depends="item.depends")

    v-tab {{ $t('common.infos') }}
    v-tab-item(lazy)
      infos-tab(:infos="item.infos", :columns-filters="columnsFilters")

    template(v-if="item.type === $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.treeOfDependencies') }}
      v-tab-item(lazy)
        tree-of-dependencies-tab(:item="item", :widget="widget")

    v-tab {{ $t('context.impactChain') }}
    v-tab-item(lazy)
      impact-chain-dependencies-tab(:item="item", :widget="widget")
</template>

<script>
import PbehaviorsListTab from './expand-panel-tabs/pbehaviors-tab.vue';
import ImpactDependsTab from './expand-panel-tabs/impact-depends-tab.vue';
import InfosTab from './expand-panel-tabs/infos-tab.vue';
import TreeOfDependenciesTab from './expand-panel-tabs/tree-of-dependencies-tab.vue';
import ImpactChainDependenciesTab from './expand-panel-tabs/impact-chain-dependencies-tab.vue';

export default {
  components: {
    ImpactChainDependenciesTab,
    PbehaviorsListTab,
    ImpactDependsTab,
    InfosTab,
    TreeOfDependenciesTab,
  },
  props: {
    item: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      required: true,
    },
    columnsFilters: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
