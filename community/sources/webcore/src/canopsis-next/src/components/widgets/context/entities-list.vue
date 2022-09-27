<template lang="pug">
  div
    c-empty-data-table-columns(v-if="!hasColumns")
    c-advanced-data-table(
      v-else,
      :items="contextEntities",
      :headers="headers",
      :loading="contextEntitiesPending || columnsFiltersPending",
      :total-items="contextEntitiesMeta.total_count",
      :pagination.sync="pagination",
      :toolbar-props="toolbarProps",
      select-all,
      expand,
      no-pagination
    )
      template(#toolbar="")
        v-flex
          c-advanced-search-field(
            :query.sync="query",
            :columns="columns",
            :tooltip="$t('search.contextAdvancedSearch')"
          )
        v-flex(v-if="hasAccessToCategory")
          c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
        v-flex
          v-layout(row, align-center)
            filter-selector(
              :label="$t('settings.selectAFilter')",
              :filters="userPreference.filters",
              :locked-filters="widget.filters",
              :value="mainFilter",
              :locked-value="lockedFilter",
              :disabled="!hasAccessToListFilters && !hasAccessToUserFilter",
              @input="updateSelectedFilter"
            )
            filters-list-btn(
              :widget-id="widget._id",
              :addable="hasAccessToAddFilter",
              :editable="hasAccessToEditFilter",
              private,
              with-alarm,
              with-entity,
              with-pbehavior
            )
        v-flex
          v-checkbox(
            :input-value="query.no_events",
            :label="$t('context.noEventsFilter')",
            color="primary",
            @change="updateNoEvents"
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
              :limit="query.limit",
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
      template(#actions="{ item }")
        actions-panel(:item="item", :editing="editing")
      template(#expand="{ item }")
        entities-list-expand-panel(
          :item="item",
          :widget="widget",
          :tab-id="tabId",
          :columns-filters="columnsFilters"
        )
      template(#mass-actions="{ selected, clearSelected }")
        mass-actions-panel.ml-3(:items="selected", @clear:items="clearSelected")

    c-table-pagination(
      :total-items="contextEntitiesMeta.total_count",
      :rows-per-page="query.limit",
      :page="query.page",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage"
    )
</template>

<script>
import { isObject } from 'lodash';

import { USERS_PERMISSIONS } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetColumnsContextMixin } from '@/mixins/widget/columns';
import { exportCsvMixinCreator } from '@/mixins/widget/export';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';
import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';
import { permissionsWidgetsContextFilters } from '@/mixins/permissions/widgets/context/filters';
import { permissionsWidgetsContextCategory } from '@/mixins/permissions/widgets/context/category';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';

import EntityColumnCell from './columns-formatting/entity-column-cell.vue';
import EntitiesListExpandPanel from './partials/entities-list-expand-panel.vue';
import ContextFab from './actions/context-fab.vue';
import ActionsPanel from './actions/actions-panel.vue';
import MassActionsPanel from './actions/mass-actions-panel.vue';

export default {
  components: {
    FilterSelector,
    FiltersListBtn,
    EntitiesListExpandPanel,
    ContextFab,
    EntityColumnCell,
    ActionsPanel,
    MassActionsPanel,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetColumnsContextMixin,
    widgetFilterSelectMixin,
    entitiesContextEntityMixin,
    entitiesAlarmColumnsFiltersMixin,
    permissionsWidgetsContextFilters,
    permissionsWidgetsContextCategory,
    exportCsvMixinCreator({
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
    editing: {
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
    updateNoEvents(noEvents) {
      this.updateContentInUserPreference({
        noEvents,
      });

      this.query = {
        ...this.query,

        no_events: noEvents,
      };
    },

    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,

        category: categoryId,
      };
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
      const {
        widgetExportColumns,
        widgetColumns,
        exportCsvSeparator,
        exportCsvDatetimeFormat,
      } = this.widget.parameters;
      const columns = widgetExportColumns && widgetExportColumns.length
        ? widgetExportColumns
        : widgetColumns;

      this.exportAsCsv({
        name: `${this.widget._id}-${new Date().toLocaleString()}`,
        widgetId: this.widget._id,
        data: {
          fields: columns.map(({ label, value }) => ({ label, name: value })),
          search: query.search,
          category: query.category,
          filters: query.filters,
          separator: exportCsvSeparator,
          /**
           * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
           */
          time_format: isObject(exportCsvDatetimeFormat) ? exportCsvDatetimeFormat.value : exportCsvDatetimeFormat,
        },
      });
    },
  },
};
</script>
