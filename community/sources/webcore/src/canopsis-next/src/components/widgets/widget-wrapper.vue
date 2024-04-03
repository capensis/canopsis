<template>
  <div :style="widgetWrapperStyles">
    <template v-if="widget.title || editing">
      <v-card-title class="widget-title pa-2">
        <v-layout
          justify-space-between
          align-center
        >
          <v-flex>
            <h4 class="ml-2 font-weight-regular">
              {{ widget.title }}
            </h4>
          </v-flex>
          <v-spacer />
        </v-layout>
      </v-card-title>
      <v-divider />
    </template>
    <v-card-text class="pa-0 position-relative">
      <component
        v-bind="widgetProps"
        :is="widgetProps.is"
        :widget="preparedWidget"
        :tab-id="tab._id"
        :visible="visible"
      />
    </v-card-text>
  </div>
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { WIDGET_TYPES, WIDGET_TYPES_RULES, WIDGET_GRID_ROW_HEIGHT, COMPONENTS_BY_WIDGET_TYPES } from '@/constants';

import featuresService from '@/services/features';

import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import { prepareContextWidget } from '@/helpers/entities/widget/forms/context';
import { prepareServiceWeatherWidget } from '@/helpers/entities/widget/forms/service-weather';
import { prepareStatsCalendarAndCounterWidget } from '@/helpers/entities/widget/forms/stats-calendar';
import { prepareMapWidget } from '@/helpers/entities/widget/forms/map';

import AlarmsListWidget from './alarm/alarms-list.vue';
import EntitiesListWidget from './context/entities-list.vue';
import ServiceWeatherWidget from './service-weather/service-weather.vue';
import TestingWeatherWidget from './testing-weather/testing-weather.vue';
import StatsCalendarWidget from './stats/calendar/stats-calendar.vue';
import TextWidget from './text/text.vue';
import CounterWidget from './counter/counter.vue';
import MapWidget from './map/map.vue';
import BarChartWidget from './chart/bar-chart-widget.vue';
import LineChartWidget from './chart/line-chart-widget.vue';
import PieChartWidget from './chart/pie-chart-widget.vue';
import NumbersWidget from './chart/numbers-widget.vue';
import UserStatisticsWidget from './statistics/user-statistics-widget.vue';
import AlarmStatisticsWidget from './statistics/alarm-statistics-widget.vue';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  components: {
    AlarmsListWidget,
    EntitiesListWidget,
    ServiceWeatherWidget,
    TestingWeatherWidget,
    StatsCalendarWidget,
    TextWidget,
    CounterWidget,
    MapWidget,
    BarChartWidget,
    LineChartWidget,
    PieChartWidget,
    NumbersWidget,
    UserStatisticsWidget,
    AlarmStatisticsWidget,

    ...featuresService.get('components.widgetWrapper.components', {}),
  },
  props: {
    widget: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
    kiosk: {
      type: Boolean,
      default: false,
    },
    visible: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    ...mapGetters(['edition']),

    /**
     * Height is multiplying on 2 for loader correct displaying
     *
     * @returns {{minHeight: string}}
     */
    widgetWrapperStyles() {
      return { minHeight: `${WIDGET_GRID_ROW_HEIGHT * 2}px` };
    },

    preparedWidget() {
      switch (this.widget.type) {
        case WIDGET_TYPES.alarmList:
          return prepareAlarmListWidget(this.widget);

        case WIDGET_TYPES.context:
          return prepareContextWidget(this.widget);

        case WIDGET_TYPES.serviceWeather:
          return prepareServiceWeatherWidget(this.widget);

        case WIDGET_TYPES.statsCalendar:
        case WIDGET_TYPES.counter:
          return prepareStatsCalendarAndCounterWidget(this.widget);

        case WIDGET_TYPES.map:
          return prepareMapWidget(this.widget);
      }

      const preparer = featuresService.get('components.widgetWrapper.computed.preparedWidget');

      return preparer ? preparer.call(this) : this.widget;
    },

    widgetProps() {
      const { type } = this.widget;
      const widgetComponentsMap = { ...COMPONENTS_BY_WIDGET_TYPES };
      let widgetSpecificsProp = {};

      if (this.kiosk) {
        widgetSpecificsProp = {
          ...this.widget.parameters.kiosk,
        };
      }

      Object.entries(WIDGET_TYPES_RULES).forEach(([key, rule]) => {
        if (rule.edition !== this.edition) {
          widgetComponentsMap[key] = 'c-alert-overlay';
          widgetSpecificsProp = {
            ...widgetSpecificsProp,

            message: this.$t('errors.statsWrongEditionError'),
            value: true,
          };
        }
      });

      const component = widgetComponentsMap[type];

      if (!component) {
        return {
          is: 'c-alert-overlay',
          message: this.$t('errors.unknownWidgetType', { type }),
          value: true,
        };
      }

      return {
        ...widgetSpecificsProp,

        is: component,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.widget-title {
  height: 37px;
}
</style>
