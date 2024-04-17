<template>
  <v-layout
    class="py-2"
    column
  >
    <kpi-widget-filters
      :widget-id="widget._id"
      :user-filters="userPreference.filters"
      :widget-filters="widget.filters"
      :locked-filter="query.lockedFilter"
      :filters="query.filter"
      :interval="query.interval"
      :sampling="query.sampling"
      :min-interval-date="minAvailableDate"
      :show-interval="showInterval"
      :show-filter="showFilter"
      :filter-disabled="filterDisabled"
      :filter-addable="filterAddable"
      :filter-editable="filterEditable"
      class="mx-3"
      @update:filters="updateSelectedFilter"
      @update:interval="updateInterval"
    />
    <v-layout
      class="kpi-widget pa-3"
      column
      align-center
    >
      <h4 class="kpi-widget__title">
        {{ widget.parameters.table_title }}
      </h4>
      <c-progress-overlay
        :pending="mainRatingSettingsPending"
        transition
      />
      <c-advanced-data-table
        :items="preparedGroupMetrics"
        :loading="groupMetricsPending"
        :headers="headers"
        class="kpi-widget__table pre-line"
        no-pagination
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { pick, find } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { DATETIME_FORMATS, KPI_RATING_CRITERIA } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone } from '@/helpers/date/date';
import { convertFiltersToQuery } from '@/helpers/entities/shared/query';
import { convertMetricValueToString } from '@/helpers/entities/metric/list';
import { isCustomCriteria } from '@/helpers/entities/metric/form';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { entitiesGroupMetricsMixin } from '@/mixins/entities/group-metrics';

import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';
import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

import KpiWidgetFilters from '../../partials/kpi-widget-filters.vue';

const { mapActions } = createNamespacedHelpers('ratingSettings');

export default {
  inject: ['$system'],
  components: {
    CProgressOverlay,
    CAdvancedDataTable,
    KpiWidgetFilters,
  },
  mixins: [
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    queryIntervalFilterMixin,
    widgetPeriodicRefreshMixin,
    entitiesGroupMetricsMixin,
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
    showInterval: {
      type: Boolean,
      default: false,
    },
    showFilter: {
      type: Boolean,
      default: false,
    },
    filterAddable: {
      type: Boolean,
      default: false,
    },
    filterEditable: {
      type: Boolean,
      default: false,
    },
    filterDisabled: {
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

        ...this.widget.parameters.widgetColumns.map(({ metric, label }) => {
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
        const preparedMetrics = Object.values(data)
          .reduce((acc, value, index) => {
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
