<template>
  <c-layout column>
    <service-dependencies-show-type-field v-if="isCustomType" v-model="showType" class="mb-3" />
    <state-settings-summary
      v-if="showStateSetting"
      :pending="pendingStateSetting"
      :state-setting="stateSetting"
      :entity="root"
      class="mb-3"
    />
    <v-table class="service-dependencies" density="compact">
      <thead>
        <tr>
          <th class="service-dependencies__icons" />
          <th v-for="column in headers" :key="column.value">
            {{ column.text }}
          </th>
        </tr>
        <tr>
          <th :colspan="headers.length + 1" class="service-dependencies__loader pa-0">
            <v-progress-linear
              :active="hasActivePending"
              color="primary"
              height="2"
              indeterminate
            />
          </th>
        </tr>
      </thead>
      <tbody>
        <service-dependencies-row
          v-for="item in items"
          :key="item.key"
          :item="item"
          :columns="headers"
          :level="0"
          :opened="opened"
          :pending-by-ids="pendingByIds"
          @toggle="toggle"
          @load="loadMore"
          @show="showTreeOfDependenciesModal"
        />
      </tbody>
    </v-table>
  </c-layout>
</template>

<script>
import get from 'lodash/get';
import uniq from 'lodash/uniq';
import { mapActions } from 'pinia';

import { PAGINATION_LIMIT } from '@/config';
import {
  ENTITY_FIELDS,
  MODALS,
  TREE_OF_DEPENDENCIES_SHOW_TYPES,
} from '@/constants';

import Observer from '@/services/observer';

import { useServiceStore } from '@/stores/modules/service';
import { useStateSettingStore } from '@/stores/modules/state-setting';

import {
  dependenciesDenormalize,
  dependencyToTreeviewDependency,
  getLoadMoreDenormalizedChild,
  normalizeDependencies,
} from '@/helpers/entities/service-dependencies/list';

import ServiceDependenciesShowTypeField from '@/components/other/service/form/fields/service-dependencies-show-type-field.vue';
import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';

import ServiceDependenciesRow from './service-dependencies-row.vue';

export default {
  inject: {
    $periodicRefresh: {
      default: () => new Observer(),
    },
  },

  components: {
    ServiceDependenciesRow,
    ServiceDependenciesShowTypeField,
    StateSettingsSummary,
  },

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
      default: () => [],
    },
    impact: {
      type: Boolean,
      default: false,
    },
    showStateSetting: {
      type: Boolean,
      default: false,
    },
    type: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
    },
  },

  data() {
    return {
      rootIds: [],
      dependenciesByIds: {},
      metaByIds: {},
      opened: [],
      pendingByIds: {},
      pendingStateSetting: false,
      stateSetting: undefined,
      showType: TREE_OF_DEPENDENCIES_SHOW_TYPES.allDependencies,
    };
  },

  computed: {
    isCustomType() {
      return this.type === TREE_OF_DEPENDENCIES_SHOW_TYPES.custom;
    },

    rootId() {
      return this.root._id;
    },

    rootHasNextPage() {
      const meta = this.metaByIds[this.rootId] || {};

      return meta.page < meta.page_count;
    },

    hasActivePending() {
      return Object.values(this.pendingByIds).some(Boolean);
    },

    treeviewRoot() {
      return dependencyToTreeviewDependency(this.root);
    },

    headers() {
      return this.columns.map(column => ({
        ...column,
        isState: column.value?.endsWith(ENTITY_FIELDS.state),
      }));
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

  watch: {
    type() {
      this.refresh();
    },

    showType() {
      this.refresh();
    },

    showStateSetting: {
      immediate: true,
      handler(value) {
        if (value) {
          this.fetchEntityStateSetting();
        }
      },
    },
  },

  mounted() {
    if (this.includeRoot) {
      this.setRootDependencies();
    }

    this.refresh();
    this.$periodicRefresh?.register?.(this.refresh);
  },

  beforeUnmount() {
    this.$periodicRefresh?.unregister?.(this.refresh);
  },

  methods: {
    ...mapActions(useServiceStore, [
      'fetchServiceDependenciesWithoutStore',
      'fetchServiceImpactsWithoutStore',
    ]),
    ...mapActions(useStateSettingStore, {
      fetchStateSettingsList: 'fetchList',
    }),

    refresh() {
      this.fetchRootDependencies({ replace: true });
    },

    async fetchEntityStateSetting() {
      this.pendingStateSetting = true;

      try {
        const { data } = await this.fetchStateSettingsList({
          params: {
            _id: this.root._id,
            limit: 1,
          },
        });

        [this.stateSetting] = data ?? [];
      } catch (err) {
        console.error(err);
      } finally {
        this.pendingStateSetting = false;
      }
    },

    async fetchRootDependencies({ replace = false } = {}) {
      if (replace) {
        const { result, dependencies, meta } = await this.requestDependencies(this.rootId);

        this.replaceRootDependencies(result, dependencies, meta);

        return;
      }

      const ids = await this.fetchDependenciesById(this.rootId);

      if (!this.includeRoot) {
        this.rootIds = ids;
      }
    },

    replaceRootDependencies(result, dependencies, meta) {
      const dependenciesByIds = {};
      const rootIds = [];
      const opened = [];

      if (this.includeRoot) {
        const treeviewRoot = dependencyToTreeviewDependency(this.root, this.impact);

        dependenciesByIds[treeviewRoot._id] = {
          ...treeviewRoot,
          children: result,
        };
        rootIds.push(treeviewRoot._id);
        opened.push(treeviewRoot.key);
      }

      Object.entries(dependencies).forEach(([dependencyId, dependency]) => {
        dependenciesByIds[dependencyId] = {
          ...dependency,
          children: dependenciesByIds[dependencyId]?.children ?? dependency.children,
        };
      });

      this.dependenciesByIds = dependenciesByIds;
      this.rootIds = this.includeRoot ? rootIds : result;
      this.metaByIds = { [this.rootId]: meta };
      this.opened = opened;
    },

    setRootDependencies() {
      const dependenciesByIds = {};
      const rootIds = [];

      this.opened = [];

      if (this.includeRoot) {
        const treeviewRoot = dependencyToTreeviewDependency(this.root, this.impact);

        dependenciesByIds[treeviewRoot._id] = treeviewRoot;
        rootIds.push(treeviewRoot._id);
        this.opened.push(treeviewRoot.key);
      }

      this.rootIds = rootIds;
      this.dependenciesByIds = dependenciesByIds;
    },

    showTreeOfDependenciesModal(dependency) {
      const { entity } = dependency;
      const hasChildren = this.impact ? !!entity.impacts_count : !!entity.depends_count;

      if (this.rootId === entity._id || !hasChildren) {
        return;
      }

      this.$modals.show({
        name: MODALS.serviceDependencies,
        config: {
          root: dependency.entity,
          columns: this.columns,
          impact: this.impact,
        },
      });
    },

    async loadMore({ parentId } = {}) {
      const isRoot = this.rootId === parentId;
      const meta = this.metaByIds[parentId] || {};
      const params = {
        page: meta.page + 1,
        limit: PAGINATION_LIMIT,
      };
      const ids = await this.fetchDependenciesById(parentId, params);

      if (!this.includeRoot && isRoot) {
        this.rootIds.push(...ids);
      }
    },

    async toggle(item) {
      if (item.loadMore || item.cycle || !item.children) {
        return;
      }

      if (this.opened.includes(item.key)) {
        this.opened = this.opened.filter(key => key !== item.key);

        return;
      }

      if (!get(this.dependenciesByIds, [item._id, 'children', 'length'])) {
        await this.fetchDependenciesById(item._id);
      }

      this.opened = [...this.opened, item.key];
    },

    getDependencyChildren(id) {
      return get(this.dependenciesByIds, [id, 'children']);
    },

    async requestDependencies(id, params = { limit: PAGINATION_LIMIT }) {
      this.pendingByIds = {
        ...this.pendingByIds,
        [id]: true,
      };

      try {
        const selectedType = this.isCustomType ? this.showType : this.type;
        const fetchList = this.impact
          ? this.fetchServiceImpactsWithoutStore
          : this.fetchServiceDependenciesWithoutStore;
        const { data, meta } = await fetchList({
          id,
          params: {
            ...params,
            define_state:
              selectedType === TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState,
            with_flags: true,
          },
        });
        const { dependencies, result } = normalizeDependencies(data, this.impact);

        return { meta, dependencies, result };
      } finally {
        this.pendingByIds = {
          ...this.pendingByIds,
          [id]: false,
        };
      }
    },

    async fetchDependenciesById(id, params = { limit: PAGINATION_LIMIT }) {
      const { meta, dependencies, result } = await this.requestDependencies(id, params);

      Object.entries(dependencies).forEach(([dependencyId, dependency]) => {
        const children = this.getDependencyChildren(dependencyId) ?? dependency.children;

        this.dependenciesByIds = {
          ...this.dependenciesByIds,
          [dependencyId]: { ...dependency, children },
        };
      });

      if (this.dependenciesByIds[id] && result.length) {
        const oldChildren = this.getDependencyChildren(id) ?? [];
        const children = uniq([...oldChildren, ...result]);

        this.dependenciesByIds = {
          ...this.dependenciesByIds,
          [id]: {
            ...this.dependenciesByIds[id],
            children,
          },
        };
      }

      this.metaByIds = {
        ...this.metaByIds,
        [id]: meta,
      };

      return result;
    },
  },
};
</script>

<style lang="scss" scoped>
.service-dependencies {
  width: 100%;

  &__icons {
    width: 1%;
  }

  &__loader {
    border-bottom: none !important;
    height: 0;
    line-height: 0;
  }
}
</style>
