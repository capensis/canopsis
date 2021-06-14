<template lang="pug">
  c-treeview-data-table.service-dependencies(
    :items="items",
    :headers="headers",
    :loading="hasActivePending",
    :load-children="loadChildren",
    :dark="dark",
    :light="light",
    item-key="key"
  )
    template(slot="expand", slot-scope="props")
      v-btn(
        v-if="props.item.entity",
        :color="$config.COLORS.impactState[props.item.impact_state]",
        icon,
        dark,
        @click="showTreeOfDependenciesModal(props.item)"
      )
        v-icon {{ props.item.entity | btnIcon }}
      v-tooltip(v-else, right)
        v-btn(
          slot="activator",
          :loading="pendingByIds[props.item.parentId]",
          icon,
          @click="loadMore(props.item.parentId)"
        )
          v-icon more_horiz
        span {{ $t('common.loadMore') }}
    template(
      slot="expand-append",
      slot-scope="props",
      v-if="includeRoot && isInRootIds(props.item._id)"
    )
      div.expand-append
        v-icon arrow_right_alt
        v-chip.ma-0(
          :color="$config.COLORS.impactState[props.item.impact_state]",
          text-color="white"
        )
          span.px-2.body-2.font-weight-bold {{ props.item.impact_state }}
    tr(slot="items", slot-scope="props")
      td(v-for="header in headers", :key="header.value")
        color-indicator-wrapper(
          v-if="props.item.entity",
          :entity="props.item.entity",
          :alarm="props.item.alarm",
          :type="header.colorIndicator"
        ) {{ props.item | get(header.value) }}
</template>

<script>
import { get, uniq } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, ENTITY_TYPES, DEFAULT_SERVICE_DEPENDENCIES_COLUMNS } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';
import {
  dependencyToTreeviewDependency,
  treeviewDependencyToDependency,
  normalizeDependencies,
  getLoadMoreDenormalizedChild,
  dependenciesDenormalize,
} from '@/helpers/treeview/service-dependencies';

import ColorIndicatorWrapper from '@/components/common/table/color-indicator-wrapper.vue';

const { mapActions } = createNamespacedHelpers('service');

export default {
  filters: {
    btnIcon(entity) {
      return entity.type === ENTITY_TYPES.service ? '$vuetify.icons.engineering' : 'person';
    },
  },
  components: { ColorIndicatorWrapper },
  props: {
    root: {
      type: Object,
      required: true,
    },
    includeRoot: {
      type: Boolean,
      default: false,
    },
    columns: {
      type: Array,
      required: false,
    },
    dark: {
      type: Boolean,
      default: false,
    },
    light: {
      type: Boolean,
      default: false,
    },
    impact: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const dependenciesByIds = {};
    const rootIds = [];

    if (this.includeRoot) {
      const treeviewRoot = dependencyToTreeviewDependency(this.root);

      dependenciesByIds[treeviewRoot._id] = treeviewRoot;
      rootIds.push(treeviewRoot._id);
    }

    return {
      rootIds,
      dependenciesByIds,

      metaByIds: {},
      pendingByIds: {},
    };
  },
  computed: {
    rootId() {
      return this.root.entity._id;
    },

    rootHasNextPage() {
      const meta = this.metaByIds[this.rootId] || {};

      return meta.page < meta.page_count;
    },

    treeviewRoot() {
      return dependencyToTreeviewDependency(this.root);
    },

    hasActivePending() {
      return Object.values(this.pendingByIds).some(pending => pending);
    },

    headers() {
      const columns = this.columns || defaultColumnsToColumns(DEFAULT_SERVICE_DEPENDENCIES_COLUMNS);

      return columns.map(({ label, value, colorIndicator }) => ({
        colorIndicator,

        sortable: false,
        text: label,
        value: value.match(/entity.|alarm./) ? value : `entity.${value}`,
      }));
    },

    items() {
      const items = dependenciesDenormalize(this.rootIds, this.dependenciesByIds, this.metaByIds);

      if (!this.includeRoot && this.rootHasNextPage) {
        items.push(getLoadMoreDenormalizedChild(this.treeviewRoot));
      }

      return items;
    },
  },
  async mounted() {
    const ids = await this.fetchDependenciesById(this.rootId);

    if (!this.includeRoot) {
      this.rootIds = ids;
    }
  },
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    isInRootIds(id) {
      return this.rootIds.includes(id);
    },

    showTreeOfDependenciesModal(dependency) {
      const isServiceEntity = dependency.entity.type === ENTITY_TYPES.service;

      if (!this.impact && !isServiceEntity) {
        return;
      }

      this.$modals.show({
        name: MODALS.serviceDependencies,
        config: {
          root: treeviewDependencyToDependency(dependency),
          columns: this.columns,
          impact: this.impact,
        },
      });
    },

    async loadMore(id) {
      const isRoot = this.rootId === id;
      const meta = this.metaByIds[id] || {};
      const params = {
        page: meta.page + 1,
        limit: PAGINATION_LIMIT,
      };

      const ids = await this.fetchDependenciesById(id, params);

      if (!this.includeRoot && isRoot) {
        this.rootIds.push(...ids);
      }
    },

    loadChildren(dependency) {
      return this.fetchDependenciesById(dependency._id);
    },

    async fetchDependenciesById(id, params = { limit: PAGINATION_LIMIT }) {
      this.$set(this.pendingByIds, id, true);

      const { data, meta } = await this.fetchDependenciesList({ id, params });
      const { dependencies, result } = normalizeDependencies(data);

      Object.entries(dependencies).forEach(([dependencyId, dependency]) => {
        this.$set(this.dependenciesByIds, dependencyId, dependency);
      });

      if (this.dependenciesByIds[id] && result.length) {
        const children = uniq([...get(this.dependenciesByIds[id], 'children', []), ...result]);

        this.$set(this.dependenciesByIds[id], 'children', children);
      }

      this.$set(this.metaByIds, id, meta);
      this.$set(this.pendingByIds, id, false);

      return result;
    },

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
</script>

<style lang="scss" scoped>
.service-dependencies /deep/ .v-treeview-node__label {
  &, .expand-append {
    display: inline-flex;
    align-items: center;

    .v-chip__content {
      height: 20px;
    }
  }
}

</style>
