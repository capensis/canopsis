<template lang="pug">
  div
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        c-advanced-search-field(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.alarmAdvancedSearch')"
        )
      v-flex(v-if="hasAccessToCategory")
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
      v-flex(v-if="hasAccessToCorrelation")
        v-switch(
          :value="query.correlation",
          :label="$t('common.correlation')",
          color="primary",
          @change="updateCorrelation"
        )
      v-flex
        v-layout(row, wrap, align-center)
          filter-selector(
            :label="$t('settings.selectAFilter')",
            :filters="userPreference.filters",
            :locked-filters="widget.filters",
            :value="mainFilter",
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
        alarms-list-remediation-instructions-filters(
          :filters.sync="remediationInstructionsFilters",
          :locked-filters.sync="widgetRemediationInstructionsFilters",
          :editable="hasAccessToEditRemediationInstructionsFilter",
          :addable="hasAccessToUserRemediationInstructionsFilter",
          :has-access-to-list-filters="hasAccessToListRemediationInstructionsFilters"
        )
      v-flex
        v-chip.primary.white--text(
          v-if="activeRange",
          close,
          label,
          @input="removeHistoryFilter"
        ) {{ $t(`quickRanges.types.${activeRange.value}`) }}
        c-action-btn(
          :tooltip="$t('liveReporting.button')",
          :color="activeRange ? 'primary' : 'black'",
          icon="schedule",
          @click="showEditLiveReportModal"
        )
      v-flex(v-if="hasAccessToExportAsCsv")
        c-action-btn(
          :loading="!!alarmsExportPending",
          :tooltip="$t('settings.exportAsCsv')",
          icon="cloud_download",
          color="black",
          @click="exportAlarmsList"
        )
    v-layout.alarms-list__top-pagination.white.px-4(row, wrap, align-center)
      c-density-btn-toggle(:value="userPreference.content.dense", @change="updateDense")
      c-pagination(
        v-if="hasColumns",
        :page="query.page",
        :limit="query.limit",
        :total="alarmsMeta.total_count",
        type="top",
        @input="updateQueryPage"
      )
    alarms-list-table(
      ref="alarmsTable",
      :widget="widget",
      :alarms="alarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="pagination",
      :loading="alarmsPending",
      :is-tour-enabled="isTourEnabled",
      :hide-children="!query.correlation",
      :columns="columns",
      :sticky-header="widget.parameters.sticky_header",
      :dense="dense",
      :refresh-alarms-list="fetchList",
      :selected-tag="query.tag",
      selectable,
      expandable,
      @select:tag="selectTag",
      @clear:tag="clearTag"
    )
      c-table-pagination(
        :total-items="alarmsMeta.total_count",
        :rows-per-page="query.limit",
        :page="query.page",
        @update:page="updateQueryPage",
        @update:rows-per-page="updateRecordsPerPage"
      )
    alarms-expand-panel-tour(v-if="isTourEnabled", :callbacks="tourCallbacks")
</template>

<script>
import { omit, pick, isObject } from 'lodash';

import { MODALS, TOURS, USERS_PERMISSIONS } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetColumnsAlarmMixin } from '@/mixins/widget/columns';
import { exportCsvMixinCreator } from '@/mixins/widget/export';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetRemediationInstructionsFilterMixin } from '@/mixins/widget/remediation-instructions-filter-select';
import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';
import { permissionsWidgetsAlarmsListCorrelation } from '@/mixins/permissions/widgets/alarms-list/correlation';
import { permissionsWidgetsAlarmsListCategory } from '@/mixins/permissions/widgets/alarms-list/category';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import { permissionsWidgetsAlarmsListRemediationInstructionsFilters }
  from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';

import AlarmsListTable from './partials/alarms-list-table.vue';
import AlarmsExpandPanelTour from './expand-panel/alarms-expand-panel-tour.vue';
import AlarmsListRemediationInstructionsFilters from './partials/alarms-list-remediation-instructions-filters.vue';

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {Object} widget - Object representing the widget
 *
 * @event openSettings#click
 */
export default {
  components: {
    FilterSelector,
    FiltersListBtn,
    AlarmsListTable,
    AlarmsExpandPanelTour,
    AlarmsListRemediationInstructionsFilters,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetColumnsAlarmMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    widgetRemediationInstructionsFilterMixin,
    entitiesAlarmMixin,
    entitiesAlarmTagMixin,
    entitiesAlarmDetailsMixin,
    permissionsWidgetsAlarmsListCategory,
    permissionsWidgetsAlarmsListCorrelation,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
    exportCsvMixinCreator({
      createExport: 'createAlarmsListExport',
      fetchExport: 'fetchAlarmsListExport',
      fetchExportFile: 'fetchAlarmsListCsvFile',
    }),
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tabId: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    tourCallbacks() {
      return {
        onPreviousStep: this.onTourPreviousStep,
        onNextStep: this.onTourNextStep,
      };
    },

    isTourEnabled() {
      return this.checkIsTourEnabled(TOURS.alarmsExpandPanel)
        && !!this.alarms.length;
    },

    activeRange() {
      const { tstart, tstop } = this.query;

      if (tstart || tstop) {
        return findQuickRangeValue(tstart, tstop);
      }

      return null;
    },

    firstAlarmExpanded() {
      const [alarm] = this.alarms;

      return alarm && this.$refs.alarmsTable.expanded[alarm._id];
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv);
    },

    dense() {
      return this.userPreference.content.dense ?? this.widget.parameters.dense;
    },
  },
  methods: {
    refreshExpanded() {
      if (this.$refs.alarmsTable?.expanded) {
        Object.entries(this.$refs.alarmsTable.expanded).forEach(([id, expanded]) => {
          if (expanded && !this.alarms.some(alarm => alarm._id === id)) {
            this.$set(this.$refs.alarmsTable.expanded, id, false);
          }
        });
      }
    },

    selectTag(tag) {
      this.query = {
        ...this.query,

        tag,
      };
    },

    clearTag() {
      this.query = omit(this.query, ['tag']);
    },

    updateCorrelation(correlation) {
      this.updateContentInUserPreference({
        isCorrelationEnabled: correlation,
      });

      this.query = {
        ...this.query,

        correlation,
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

    updateRecordsPerPage(limit) {
      this.updateContentInUserPreference({
        itemsPerPage: limit,
      });

      this.query = {
        ...this.query,

        limit,
      };
    },

    updateDense(dense) {
      this.updateContentInUserPreference({
        dense,
      });
    },

    expandFirstAlarm() {
      if (!this.firstAlarmExpanded) {
        this.$set(this.$refs.alarmsTable.expanded, this.alarms[0]._id, true);
      }
    },

    onTourPreviousStep(currentStep) {
      if (currentStep !== 1) {
        this.expandFirstAlarm();
      }

      return this.$nextTick();
    },

    onTourNextStep() {
      this.expandFirstAlarm();

      return this.$nextTick();
    },

    removeHistoryFilter() {
      this.query = omit(this.query, ['tstart', 'tstop']);
    },

    showEditLiveReportModal() {
      this.$modals.show({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['tstart', 'tstop', 'time_field']),
          action: params => this.query = { ...this.query, ...params },
        },
      });
    },

    async fetchList() {
      if (this.hasColumns) {
        const params = this.getQuery();

        this.fetchAlarmsDetailsList({ widgetId: this.widget._id });

        if (!this.alarmTagsPending) {
          this.fetchAlarmTagsList({ params: { paginate: false } });
        }

        if (!this.alarmsPending) {
          await this.fetchAlarmsList({
            widgetId: this.widget._id,
            params,
          });

          this.refreshExpanded();
        }
      }
    },

    exportAlarmsList() {
      const query = this.getQuery();
      const {
        widgetExportColumns,
        widgetColumns,
        exportCsvSeparator,
        exportCsvDatetimeFormat,
      } = this.widget.parameters;
      const columns = widgetExportColumns?.length
        ? widgetExportColumns
        : widgetColumns;

      this.exportAsCsv({
        name: `${this.widget._id}-${new Date().toLocaleString()}`,
        widgetId: this.widget._id,
        data: {
          ...pick(query, ['search', 'category', 'correlation', 'opened', 'tstart', 'tstop']),

          fields: columns.map(({ label, value }) => ({ label, name: value })),
          filter: JSON.stringify(query.filter),
          separator: exportCsvSeparator,
          /**
           * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
           */
          time_format: isObject(exportCsvDatetimeFormat)
            ? exportCsvDatetimeFormat.value
            : exportCsvDatetimeFormat,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.alarms-list__top-pagination {
  min-height: 46px;
}
</style>
