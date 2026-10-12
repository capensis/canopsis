<template>
  <c-layout class="py-2" column>
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreferenceFilters"
      :widget-filters="widget.filters"
      :locked-filters="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
      :sampling="query.sampling"
      :min-interval-date="minAvailableDate"
      :show-interval="showInterval"
      :show-filter="showFilter"
      class="mx-3"
      @update:filters="updateSelectedFilter"
      @update:interval="updateInterval"
    />
    <c-layout class="kpi-widget pa-3" align="center" column>
      <h4 class="kpi-widget__title">
        {{ widget.parameters.table_title }}
      </h4>
      <c-progress-overlay :pending="mainRatingSettingsPending" transition />
      <c-advanced-data-table
        :items="preparedGroupMetrics"
        :loading="groupMetricsPending"
        :headers="headers"
        class="kpi-widget__table pre-line"
        no-pagination
      />
    </c-layout>
  </c-layout>
</template>

<script>
import { find, pick } from 'lodash';
import { mapActions, mapState } from 'pinia';

import { DATETIME_FORMATS, KPI_RATING_CRITERIA } from '@/constants';

import { useMetricsStore } from '@/stores/modules/metrics';
import { useRatingSettingsStore } from '@/stores/modules/rating-settings';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertMetricValueToString } from '@/helpers/entities/metric/list';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';

import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import KpiWidgetFilters from '@/components/widgets/partials/kpi-widget-filters.vue';

export default {
  inject: ['$system'],

  components: {
    KpiWidgetFilters,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetPeriodicRefreshMixin,
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
    type: {
      type: Number,
      required: true,
    },
    showInterval: {
      type: Boolean,
      default: false,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
  },

  data() {
    return {
      mainRatingSettingsPending: false,
      mainRatingSettings: [],
    };
  },

  computed: {
    ...mapState(useMetricsStore, {
      getGroupMetricsListByWidgetId: 'getGroupMetricsListByWidgetId',
      getGroupMetricsPendingByWidgetId: 'getGroupMetricsPendingByWidgetId',
      getGroupMetricsMetaByWidgetId: 'getGroupMetricsMetaByWidgetId',
    }),

    groupMetrics() {
      return this.getGroupMetricsListByWidgetId(this.widget._id);
    },

    groupMetricsPending() {
      return this.getGroupMetricsPendingByWidgetId(this.widget._id);
    },

    groupMetricsMeta() {
      return this.getGroupMetricsMetaByWidgetId(this.widget._id);
    },

    mainParameter() {
      return this.widget.parameters?.mainParameter ?? {};
    },

    firstColumnText() {
      const criteriaObject = find(this.mainRatingSettings, {
        id: this.mainParameter.criteria,
      });

      const labelsMap = {
        [KPI_RATING_CRITERIA.user]: this.$t('common.username'),
        [KPI_RATING_CRITERIA.role]: this.$t('common.role'),
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

        ...(this.widget.parameters.widgetColumns ?? []).map(({ metric, label }) => {
          const kpiKey = `kpi.statisticsWidgets.metrics.${metric}`;
          const alarmKey = `alarm.metrics.${metric}`;

          let text;

          if (label) {
            text = label;
          } else if (this.$te(kpiKey)) {
            text = this.$t(kpiKey);
          } else if (this.$te(alarmKey)) {
            text = this.$t(alarmKey);
          }

          return {
            text: text ?? this.$t(`user.metrics.${metric}`),
            value: metric,
            sortable: false,
          };
        }),
      ];
    },

    preparedGroupMetrics() {
      if (!this.query.parameters?.length) {
        return [];
      }

      return this.groupMetrics.map(({ title, data = [] }) => {
        const preparedMetrics = Object.values(data).reduce((acc, value, index) => {
          const { metric } = this.query.parameters[index];

          acc[metric] = value.reduce((secondAcc, item) => {
            const preparedValue = convertMetricValueToString({
              value: item.value,
              metric,
              format: DATETIME_FORMATS.refreshFieldFormat,
            });

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

    userPreferenceFilters() {
      return this.userPreference?.filters ?? [];
    },
  },

  watch: {
    minAvailableDate() {
      const { from } = this.getIntervalQuery();

      if (this.minAvailableDate && from < this.minAvailableDate) {
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
    ...mapActions(useMetricsStore, {
      fetchGroupMetricsList: 'fetchGroupMetricsList',
    }),
    ...mapActions(useRatingSettingsStore, {
      fetchRatingSettingsWithoutStore: 'fetchListWithoutStore',
    }),

    getQuery() {
      const { filter, lockedFilter } = this.query;

      return {
        ...this.getIntervalQuery(),
        ...pick(this.query, ['parameters', 'criteria', 'entity_patterns', 'limit', 'page']),
        widget_filters: convertFiltersToQuery(filter, lockedFilter),
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
