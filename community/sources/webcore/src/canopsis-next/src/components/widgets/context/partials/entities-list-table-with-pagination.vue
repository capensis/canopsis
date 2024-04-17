<template>
  <div>
    <c-empty-data-table-columns v-if="!columns.length" />
    <c-advanced-data-table
      v-else
      :items="entities"
      :headers="headers"
      :loading="pending || columnsFiltersPending"
      :total-items="meta.total_count"
      :options.sync="options"
      :toolbar-props="toolbarProps"
      :select-all="selectable"
      expand
      no-pagination
    >
      <template #toolbar="">
        <slot name="toolbar" />
        <v-flex
          v-if="columns.length"
          xs12
        >
          <v-layout
            wrap
            align-center
          >
            <c-pagination
              :page="options.page"
              :limit="options.itemsPerPage"
              :total="meta.total_count"
              type="top"
              @input="updatePage"
            />
          </v-layout>
        </v-flex>
      </template>
      <template
        v-for="column in columns"
        #[column.value]="{ item }"
      >
        <entity-column-cell
          :key="column.value"
          :entity="item"
          :column="column"
          :columns-filters="columnsFilters"
          :show-root-cause-by-state-click="showRootCauseByStateClick"
          @click:state="openRootCauseDiagram"
        />
      </template>
      <template #actions="{ item }">
        <actions-panel :item="item" />
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
          class="ml-3"
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
    />
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { getPageForNewItemsPerPage } from '@/helpers/pagination';

import { authMixin } from '@/mixins/auth';
import { widgetOptionsMixin } from '@/mixins/widget/options';
import { entitiesAlarmColumnsFiltersMixin } from '@/mixins/entities/associative-table/alarm-columns-filters';

import EntityColumnCell from '../columns-formatting/entity-column-cell.vue';
import ActionsPanel from '../actions/actions-panel.vue';
import MassActionsPanel from '../actions/mass-actions-panel.vue';

import EntitiesListExpandPanel from './entities-list-expand-panel.vue';

export default {
  components: {
    EntitiesListExpandPanel,
    EntityColumnCell,
    ActionsPanel,
    MassActionsPanel,
  },
  mixins: [
    authMixin,
    widgetOptionsMixin,
    entitiesAlarmColumnsFiltersMixin,
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
      return this.columns.length
        ? [
          ...this.columns,

          { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
        ]
        : [];
    },

    showRootCauseByStateClick() {
      return this.widget.parameters.showRootCauseByStateClick ?? true;
    },
  },
  async mounted() {
    this.columnsFiltersPending = true;
    this.columnsFilters = await this.fetchAlarmColumnsFiltersList();
    this.columnsFiltersPending = false;
  },
  methods: {
    updateItemsPerPage(itemsPerPage) {
      this.$emit('update:query', {
        ...this.query,

        itemsPerPage,
        page: getPageForNewItemsPerPage(itemsPerPage, this.query.itemsPerPage, this.query.page),
      });
    },

    updatePage(page) {
      this.$emit('update:query', {
        ...this.query,

        page,
      });
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
