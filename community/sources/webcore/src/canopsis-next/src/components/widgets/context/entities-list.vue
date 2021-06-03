<template lang="pug">
  div
    c-empty-data-table-columns(v-if="!hasColumns")
    c-advanced-data-table(
      v-else,
      :items="contextEntities",
      :headers="headers",
      :loading="contextEntitiesPending || columnsFiltersPending",
      :total-items="contextEntitiesMeta.total_count",
      :pagination.sync="vDataTablePagination",
      :toolbar-props="toolbarProps",
      select-all,
      expand,
      no-pagination
    )
      template(slot="toolbar", slot-scope="props")
        v-flex
          c-advanced-search(
            :query.sync="query",
            :columns="columns",
            :tooltip="$t('search.contextAdvancedSearch')"
          )
        v-flex(v-if="hasAccessToCategory")
          c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
        v-flex
          filter-selector(
            :label="$t('settings.selectAFilter')",
            :filters="viewFilters",
            :locked-filters="widgetViewFilters",
            :value="mainFilter",
            :condition="mainFilterCondition",
            :has-access-to-edit-filter="hasAccessToEditFilter",
            :has-access-to-user-filter="hasAccessToUserFilter",
            :has-access-to-list-filters="hasAccessToListFilters",
            :entitiesType="$constants.ENTITIES_TYPES.entity",
            @input="updateSelectedFilter",
            @update:condition="updateSelectedCondition",
            @update:filters="updateFilters"
          )
        v-flex(v-if="hasAccessToCreateEntity")
          context-fab
        v-flex(v-if="hasAccessToExportAsCsv")
          c-action-btn(
            :loading="!!contextExportPending",
            :tooltip="$t('settings.exportAsCsv')",
            icon="cloud_download",
            color="black",
            @click="exportContextList"
          )
        v-flex(v-if="hasColumns", xs12)
          v-layout(row, wrap, align-center)
            c-pagination(
              :page="query.page",
              :limit="query.rowsPerPage",
              :total="contextEntitiesMeta.total_count",
              type="top",
              @input="updateQueryPage"
            )

      template(v-for="column in columns", :slot="column.value", slot-scope="props")
        entity-column-cell(
          :entity="props.item",
          :column="column",
          :columns-filters="columnsFilters"
        )
      template(slot="actions", slot-scope="props")
        actions-panel(:item="props.item", :is-editing-mode="isEditingMode")

      template(slot="expand", slot-scope="props")
        entities-list-expand-panel(
          :item="props.item",
          :widget="widget",
          :tab-id="tabId",
          :columns-filters="columnsFilters"
        )

      template(slot="mass-actions", slot-scope="props")
        mass-actions-panel.ml-3(:items-ids="props.selected")
    c-table-pagination(
      :total-items="contextEntitiesMeta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { omit, isString } from 'lodash';

import { USERS_PERMISSIONS } from '@/constants';

import { prepareMainFilterToQueryFilter } from '@/helpers/filter';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetExportMixinCreator from '@/mixins/widget/export';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';
import { permissionsWidgetsContextEntityFilters } from '@/mixins/permissions/widgets/context-entity/filters';
import { permissionsWidgetsContextEntityCategory } from '@/mixins/permissions/widgets/context-entity/category';

import EntityColumnCell from './columns-formatting/entity-column-cell.vue';
import EntitiesListExpandPanel from './partials/entities-list-expand-panel.vue';
import ContextFab from './actions/context-fab.vue';
import ActionsPanel from './actions/actions-panel.vue';
import MassActionsPanel from './actions/mass-actions-panel.vue';

export default {
  components: {
    FilterSelector,
    EntitiesListExpandPanel,
    ContextFab,
    EntityColumnCell,
    ActionsPanel,
    MassActionsPanel,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetColumnsMixin,
    widgetFilterSelectMixin,
    entitiesContextEntityMixin,
    entitiesAlarmColumnsFiltersMixin,
    permissionsWidgetsContextEntityFilters,
    permissionsWidgetsContextEntityCategory,
    widgetExportMixinCreator({
      createExport: 'createContextExport',
      fetchExport: 'fetchContextExport',
      fetchExportFile: 'fetchContextCsvFile',
    }),
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      required: true,
    },
  },
  data() {
    return {
      columnsFilters: [],
      columnsFiltersPending: false,
    };
  },
  computed: {
    toolbarProps() {
      return {
        'justify-space-between': true,
        'align-center': true,
      };
    },

    headers() {
      if (this.hasColumns) {
        return [
          ...this.columns,
          { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
        ];
      }

      return [];
    },

    hasAccessToCreateEntity() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.createEntity);
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.context.actions.exportAsCsv);
    },
  },
  async mounted() {
    this.columnsFiltersPending = true;
    this.columnsFilters = await this.fetchAlarmColumnsFiltersList();
    this.columnsFiltersPending = false;
  },
  methods: {
    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        category: categoryId,
      });

      this.query = {
        ...this.query,

        category: categoryId,
      };
    },

    updateQueryBySelectedFilterAndCondition(filter, condition) {
      this.query = {
        ...this.query,

        page: 1,
        mainFilter: prepareMainFilterToQueryFilter(filter, condition),
      };
    },

    getQuery() {
      const query = omit(this.query, [
        'sortKey',
        'sortDir',
        'mainFilter',
        'searchFilter',
        'typesFilter',
      ]);

      if (this.query.sortKey) {
        query.sort = this.query.sortDir.toLowerCase();
        query.sort_by = this.query.sortKey;
      }

      const filters = ['mainFilter', 'typesFilter'].reduce((acc, filterKey) => {
        const queryFilter = isString(this.query[filterKey]) ? JSON.parse(this.query[filterKey]) : this.query[filterKey];

        if (queryFilter) {
          acc.push(queryFilter);
        }

        return acc;
      }, []);

      if (filters.length) {
        query.filter = {
          $and: filters,
        };
      }

      return query;
    },


    fetchList() {
      if (this.hasColumns) {
        const params = this.getQuery();
        params.with_flags = true;

        this.fetchContextEntitiesList({
          widgetId: this.widget._id,
          params,
        });
      }
    },

    exportContextList() {
      const query = this.getQuery();

      this.exportWidgetAsCsv({
        name: `${this.widget._id}-${new Date().toLocaleString()}`,
        params: {
          filter: query._filter,
          search: query.search,
          active_columns: query.active_columns,
          separator: this.widget.parameters.exportCsvSeparator,
          time_format: this.widget.parameters.exportCsvDatetimeFormat,
        },
      });
    },
  },
};
</script>