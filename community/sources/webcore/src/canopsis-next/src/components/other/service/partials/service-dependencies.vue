<template>
  <c-treeview-data-table
    class="service-dependencies"
    :items="items"
    :headers="headers"
    :loading="hasActivePending"
    :load-children="loadChildren"
    item-key="key"
  >
    <template #expand="{ item }">
      <v-tooltip
        v-if="item.loadMore"
        right
      >
        <template #activator="{ on }">
          <v-btn
            v-on="on"
            :loading="pendingByIds[item.parentId]"
            fab
            small
            depressed
            @click="loadMore(item.parentId)"
          >
            <v-icon>more_horiz</v-icon>
          </v-btn>
        </template>
        <span>{{ $t('common.loadMore') }}</span>
      </v-tooltip>
      <v-btn
        v-else
        :color="getEntityColor(item.entity)"
        fab
        small
        depressed
        dark
        @click="showTreeOfDependenciesModal(item)"
      >
        <v-icon>{{ getIconByEntity(item.entity) }}</v-icon>
      </v-btn>
      <v-tooltip
        v-if="item.cycle"
        top
      >
        <template #activator="{ on }">
          <v-icon
            class="mx-1"
            v-on="on"
            color="error"
            size="14"
          >
            autorenew
          </v-icon>
        </template>
        <span>{{ $t('common.cycleDependency') }}</span>
      </v-tooltip>
    </template>
    <template #expand-append="{ item }">
      <div
        class="expand-append"
        v-if="includeRoot && isInRootIds(item._id)"
      >
        <v-icon>arrow_right_alt</v-icon>
        <v-chip
          class="ma-0"
          :color="getEntityColor(item.entity)"
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
</template>

<script>
import { get, uniq } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, ENTITY_TYPES, ENTITY_FIELDS, COLOR_INDICATOR_TYPES } from '@/constants';

import { getIconByEntityType } from '@/helpers/entities/entity/icons';
import { getEntityColor } from '@/helpers/entities/entity/color';
import {
  dependencyToTreeviewDependency,
  normalizeDependencies,
  getLoadMoreDenormalizedChild,
  dependenciesDenormalize,
} from '@/helpers/entities/service-dependencies/list';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

import ServiceDependenciesEntityCell from './service-dependencies-entity-cell.vue';

export default {
  components: { ServiceDependenciesEntityCell },
  mixins: [entitiesEntityDependenciesMixin],
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
  },
  data() {
    const dependenciesByIds = {};
    const rootIds = [];

    if (this.includeRoot) {
      const treeviewRoot = dependencyToTreeviewDependency(this.root, this.impact);

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
          value: `entity.${column.value}`,
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
  async mounted() {
    const ids = await this.fetchDependenciesById(this.rootId);

    if (!this.includeRoot) {
      this.rootIds = ids;
    }
  },
  methods: {
    getIconByEntity(entity) {
      return getIconByEntityType(entity.type);
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

      const { data, meta } = await this.fetchDependenciesList({
        id,
        params: {
          ...params,

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
  &, .expand-append {
    display: inline-flex;
    align-items: center;

    .v-chip__content {
      height: 20px;
    }
  }
}

</style>
