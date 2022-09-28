<template lang="pug">
  v-tabs(color="secondary lighten-1", slider-color="primary", dark, centered)
    v-tab {{ $tc('common.pbehavior', 2) }}
    v-tab-item
      v-card.secondary.lighten-2(flat)
        v-card-text
          pbehaviors-simple-list(
            :entity="item",
            :deletable="hasDeleteAnyPbehaviorAccess",
            :editable="hasUpdateAnyPbehaviorAccess"
          )

    template(v-if="item.type !== $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.impactDepends') }}
      v-tab-item(lazy)
        impact-depends-tab(:entity="item")

    v-tab {{ $t('common.infos') }}
    v-tab-item(lazy)
      infos-tab(:infos="item.infos", :columns-filters="columnsFilters")

    template(v-if="item.type === $constants.ENTITY_TYPES.service")
      v-tab {{ $t('context.treeOfDependencies') }}
      v-tab-item(lazy)
        tree-of-dependencies-tab(:item="item", :columns="serviceDependenciesColumns")

    v-tab {{ $t('context.impactChain') }}
    v-tab-item(lazy)
      impact-chain-dependencies-tab(:item="item", :columns="serviceDependenciesColumns")
</template>

<script>
import { permissionsTechnicalExploitationPbehaviorMixin } from '@/mixins/permissions/technical/exploitation/pbehavior';

import PbehaviorsSimpleList from '@/components/other/pbehavior/partials/pbehaviors-simple-list.vue';

import ImpactDependsTab from './expand-panel-tabs/impact-depends-tab.vue';
import InfosTab from './expand-panel-tabs/infos-tab.vue';
import TreeOfDependenciesTab from './expand-panel-tabs/tree-of-dependencies-tab.vue';
import ImpactChainDependenciesTab from './expand-panel-tabs/impact-chain-dependencies-tab.vue';

export default {
  components: {
    PbehaviorsSimpleList,
    ImpactChainDependenciesTab,
    ImpactDependsTab,
    InfosTab,
    TreeOfDependenciesTab,
  },
  mixins: [permissionsTechnicalExploitationPbehaviorMixin],
  props: {
    item: {
      type: Object,
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
  },
};
</script>
