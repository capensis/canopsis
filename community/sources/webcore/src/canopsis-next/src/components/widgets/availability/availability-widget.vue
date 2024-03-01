<template>
  <v-layout
    class="gap-3"
    column
  >
    <availability-widget-filters
      :type="query.showType"
      :display-parameter="query.displayParameter"
      :trend="query.showTrend"
      :interval="query.interval"
      :widget-id="widget._id"
      :user-filters="userPreference.filters"
      :widget-filters="widget.filters"
      :locked-filter="lockedFilter"
      :filters="mainFilter"
      :show-interval="hasAccessToInterval"
      :show-filter="hasAccessToListFilters"
      :filter-addable="hasAccessToAddFilter"
      :filter-editable="hasAccessToEditFilter"
      :min-interval-date="minAvailableDate"
      :exporting="exporting"
      class="px-3 pt-3"
      @export="exportAvailabilityList"
      @update:filters="updateSelectedFilter"
      @update:interval="updateInterval"
      @update:trend="updateTrend"
      @update:type="updateShowType"
      @update:display-parameter="updateDisplayParameter"
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
      :options.sync="options"
    />
  </v-layout>
</template>

<script>
import { pick } from 'lodash';

import { getAvailabilityDownloadFileUrl } from '@/helpers/entities/availability/url';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';
import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';

import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { permissionsWidgetsAvailabilityFilters } from '@/mixins/permissions/widgets/availability/filters';
import { permissionsWidgetsAlarmStatisticsInterval } from '@/mixins/permissions/widgets/availability/interval';
import { exportMixinCreator } from '@/mixins/widget/export';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { entitiesAvailabilityMixin } from '@/mixins/entities/availability';
import { widgetOptionsMixin } from '@/mixins/widget/options';

import AvailabilityWidgetFilters from '@/components/widgets/availability/partials/availability-widget-filters.vue';
import AvailabilityList from '@/components/other/availability/availability-list.vue';

export default {
  components: { AvailabilityList, AvailabilityWidgetFilters },
  mixins: [
    widgetPeriodicRefreshMixin,
    widgetFilterSelectMixin,
    widgetFetchQueryMixin,
    permissionsWidgetsAvailabilityFilters,
    permissionsWidgetsAlarmStatisticsInterval,
    queryIntervalFilterMixin,
    entitiesAvailabilityMixin,
    widgetOptionsMixin,
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
  },
  data() {
    return {
      exporting: false,
    };
  },
  computed: {
    minAvailableDate() {
      const { min_date: minDate } = this.availabilitiesMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },
  },
  methods: {
    customQueryCondition(query, oldQuery) {
      const fields = ['showTrend', 'interval', 'filter'];

      return isMetricsQueryChanged(
        pick(query, fields),
        pick(oldQuery, fields),
        this.minAvailableDate,
      );
    },

    updateTrend(value) {
      this.updateContentInUserPreference({ show_trend: value });
      this.updateQueryField('showTrend', value);
    },

    updateShowType(value) {
      this.updateContentInUserPreference({ show_type: value });
      this.updateQueryField('showType', value);
    },

    updateDisplayParameter(value) {
      this.updateContentInUserPreference({ display_parameter: value });
      this.updateQueryField('displayParameter', value);
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['page']),
        with_history: this.query.showTrend,
        widget_filters: convertFiltersToQuery(this.query.filter, this.lockedFilter),
      };
    },

    async fetchList() {
      await this.fetchAvailabilityList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });
    },

    getExportQuery() {
      /**
       * TODO: Fix it, when API will be integrated
       */
      return {};
    },

    async exportAvailabilityList() {
      this.exporting = true;

      try {
        const fileData = await this.generateFile({
          widgetId: this.widget._id,
          data: this.getExportQuery(),
        });

        this.downloadFile(getAvailabilityDownloadFileUrl(fileData._id));
      } catch (err) {
        this.$popups.error({ text: this.$t('availability.popups.exportCSVFailed') });
      } finally {
        this.exporting = false;
      }
    },
  },
};
</script>
