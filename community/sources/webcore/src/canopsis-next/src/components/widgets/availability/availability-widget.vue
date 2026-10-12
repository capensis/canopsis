<template>
  <c-layout class="gap-3" column>
    <availability-widget-filters
      :query="query"
      :searches="searches"
      :widget-id="widget._id"
      :user-filters="userPreferenceFilters"
      :widget-filters="widget.filters"
      :show-interval="hasAccessToInterval"
      :show-filter="hasAccessToUserFilter"
      :show-export="!hideActions && hasAccessToExportAsCsv"
      :exporting="exporting"
      :max-value-filter-seconds="maxValueFilterSeconds"
      class="px-3 pt-3"
      @export="exportAvailabilityList"
      @update:filters="updateSelectedFilter"
      @update:interval="updateInterval"
      @update:trend="updateTrend"
      @update:type="updateShowType"
      @update:display-parameter="updateDisplayParameter"
      @update:value-filter="updateValueFilter"
      @submit:search="updateSearch"
      @reset:search="resetSearch"
      @toggle-pin:search="togglePinSearch"
      @remove:search="removeSearch"
    />

    <availability-list
      :availabilities="availabilities"
      :pending="availabilitiesPending"
      :total-items="availabilitiesMeta.total_count"
      :columns="widget.parameters.widget_columns"
      :display-parameter="query.displayParameter"
      :active-alarms-columns="widget.parameters.active_alarms_columns"
      :resolved-alarms-columns="widget.parameters.resolved_alarms_columns"
      :show-type="query.showType"
      :options="options"
      :interval="interval"
      :show-trend="query.showTrend"
      @update:options="updateOptions"
    />
  </c-layout>
</template>

<script>
import { isEmpty, pick } from 'lodash';

import {
  ADVANCED_SEARCH_FIELDS,
  AVAILABILITY_DISPLAY_PARAMETERS,
  AVAILABILITY_VALUE_FILTER_METHODS,
  TIME_UNITS,
} from '@/constants';

import { isOmitEqual } from '@/helpers/collection';
import { convertQueryIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { toSeconds } from '@/helpers/date/duration';
import { getAvailabilityFieldByDisplayParameterAndShowType } from '@/helpers/entities/availability/entity';
import {
  getAvailabilitiesTrendByInterval,
  prepareAvailabilityWidgetColumnValue,
} from '@/helpers/entities/availability/query';
import { getExportMetricDownloadFileUrl } from '@/helpers/entities/metric/url';
import { convertFiltersToQuery, convertSortToRequest } from '@/helpers/entities/shared/query';

import { entitiesAvailabilityMixin } from '@/mixins/entities/availability';
import { permissionsWidgetsAvailabilityExport } from '@/mixins/permissions/widgets/availability/export';
import { permissionsWidgetsAvailabilityFilters } from '@/mixins/permissions/widgets/availability/filters';
import { permissionsWidgetsAvailabilityInterval } from '@/mixins/permissions/widgets/availability/interval';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetAdvancedSearchSavedItemsMixin } from '@/mixins/widget/advanced-search/saved-items';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import AvailabilityList from '@/components/other/availability/availability-list.vue';
import AvailabilityWidgetFilters from '@/components/widgets/availability/partials/availability-widget-filters.vue';

export default {
  inject: ['$system'],

  components: { AvailabilityList, AvailabilityWidgetFilters },

  mixins: [
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
    widgetAdvancedSearchSavedItemsMixin,
    permissionsWidgetsAvailabilityFilters,
    permissionsWidgetsAvailabilityInterval,
    permissionsWidgetsAvailabilityExport,
    queryIntervalFilterMixin,
    entitiesAvailabilityMixin,
    exportMixinCreator({
      createExport: 'createAvailabilityExport',
      fetchExport: 'fetchAvailabilityExport',
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
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      exporting: false,
    };
  },

  computed: {
    interval() {
      return this.getIntervalQuery();
    },

    maxValueFilterSeconds() {
      if (this.interval.to === this.interval.from) {
        return toSeconds(1, TIME_UNITS.day);
      }

      return this.interval.to - this.interval.from;
    },

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },
  },

  methods: {
    customQueryCondition(query, oldQuery) {
      const omitFields = [];

      if (!query.valueFilter) {
        omitFields.push('showType', 'displayParameter');
      }

      return !isOmitEqual(query, oldQuery, omitFields) && !isEmpty(query);
    },

    getIntervalQuery() {
      const { interval } = this.query;

      if (!interval) {
        return {};
      }

      return convertQueryIntervalToTimestamp({
        interval,
        timezone: this.$system.timezone,
      });
    },

    updateInterval(interval) {
      this.updateQueryField('interval', interval);

      if (this.query.valueFilter) {
        const { valueFilter } = this.query;
        const { from, to } = this.getIntervalQuery();

        this.updateQueryField('valueFilter', {
          ...valueFilter,
          value: Math.min(valueFilter.value, to - from),
        });
      }
    },

    updateTrend(value) {
      this.updateContentInUserPreference({ show_trend: value });
      this.updateQueryField('showTrend', value);
    },

    updateShowType(value) {
      this.updateContentInUserPreference({ show_type: value });
      this.updateQueryField('showType', value);

      if (this.query.valueFilter) {
        this.updateQueryField('valueFilter', {
          ...this.query.valueFilter,
          method: AVAILABILITY_VALUE_FILTER_METHODS.greater,
          value: 0,
        });
      }
    },

    updateDisplayParameter(value) {
      this.updateContentInUserPreference({ display_parameter: value });
      this.updateQueryField('displayParameter', value);
    },

    updateValueFilter(value) {
      this.updateQueryField('valueFilter', value);
    },

    getQuery() {
      const {
        sortBy = [],
        sortDesc = [],
        showTrend,
        showType,
        displayParameter,
        filter,
        lockedFilter,
        valueFilter,
        itemsPerPage,
      } = this.query;

      const query = {
        ...this.interval,
        ...pick(this.query, ['page', 'search', 'active_columns', ADVANCED_SEARCH_FIELDS.search]),
        ...convertSortToRequest(
          sortBy.map(prepareAvailabilityWidgetColumnValue),
          sortDesc,
        ),
        limit: itemsPerPage,
        filters: convertFiltersToQuery(filter, lockedFilter),
      };

      if (valueFilter) {
        query.value_filter_parameter = getAvailabilityFieldByDisplayParameterAndShowType(
          displayParameter,
          showType,
        );
        query.value_filter_value = valueFilter.value;
        query.value_filter_method = valueFilter.method;
      }

      if (showTrend) {
        query.trend = getAvailabilitiesTrendByInterval(this.query.interval);
      }

      return query;
    },

    fetchList() {
      this.fetchAvailabilityList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });
    },

    getExportQuery() {
      const {
        widget_columns: widgetColumns,
        export_settings: {
          widget_export_columns: widgetExportColumns,
          export_csv_separator: exportCsvSeparator,
        } = {},
      } = this.widget.parameters;

      const columns = widgetExportColumns?.length ? widgetExportColumns : widgetColumns;
      const { displayParameter, showType } = this.query;
      const valueField = getAvailabilityFieldByDisplayParameterAndShowType(
        displayParameter,
        showType,
      );
      const isUptimeParameter = displayParameter === AVAILABILITY_DISPLAY_PARAMETERS.uptime;
      const fields = [
        { name: valueField, label: this.$t(`common.${isUptimeParameter ? 'uptime' : 'downtime'}`) },
        ...columns.map(({ value, text }) => ({ name: value, label: text })),
      ];

      const { page, limit, trend, ...restQuery } = this.getQuery();

      return {
        fields,
        separator: exportCsvSeparator,
        ...restQuery,
      };
    },

    async exportAvailabilityList() {
      this.exporting = true;

      try {
        const fileData = await this.generateFile({
          widgetId: this.widget._id,
          data: this.getExportQuery(),
        });

        this.downloadFile(getExportMetricDownloadFileUrl(fileData._id));
      } catch {
        this.$popups.error({ text: this.$t('availability.popups.exportCSVFailed') });
      } finally {
        this.exporting = false;
      }
    },
  },
};
</script>
