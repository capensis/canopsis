<template>
  <div :style="widgetWrapperStyles">
    <template v-if="widget.title || activeViewEditing">
      <v-card-title class="widget-title pa-2">
        <c-layout justify="space-between" align="center">
          <div>
            <h4 class="widget-title__text ml-2 font-weight-regular">
              {{ widget.title }}
            </h4>
          </div>

          <v-spacer />
        </c-layout>
      </v-card-title>

      <v-divider />
    </template>

    <v-card-text v-if="component" class="pa-0 position-relative">
      <component
        :is="component"
        v-bind="widgetSpecificProps"
        :widget="preparedWidget"
        :tab-id="tab._id"
        :visible="visible"
      />
    </v-card-text>

    <v-card-text v-else class="widget-placeholder pa-3">
      <div class="widget-placeholder__content">
        <v-icon class="widget-placeholder__icon" icon="widgets" size="32" />

        <div class="widget-placeholder__body">
          <div class="widget-placeholder__title">
            {{ placeholderTitle }}
          </div>

          <div class="widget-placeholder__meta">
            {{ widget.type }}
          </div>

          <div v-if="widget._id" class="widget-placeholder__id">
            {{ widget._id }}
          </div>
        </div>
      </div>
    </v-card-text>
  </div>
</template>

<script>
import { WIDGET_GRID_ROW_HEIGHT, WIDGET_TYPES } from '@/constants';

import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import { prepareAvailabilityWidget } from '@/helpers/entities/widget/forms/availability';
import { prepareContextWidget } from '@/helpers/entities/widget/forms/context';
import { prepareServiceWeatherWidget } from '@/helpers/entities/widget/forms/service-weather';
import { prepareStatsCalendarWidget } from '@/helpers/entities/widget/forms/stats-calendar';

import { activeViewMixin } from '@/mixins/active-view';
import { authMixin } from '@/mixins/auth';

import AlarmsListWidget from './alarm/alarms-list.vue';
import AvailabilityWidget from './availability/availability-widget.vue';
import BarChartWidget from './chart/bar-chart-widget.vue';
import LineChartWidget from './chart/line-chart-widget.vue';
import NumbersWidget from './chart/numbers-widget.vue';
import PieChartWidget from './chart/pie-chart-widget.vue';
import EntitiesListWidget from './context/entities-list.vue';
import CounterWidget from './counter/counter.vue';
import ExternalDataTableWidget from './external-data-table/external-data-table-widget.vue';
import MapWidget from './map/map.vue';
import ServiceWeatherWidget from './service-weather/service-weather.vue';
import AlarmStatisticsWidget from './statistics/alarm-statistics-widget.vue';
import UserStatisticsWidget from './statistics/user-statistics-widget.vue';
import StatsCalendarWidget from './stats/calendar/stats-calendar.vue';
import TextWidget from './text/text.vue';

export default {
  provide() {
    return {
      $widget: this.widget,
    };
  },
  components: {
    AlarmsListWidget,
    AvailabilityWidget,
    BarChartWidget,
    CounterWidget,
    EntitiesListWidget,
    ExternalDataTableWidget,
    LineChartWidget,
    MapWidget,
    NumbersWidget,
    PieChartWidget,
    ServiceWeatherWidget,
    StatsCalendarWidget,
    TextWidget,
    AlarmStatisticsWidget,
    UserStatisticsWidget,
  },
  mixins: [authMixin, activeViewMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    placeholderTitle() {
      return this.widget.title || this.$t('common.widget');
    },

    component() {
      return {
        [WIDGET_TYPES.alarmList]: AlarmsListWidget,
        [WIDGET_TYPES.barChart]: BarChartWidget,
        [WIDGET_TYPES.context]: EntitiesListWidget,
        [WIDGET_TYPES.text]: TextWidget,
        [WIDGET_TYPES.counter]: CounterWidget,
        [WIDGET_TYPES.statsCalendar]: StatsCalendarWidget,
        [WIDGET_TYPES.map]: MapWidget,
        [WIDGET_TYPES.externalDataTable]: ExternalDataTableWidget,
        [WIDGET_TYPES.lineChart]: LineChartWidget,
        [WIDGET_TYPES.pieChart]: PieChartWidget,
        [WIDGET_TYPES.numbers]: NumbersWidget,
        [WIDGET_TYPES.serviceWeather]: ServiceWeatherWidget,
        [WIDGET_TYPES.userStatistics]: UserStatisticsWidget,
        [WIDGET_TYPES.alarmStatistics]: AlarmStatisticsWidget,
        [WIDGET_TYPES.availability]: AvailabilityWidget,
      }[this.widget.type];
    },

    preparedWidget() {
      if (this.widget.type === WIDGET_TYPES.alarmList) {
        return prepareAlarmListWidget(this.widget);
      }

      if (this.widget.type === WIDGET_TYPES.context) {
        return prepareContextWidget(this.widget);
      }

      if (this.widget.type === WIDGET_TYPES.serviceWeather) {
        return prepareServiceWeatherWidget(this.widget);
      }

      if (this.widget.type === WIDGET_TYPES.statsCalendar) {
        return prepareStatsCalendarWidget(this.widget);
      }

      if (this.widget.type === WIDGET_TYPES.availability) {
        return prepareAvailabilityWidget(this.widget);
      }

      return this.widget;
    },

    widgetSpecificProps() {
      let props = {};

      if (this.activeViewIsKioskScreenMode) {
        props = {
          ...this.preparedWidget.parameters?.kiosk,
        };
      }

      if (!this.hasCurrentViewActionsAccess) {
        props = {
          ...props,
          hideActions: true,
          hideMassSelection: true,
        };
      }

      return props;
    },

    /**
     * Height is multiplying on 2 for loader correct displaying
     *
     * @returns {{minHeight: string}}
     */
    widgetWrapperStyles() {
      return { minHeight: `${WIDGET_GRID_ROW_HEIGHT * 2}px` };
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-title {
  height: 37px;
}

.widget-title__text {
  overflow: hidden;
  margin: 0;
  font-size: 1.25rem;
  line-height: 1.75rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.widget-placeholder {
  display: flex;
  min-height: 100px;
  height: 100%;
  border: 1px dashed rgba(var(--v-theme-on-surface), 0.28);
  background: rgba(var(--v-theme-primary), 0.05);
}

.widget-placeholder__content {
  display: flex;
  align-items: center;
  width: 100%;
  min-width: 0;
}

.widget-placeholder__icon {
  flex: 0 0 auto;
  color: rgba(var(--v-theme-primary), 0.62);
}

.widget-placeholder__body {
  min-width: 0;
  margin-left: 12px;
}

.widget-placeholder__title {
  overflow: hidden;
  font-weight: 500;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.widget-placeholder__meta,
.widget-placeholder__id {
  overflow: hidden;
  color: rgba(var(--v-theme-on-surface), 0.62);
  font-size: 12px;
  line-height: 18px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
