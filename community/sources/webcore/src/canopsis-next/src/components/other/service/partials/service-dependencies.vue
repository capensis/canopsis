<template>
  <v-layout column>
    <service-dependencies-show-type-field
      v-if="isCustomType"
      v-model="showType"
      class="mb-3"
    />
    <state-settings-summary
      v-if="showStateSetting"
      :state-setting="stateSetting"
      :pending="stateSettingPending"
      :entity="root"
    />
    <c-treeview-data-table
      ref="treeviewDataTable"
      :items="items"
      :headers="headers"
      :loading="hasActivePending"
      :load-children="loadChildren"
      class="service-dependencies"
      item-key="key"
    >
      <template #expand="{ item }">
        <service-dependencies-expand
          :item="item"
          :pending="pendingByIds[item.parentId]"
          @load="loadMore"
          @show="showTreeOfDependenciesModal"
        />
      </template>
      <template #expand-append="{ item }">
        <div
          v-if="includeRoot && isInRootIds(item._id)"
          class="expand-append"
        >
          <v-icon>arrow_right_alt</v-icon>
          <v-chip
            :color="getEntityColor(item.entity)"
            class="ma-0"
            text-color="white"
          >
            <span class="px-2 text-body-2 font-weight-bold">{{ item.entity.impact_state }}</span>
          </v-chip>
        </div>
      </template>
      <template #items="{ item }">
        <tr>
          <td
            v-for="(header, index) in headers"
            :key="header.value"
          >
            <c-no-events-icon
              v-if="!index"
              :value="item.entity | get('idle_since')"
              top
            />
            <service-dependencies-entity-cell
              v-else-if="item.entity"
              :item="item"
              :column="header"
            />
          </td>
        </tr>
      </template>
    </c-treeview-data-table>
  </v-layout>
</template>

<script>
import { get, uniq } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import {
  MODALS,
  ENTITY_TYPES,
  ENTITY_FIELDS,
  COLOR_INDICATOR_TYPES,
  TREE_OF_DEPENDENCIES_SHOW_TYPES,
} from '@/constants';

import { getEntityColor } from '@/helpers/entities/entity/color';
import {
  dependencyToTreeviewDependency,
  normalizeDependencies,
  getLoadMoreDenormalizedChild,
  dependenciesDenormalize,
} from '@/helpers/entities/service-dependencies/list';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';
import { checkStateSettingMixin } from '@/mixins/entities/check-state-setting';

import StateSettingsSummary from '@/components/other/state-setting/state-settings-summary.vue';
import ServiceDependenciesShowTypeField
  from '@/components/other/service/form/fields/service-dependencies-show-type-field.vue';

import ServiceDependenciesExpand from './service-dependencies-expand.vue';
import ServiceDependenciesEntityCell from './service-dependencies-entity-cell.vue';

export default {
  components: {
    ServiceDependenciesShowTypeField,
    StateSettingsSummary,
    ServiceDependenciesExpand,
    ServiceDependenciesEntityCell,
  },
  mixins: [entitiesEntityDependenciesMixin, checkStateSettingMixin],
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
    impact: {
      type: Boolean,
      default: false,
    },
    openableRoot: {
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
      pendingByIds: {},

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

    treeviewRoot() {
      return dependencyToTreeviewDependency(this.root);
    },

    hasActivePending() {
      return Object.values(this.pendingByIds).some(pending => pending);
    },

    headers() {
      return [
        { sortable: false, text: '', value: 'no-events-icon' },

        ...this.columns.map(column => ({
          ...column,
          isState: column.value?.endsWith(ENTITY_FIELDS.state),
        })),
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
  watch: {
    type() {
      this.setRootDependencies();
      this.fetchRootDependencies();
    },
    showType() {
      this.setRootDependencies();
      this.fetchRootDependencies();
    },
    showStateSetting: {
      immediate: true,
      handler(value) {
        if (value) {
          this.checkStateSetting(this.root);
        }
      },
    },
  },
  mounted() {
    this.setRootDependencies();
    this.fetchRootDependencies();
  },
  methods: {
    async fetchRootDependencies() {
      const ids = await this.fetchDependenciesById(this.rootId);

      if (!this.includeRoot) {
        this.rootIds = ids;
      }
    },

    setRootDependencies() {
      this.$refs.treeviewDataTable.clearOpened();

      const dependenciesByIds = {};
      const rootIds = [];

      if (this.includeRoot) {
        const treeviewRoot = dependencyToTreeviewDependency(this.root, this.impact);

        dependenciesByIds[treeviewRoot._id] = treeviewRoot;
        rootIds.push(treeviewRoot._id);
      }

      this.rootIds = rootIds;
      this.dependenciesByIds = dependenciesByIds;
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

    loadChildren(dependency) {
      return this.fetchDependenciesById(dependency._id);
    },

    getDependencyChildren(id) {
      return get(this.dependenciesByIds, [id, 'children']);
    },

    async fetchDependenciesById(id, params = { limit: PAGINATION_LIMIT }) {
      this.$set(this.pendingByIds, id, true);

      const selectedType = this.isCustomType
        ? this.showType
        : this.type;

      const { data, meta } = await this.fetchDependenciesList({
        id,
        params: {
          ...params,

          define_state: selectedType === TREE_OF_DEPENDENCIES_SHOW_TYPES.dependenciesDefiningTheState,
          with_flags: true,
        },
      });

      const { dependencies, result } = normalizeDependencies(data, this.impact);

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
  },
};
</script>

<style lang="scss" scoped>
.service-dependencies ::v-deep .v-treeview-node__label {
  overflow: initial;

  &, .expand-append {
    display: inline-flex;
    align-items: center;

    .v-chip__content {
      height: 20px;
    }
  }
}
</style>
