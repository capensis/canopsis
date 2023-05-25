<template lang="pug">
  v-layout.py-2(column)
    statistics-widget-filters.mx-3(
      :widget-id="widget._id",
      :user-filters="userPreference.filters",
      :widget-filters="widget.filters",
      :locked-value="lockedFilter",
      :filters="mainFilter",
      :interval="query.interval",
      :sampling="query.sampling",
      :min-interval-date="minAvailableDate",
      :show-filter="hasAccessToUserFilter",
      :show-interval="hasAccessToInterval",
      :show-sampling="hasAccessToSampling",
      :filter-disabled="!hasAccessToListFilters",
      :filter-addable="hasAccessToAddFilter",
      :filter-editable="hasAccessToEditFilter",
      @update:filters="updateSelectedFilter",
      @update:sampling="updateSampling",
      @update:interval="updateInterval"
    )
    v-layout(column)
      c-progress-overlay(:pending="mainRatingSettingsPending", transition)
      c-advanced-data-table.pre-line(
        :items="preparedGroupMetrics",
        :loading="groupMetricsPending",
        :headers="headers",
        :total-items="groupMetricsMeta.total_count",
        no-pagination
      )
      c-table-pagination(
        :total-items="groupMetricsMeta.total_count",
        :rows-per-page.sync="query.limit",
        :page.sync="query.page"
      )
</template>

<script>
import { pick, find } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { KPI_RATING_CRITERIA } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertFilterToQuery } from '@/helpers/query';
import { convertMetricValueToString, isCustomCriteria } from '@/helpers/metrics';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { metricsIntervalFilterMixin } from '@/mixins/widget/metrics/interval';
import { widgetSamplingFilterMixin } from '@/mixins/widget/chart/sampling';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesGroupMetricsMixin } from '@/mixins/entities/group-metrics';
import { permissionsWidgetsBarChartInterval } from '@/mixins/permissions/widgets/chart/bar/interval';
import { permissionsWidgetsBarChartSampling } from '@/mixins/permissions/widgets/chart/bar/sampling';
import { permissionsWidgetsBarChartFilters } from '@/mixins/permissions/widgets/chart/bar/filters';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

import StatisticsWidgetFilters from './partials/statistics-widget-filters.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$system'],
  components: {
    CProgressOverlay,
    CAdvancedDataTable,
    StatisticsWidgetFilters,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    metricsIntervalFilterMixin,
    widgetSamplingFilterMixin,
    widgetPeriodicRefreshMixin,
    entitiesGroupMetricsMixin,
    permissionsWidgetsBarChartInterval,
    permissionsWidgetsBarChartSampling,
    permissionsWidgetsBarChartFilters,
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
      mainRatingSettingsPending: false,
      mainRatingSettings: [],
    };
  },
  computed: {
    mainParameter() {
      return this.widget.parameters?.mainParameter ?? {};
    },

    isCustomCriteria() {
      return isCustomCriteria(this.mainParameter.criteria);
    },

    firstColumnText() {
      const criteriaObject = find(this.mainRatingSettings, {
        id: this.mainParameter.criteria,
      });

      const labelsMap = {
        [KPI_RATING_CRITERIA.user]: this.$t('common.username'),
        [KPI_RATING_CRITERIA.role]: this.$tc('common.role'),
        [KPI_RATING_CRITERIA.category]: this.$t('common.category'),
        [KPI_RATING_CRITERIA.impactLevel]: this.$t('common.impactLevel'),
      };

      return this.mainParameter.columnName || labelsMap[criteriaObject?.label] || '';
    },

    headers() {
      if (this.mainRatingSettingsPending) {
        return [];
      }

      return [
        {
          text: this.firstColumnText,
          value: 'title',
          sortable: false,
        },

        ...this.widget.parameters.widgetColumns.map(({ metric }) => {
          const alarmKey = `alarm.metrics.${metric}`;

          return {
            text: this.$te(alarmKey) ? this.$t(alarmKey) : this.$t(`user.metrics.${metric}`),
            value: metric,
            sortable: false,
          };
        }),
      ];
    },

    preparedGroupMetrics() {
      return this.groupMetrics.map(({ title, data = [] }) => {
        const preparedMetrics = Object.entries(data)
          .reduce((acc, [metric, value]) => {
            acc[metric] = value.reduce((secondAcc, item) => {
              const preparedValue = convertMetricValueToString(item.value, metric);

              return secondAcc + (item.title ? `${item.title}: ${preparedValue}\n` : preparedValue);
            }, '');

            return acc;
          }, {});

        return {
          title,

          ...preparedMetrics,
        };
      });
    },

    minAvailableDate() {
      const { min_date: minDate } = this.groupMetricsMeta;

      return minDate
        ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone)
        : null;
    },
  },
  watch: {
    minAvailableDate() {
      const { from } = this.getIntervalQuery();

      if (from < this.minAvailableDate) {
        this.query = {
          ...this.query,
          interval: {
            from: this.minAvailableDate,
            to: this.query.interval.to,
          },
        };
      }
    },
  },
  mounted() {
    this.fetchMainRatingSettingsList();
  },
  methods: {
    ...mapActions({ fetchRatingSettingsWithoutStore: 'fetchListWithoutStore' }),

    getQuery() {
      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'criteria', 'entity_patterns', 'limit', 'page']),
        widget_filters: convertFilterToQuery(this.query.filter),
      };
    },

    async fetchMainRatingSettingsList() {
      this.mainRatingSettingsPending = true;

      const { data: mainRatingSettings = [] } = await this.fetchRatingSettingsWithoutStore({
        params: {
          type: this.type,
          main: true,
          paginate: false,
        },
      });

      this.mainRatingSettings = mainRatingSettings;
      this.mainRatingSettingsPending = false;
    },

    fetchList() {
      return this.fetchGroupMetricsList({
        widgetId: this.widget._id,
        params: this.getQuery(),
      });
    },
  },
};
</script>
