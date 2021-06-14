<template lang="pug">
  div
    v-layout.white(justify-space-between, align-center)
      v-flex
        advanced-search(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.contextAdvancedSearch')"
        )
      v-flex
        c-pagination(
          v-if="hasColumns",
          :page="query.page",
          :limit="query.limit",
          :total="contextEntitiesMeta.total",
          type="top",
          @input="updateQueryPage"
        )
      v-flex
        filter-selector(
          :label="$t('settings.selectAFilter')",
          :filters="viewFilters",
          :lockedFilters="widgetViewFilters",
          :value="mainFilter",
          :condition="mainFilterCondition",
          :hasAccessToEditFilter="hasAccessToEditFilter",
          :hasAccessToUserFilter="hasAccessToUserFilter",
          :hasAccessToListFilters="hasAccessToListFilters",
          @input="updateSelectedFilter",
          @update:condition="updateSelectedCondition",
          @update:filters="updateFilters",
          :entitiesType="$constants.ENTITIES_TYPES.entity"
        )
      v-flex.ml-4
        mass-actions-panel(:itemsIds="selectedIds")
      v-flex
        context-fab(v-if="hasAccessToCreateEntity")
      v-flex(v-if="hasAccessToExportAsCsv")
        c-action-btn(
          :loading="!!contextExportPending",
          :tooltip="$t('settings.exportAsCsv')",
          icon="cloud_download",
          color="black",
          @click="exportContextList"
        )
    c-empty-data-table-columns(v-if="!hasColumns")
    div(v-else)
      v-data-table(
        v-model="selected",
        :items="contextEntities",
        :headers="headers",
        :loading="contextEntitiesPending || alarmColumnFiltersPending",
        :total-items="contextEntitiesMeta.total",
        :pagination.sync="vDataTablePagination",
        item-key="_id",
        select-all,
        hide-actions
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="headerCell", slot-scope="props")
          span {{ props.header.text }}
        template(slot="items", slot-scope="props")
          td
            v-checkbox(primary, hide-details, v-model="props.selected")
          td(
            v-for="column in columns",
            @click="props.expanded = !props.expanded"
          )
            entity-column-cell(:entity="props.item", :column="column")
          td
            actions-panel(:item="props.item", :isEditingMode="isEditingMode")
        template(slot="expand", slot-scope="props")
          more-infos(:item="props.item", :tabId="tabId")
      c-table-pagination(
        :total-items="contextEntitiesMeta.total",
        :rows-per-page="query.limit",
        :page="query.page",
        @update:page="updateQueryPage",
        @update:rows-per-page="updateRecordsPerPage"
      )
</template>

<script>
import { omit, isString } from 'lodash';

import { USERS_RIGHTS } from '@/constants';
import { prepareMainFilterToQueryFilter } from '@/helpers/filter';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import AdvancedSearch from '@/components/common/search/advanced-search.vue';

import authMixin from '@/mixins/auth';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import alarmColumnFiltersMixin from '@/mixins/entities/alarm-column-filters';
import widgetExportMixinCreator from '@/mixins/widget/export';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

import EntityColumnCell from './columns-formatting/entity-column-cell.vue';
import MoreInfos from './more-infos/more-infos.vue';
import ContextFab from './actions/context-fab.vue';
import ActionsPanel from './actions/actions-panel.vue';
import MassActionsPanel from './actions/mass-actions-panel.vue';

/**
 * Entities list
 *
 * @module context
 *
 * @prop {Object} widget - Object representing the widget
 * @prop {Array} columns - List of entities columns
 *
 * @event openSettings#click
 */
export default {
  components: {
    FilterSelector,
    MoreInfos,
    ContextFab,
    EntityColumnCell,
    ActionsPanel,
    MassActionsPanel,
    AdvancedSearch,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetColumnsMixin,
    alarmColumnFiltersMixin,
    widgetFilterSelectMixin,
    entitiesContextEntityMixin,
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
      selected: [],
    };
  },
  computed: {
    selectedIds() {
      return this.selected.map(item => item._id);
    },

    headers() {
      if (this.hasColumns) {
        return [...this.columns, { text: this.$t('common.actionsLabel'), sortable: false }];
      }

      return [];
    },

    hasAccessToCreateEntity() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.createEntity);
    },

    hasAccessToListFilters() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.editFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.userFilter);
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.exportAsCsv);
    },
  },
  mounted() {
    this.fetchAlarmColumnFilters();
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
        'mainFilter',
        'searchFilter',
        'typesFilter',
      ]);

      query.start = ((this.query.page - 1) * this.query.limit) || 0;

      if (this.query.sortKey) {
        query.sort = [{
          property: this.query.sortKey,
          direction: this.query.sortDir,
        }];
      }

      const filters = ['mainFilter', 'typesFilter'].reduce((acc, filterKey) => {
        const queryFilter = isString(this.query[filterKey]) ? JSON.parse(this.query[filterKey]) : this.query[filterKey];

        if (queryFilter) {
          acc.push(queryFilter);
        }

        return acc;
      }, []);

      if (filters.length) {
        query._filter = {
          $and: filters,
        };
      }

      return query;
    },
    fetchList() {
      if (this.hasColumns) {
        this.fetchContextEntitiesList({
          widgetId: this.widget._id,
          params: this.getQuery(),
        });
      }
    },

    updateQueryBySelectedFilterAndCondition(filter, condition) {
      this.query = {
        ...this.query,

        mainFilter: prepareMainFilterToQueryFilter(filter, condition),
      };
    },

    exportContextList() {
      const query = this.getQuery();

      this.exportWidgetAsCsv({
        params: {
          filter: query._filter,
          search: query.search,
          active_columns: query.active_columns,
          separator: this.widget.parameters.exportCsvSeparator,
        },
        name: `${this.widget._id}-${new Date().toLocaleString()}`,
      });
    },
  },
};
</script>
