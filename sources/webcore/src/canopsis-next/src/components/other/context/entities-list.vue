<template lang="pug">
  div
    v-layout.white(justify-space-between, align-center)
      v-flex
        context-search(:query.sync="query")
      v-flex
        pagination(
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
          :hasAccessToListFilter="hasAccessToListFilter",
          @input="updateSelectedFilter",
          @update:condition="updateSelectedCondition",
          @update:filters="updateFilters",
          :entitiesType="$constants.ENTITIES_TYPES.entity"
        )
      v-flex.ml-4
        mass-actions-panel(:itemsIds="selectedIds")
      v-flex
        context-fab(v-if="hasAccessToCreateEntity")
    no-columns-table(v-if="!hasColumns")
    div(v-else)
      v-data-table(
        v-model="selected",
        :items="contextEntities",
        :headers="headers",
        :loading="contextEntitiesPending",
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
            div(v-if="column.value === 'enabled'")
              enabled-column(:value="props.item.enabled")
            ellipsis(
              v-else,
              :text="props.item | get(column.value, null, '')",
              :maxLetters="column.maxLetters"
            )
          td
            actions-panel(:item="props.item", :isEditingMode="isEditingMode")
        template(slot="expand", slot-scope="props")
          more-infos(:item="props.item", :tabId="tabId")
      v-layout.white(align-center)
        v-flex(xs10)
          pagination(
            :page="query.page",
            :limit="query.limit",
            :total="contextEntitiesMeta.total",
            @input="updateQueryPage"
          )
        v-flex(xs2)
          records-per-page(:value="query.limit", @input="updateRecordsPerPage")
</template>

<script>
import { omit, isString } from 'lodash';

import { USERS_RIGHTS } from '@/constants';
import prepareMainFilterToQueryFilter from '@/helpers/filter';

import Ellipsis from '@/components/tables/ellipsis.vue';
import ContextSearch from '@/components/other/context/search/context-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';
import EnabledColumn from '@/components/tables/enabled-column.vue';

import authMixin from '@/mixins/auth';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import widgetRecordsPerPageMixin from '@/mixins/widget/records-per-page';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

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
    Ellipsis,
    ContextSearch,
    RecordsPerPage,
    NoColumnsTable,
    FilterSelector,
    MoreInfos,
    ContextFab,
    ActionsPanel,
    MassActionsPanel,
    EnabledColumn,
  },
  mixins: [
    authMixin,
    widgetQueryMixin,
    widgetColumnsMixin,
    widgetPaginationMixin,
    widgetFilterSelectMixin,
    widgetRecordsPerPageMixin,
    entitiesContextEntityMixin,
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

    hasAccessToListFilter() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.editFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_RIGHTS.business.context.actions.userFilter);
    },
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

      const filters = ['mainFilter', 'searchFilter', 'typesFilter'].reduce((acc, filterKey) => {
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
  },
};
</script>
