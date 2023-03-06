<template lang="pug">
  div(:style="widgetWrapperStyles")
    template(v-if="widget.title || editing")
      v-card-title.widget-title.pa-2
        v-layout(justify-space-between, align-center)
          v-flex
            h4.ml-2.font-weight-regular {{ widget.title }}
          v-spacer
      v-divider
    v-card-text.pa-0.position-relative
      component(
        v-bind="widgetProps",
        :widget="preparedWidget",
        :tab-id="tab._id",
        :editing="editing"
      )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import {
  WIDGET_TYPES,
  WIDGET_TYPES_RULES,
  WIDGET_GRID_ROW_HEIGHT,
  COMPONENTS_BY_WIDGET_TYPES,
} from '@/constants';

import {
  prepareAlarmListWidget,
  prepareContextWidget,
  prepareServiceWeatherWidget,
  prepareStatsCalendarAndCounterWidget,
  prepareMapWidget,
} from '@/helpers/widgets';

import featuresService from '@/services/features';

import AlarmsListWidget from './alarm/alarms-list.vue';
import EntitiesListWidget from './context/entities-list.vue';
import ServiceWeatherWidget from './service-weather/service-weather.vue';
import TestingWeatherWidget from './testing-weather/testing-weather.vue';
import StatsCalendarWidget from './stats/calendar/stats-calendar.vue';
import TextWidget from './text/text.vue';
import CounterWidget from './counter/counter.vue';
import MapWidget from './map/map.vue';

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

      Object.entries(WIDGET_TYPES_RULES).forEach(([key, rule]) => {
        if (rule.edition !== this.edition) {
          widgetComponentsMap[key] = 'c-alert-overlay';
          widgetSpecificsProp = {
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
