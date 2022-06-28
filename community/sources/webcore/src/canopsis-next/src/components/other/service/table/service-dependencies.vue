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
    template(#expand="{ item }")
      v-btn(
        v-if="item.entity",
        :color="getEntityColor(item)",
        icon,
        dark,
        @click="showTreeOfDependenciesModal(item)"
      )
        v-icon {{ getIconByEntity(item.entity) }}
      v-tooltip(v-else, right)
        v-btn(
          slot="activator",
          :loading="pendingByIds[item.parentId]",
          icon,
          @click="loadMore(item.parentId)"
        )
          v-icon more_horiz
        span {{ $t('common.loadMore') }}
      v-tooltip(v-if="item.cycle", top)
        template(#activator="{ on }")
          v-icon(v-on="on", color="error", size="14") autorenew
        span {{ $t('common.cycleDependency') }}
    template(
      slot="expand-append",
      slot-scope="{ item }",
      v-if="includeRoot && isInRootIds(item._id)"
    )
      div.expand-append
        v-icon arrow_right_alt
        v-chip.ma-0(
          :color="getEntityColor(item)",
          text-color="white"
        )
          span.px-2.body-2.font-weight-bold {{ item.impact_state }}
    template(#items="{ item }")
      tr
        td(v-for="(header, index) in headers", :key="header.value")
          c-no-events-icon(v-if="!index", :value="item.entity | get('idle_since')", top)
          color-indicator-wrapper(
            v-else-if="item.entity",
            :entity="item.entity",
            :alarm="item.alarm",
            :type="header.colorIndicator"
          ) {{ item | get(header.value) }}
</template>

<script>
import { get, uniq } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { PAGINATION_LIMIT } from '@/config';

import { MODALS, ENTITY_TYPES, DEFAULT_SERVICE_DEPENDENCIES_COLUMNS, COLOR_INDICATOR_TYPES } from '@/constants';

import { defaultColumnsToColumns } from '@/helpers/entities';
import { getEntityColor } from '@/helpers/color';
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
    openableRoot: {
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
      const headers = columns.map(({ label, value, colorIndicator }) => ({
        colorIndicator,

        sortable: false,
        text: label,
        value: value.match(/entity.|alarm./) ? value : `entity.${value}`,
      }));

      return [
        { sortable: false, text: '', value: 'no-events-icon' },

        ...headers,
      ];
    },

    items() {
      const items = dependenciesDenormalize({
        ids: this.rootIds,
        dependenciesByIds: this.dependenciesByIds,
        metaByIds: this.metaByIds,
      });

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

    getIconByEntity(entity) {
      return entity.type === ENTITY_TYPES.service ? '$vuetify.icons.engineering' : 'person';
    },

    getEntityColor(entity) {
      return getEntityColor(entity, COLOR_INDICATOR_TYPES.impactState);
    },

    isInRootIds(id) {
      return this.rootIds.includes(id);
    },

    showTreeOfDependenciesModal(dependency) {
      const { entity } = dependency;

      if (
        (!this.openableRoot && this.rootId === entity._id)
        || (!this.impact && entity.type !== ENTITY_TYPES.service)
      ) {
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

    getDependencyChildren(id) {
      return get(this.dependenciesByIds, [id, 'children']);
    },

    async fetchDependenciesById(id, params = { limit: PAGINATION_LIMIT }) {
      this.$set(this.pendingByIds, id, true);

      const { data, meta } = await this.fetchDependenciesList({ id, params });
      const { dependencies, result } = normalizeDependencies(data);

      Object.entries(dependencies).forEach(([dependencyId, dependency]) => {
        const children = this.getDependencyChildren(dependencyId) ?? dependency.children;

        this.$set(this.dependenciesByIds, dependencyId, { ...dependency, children });
      });

      if (this.dependenciesByIds[id] && result.length) {
        const oldChildren = this.getDependencyChildren(id) ?? [];
        const children = uniq([...oldChildren, ...result]);

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
