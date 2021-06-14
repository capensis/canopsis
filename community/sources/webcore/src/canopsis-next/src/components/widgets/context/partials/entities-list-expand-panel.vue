<template lang="pug">
  v-tabs.visible(v-model="activeTab", color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $t('context.expandPanel.tabs.pbehaviors') }}
    v-tab-item
      pbehaviors-list-tab(:item-id="item._id", :tab-id="tabId")
    template(v-if="item.type !== $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.expandPanel.tabs.impactDepends') }}
      v-tab-item
        impact-depends-tab(:impact="item.impact", :depends="item.depends")
    v-tab {{ $t('context.expandPanel.tabs.infos') }}
    v-tab-item
      infos-tab(:infos="item.infos", :columns-filters="columnsFilters")
    template(v-if="item.type === $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.expandPanel.tabs.treeOfDependencies') }}
      v-tab-item
        tree-of-dependencies-tab(:item="item", :widget="widget")
    template
      v-tab {{ $t('context.expandPanel.tabs.impactChain') }}
      v-tab-item
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
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-tabs.visible {
  & /deep/ > .v-tabs__bar {
    display: block;
  }
}
</style>
