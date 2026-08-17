<template>
  <div>
    <c-empty-data-table-columns v-if="!preparedColumns.length" />
    <c-advanced-data-table
      v-else
      :items="entities"
      :headers="headers"
      :loading="pending || columnsFiltersPending || entityInfoPropertyPending"
      :total-items="meta.total_count"
      :options="options"
      :toolbar-props="toolbarProps"
      :select-all="selectable"
      :mass-actions-grid-parameters="widgetCurrentGridParameters"
      :default-keep-selected-after-action="defaultKeepSelectedAfterAction"
      expand
      no-pagination
      @update:options="updateOptions"
    >
      <template #toolbar>
        <slot name="toolbar" />
        <div v-if="preparedColumns.length" class="w-100">
          <c-layout align="center" justify="center" wrap>
            <c-pagination
              :page="options.page"
              :limit="options.itemsPerPage"
              :total="meta.total_count"
              type="top"
              @update:page="updatePage"
            />
          </c-layout>
        </div>
      </template>
      <template
        v-for="column in preparedColumns"
        :key="column.value"
        #[column.value]="{ item }"
      >
        <entity-column-cell
          :entity="item"
          :column="column"
          :columns-filters="columnsFilters"
          :show-root-cause-by-state-click="showRootCauseByStateClick"
          @click:state="openRootCauseDiagram"
        />
      </template>
      <template #actions="{ item }">
        <actions-panel
          :item="item"
          :default-inherited-pbehavior="widget.parameters.defaultInheritedPbehavior"
        />
      </template>
      <template #expand="{ item }">
        <entities-list-expand-panel
          :item="item"
          :columns-filters="columnsFilters"
          :service-dependencies-columns="widget.parameters.serviceDependenciesColumns"
          :resolved-alarms-columns="widget.parameters.resolvedAlarmsColumns"
          :active-alarms-columns="widget.parameters.activeAlarmsColumns"
          :charts="widget.parameters.charts"
          :tree-of-dependencies-show-type="widget.parameters.treeOfDependenciesShowType"
          :expand-grid-range-size="widget.parameters.expandGridRangeSize"
          :availability="widget.parameters.availability"
        />
      </template>
      <template #mass-actions="{ selected, clearSelected }">
        <mass-actions-panel
          :items="selected"
          :default-inherited-pbehavior="widget.parameters.defaultInheritedPbehavior"
          small
          @clear:items="clearSelected"
        />
      </template>
    </c-advanced-data-table>
    <c-table-pagination
      :total-items="meta.total_count"
      :items-per-page="options.itemsPerPage"
      :page="options.page"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @update:pagination="updatePaginationOptions"
    />
  </div>
</template>

<script>
import { isEqual, pick } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP, WIDGET_GRID_SIZES_KEYS } from '@/constants';

import {
  convertDataTableSortToQuery,
  convertQuerySortToDataTableSort,
} from '@/helpers/entities/shared/query';

import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';
import { widgetColumnsEntityInfoPropertyMixin } from '@/mixins/widget/columns/entity-info-property';

import ActionsPanel from '../actions/actions-panel.vue';
import MassActionsPanel from '../actions/mass-actions-panel.vue';
import EntityColumnCell from '../columns-formatting/entity-column-cell.vue';

import EntitiesListExpandPanel from './entities-list-expand-panel.vue';

export default {
  components: {
    ActionsPanel,
    EntitiesListExpandPanel,
    EntityColumnCell,
    MassActionsPanel,
  },

  mixins: [
    entitiesAlarmColumnsFiltersMixin,
    widgetColumnsEntityInfoPropertyMixin,
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    entities: {
      type: Array,
      required: true,
    },
    meta: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:query'],

  data() {
    return {
      columnsFilters: [],
      columnsFiltersPending: false,
    };
  },

  computed: {
    displayGridSize() {
      return DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.$vuetify.display.name] ?? WIDGET_GRID_SIZES_KEYS.desktop;
    },

    widgetCurrentGridParameters() {
      const gridParameters = this.widget.grid_parameters ?? {};

      return pick(gridParameters[this.displayGridSize] ?? gridParameters, ['x', 'y', 'w', 'h']);
    },

    defaultKeepSelectedAfterAction() {
      return this.widget.parameters.keepSelectedAfterAction ?? false;
    },

    toolbarProps() {
      return {
        justify: 'space-between',
        align: 'center',
      };
    },

    preparedColumns() {
      return this.columns.map(column => ({
        ...column,

        text: column.label || this.findAliasByColumnValue(column.value) || column.text,
      }));
    },

    headers() {
      if (!this.preparedColumns.length) {
        return [];
      }

      return this.hideActions
        ? this.preparedColumns
        : [
            ...this.preparedColumns,

            { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
          ];
    },

    showRootCauseByStateClick() {
      return this.widget.parameters.showRootCauseByStateClick ?? true;
    },

    options: {
      get() {
        const {
          page = 1,
          itemsPerPage,
          limit,
          sortBy = [],
          sortDesc = [],
        } = this.query;

        return {
          page,
          itemsPerPage: itemsPerPage || limit || PAGINATION_LIMIT,
          sortBy: convertQuerySortToDataTableSort(sortBy, sortDesc),
        };
      },

      set(newOptions = {}) {
        const {
          page = 1,
          itemsPerPage = PAGINATION_LIMIT,
          sortBy = [],
          sortDesc = [],
        } = newOptions;
        const sort = sortBy[0]?.key
          ? convertDataTableSortToQuery(sortBy)
          : { sortBy, sortDesc };
        const convertedOptions = {
          page,
          itemsPerPage,
          sortBy: sort.sortBy,
          sortDesc: sort.sortDesc,
        };
        const oldOptions = {
          page: this.query.page ?? 1,
          itemsPerPage: this.query.itemsPerPage ?? PAGINATION_LIMIT,
          sortBy: this.query.sortBy ?? [],
          sortDesc: this.query.sortDesc ?? [],
        };

        if (isEqual(convertedOptions, oldOptions)) {
          return;
        }

        this.$emit('update:query', {
          ...this.query,
          ...convertedOptions,
        });
      },
    },
  },

  async mounted() {
    this.fetchEntityInfoPropertiesList({ params: { paginate: false } });

    this.columnsFiltersPending = true;
    this.columnsFilters = await this.fetchAlarmColumnsFiltersList();
    this.columnsFiltersPending = false;
  },

  methods: {
    updateOptions(options) {
      this.options = options;
    },

    updateItemsPerPage(itemsPerPage) {
      this.$emit('update:query', { ...this.query, itemsPerPage });
    },

    updatePage(page) {
      this.$emit('update:query', { ...this.query, page });
    },

    updatePaginationOptions({ page, itemsPerPage }) {
      this.$emit('update:query', { ...this.query, page, itemsPerPage });
    },

    openRootCauseDiagram(entity) {
      this.$modals.show({
        name: MODALS.entitiesRootCauseDiagram,
        config: {
          entity,
          colorIndicator: this.widget.parameters.rootCauseColorIndicator,
        },
      });
    },
  },
};
</script>
