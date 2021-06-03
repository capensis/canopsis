<template lang="pug">
  v-card.white--text.cursor-pointer.weather__item.ma-1(
    :style="{ backgroundColor: color }",
    :class="{ 'v-card__with-see-alarms-btn': hasAlarmsListAccess }",
    tile,
    @click="showTestSuiteInformationModal"
  )
    v-layout.fill-height(row)
      v-flex.pa-2
        h3.text-md-center {{ testSuite.name }}
        v-divider.white(light)
        v-layout.pt-1(justify-start)
          v-flex(xs6)
            v-layout.fill-height(column, justify-space-between)
              c-mini-bar-chart(:history="testSuite.mini_chart", :unit="$constants.TIME_UNITS.second")
              span.pre-wrap {{ testSuite.timestamp | date('testSuiteFormat', true) }}
          v-flex(xs6)
            test-suite-statistics(:test-suite="testSuite")
    v-btn.see-alarms-btn(
      v-if="hasAlarmsListAccess",
      flat,
      @click.stop="showAlarmListModal"
    ) {{ $t('serviceWeather.seeAlarms') }}
</template>

<script>
import VRuntimeTemplate from 'v-runtime-template';

import {
  MODALS,
  USERS_PERMISSIONS,
  WIDGET_TYPES,
  TEST_SUITE_COLORS,
} from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';

import { authMixin } from '@/mixins/auth';

import TestSuiteStatistics from './test-suite-statistics.vue';

export default {
  components: {
    VRuntimeTemplate,
    TestSuiteStatistics,
  },
  mixins: [authMixin],
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasAlarmsListAccess() {
      return this.checkAccess(USERS_PERMISSIONS.business.testingWeather.actions.alarmsList);
    },

    color() {
      return TEST_SUITE_COLORS[this.testSuite.state];
    },
  },
  methods: {
    showTestSuiteInformationModal() {
      this.$modals.show({
        name: MODALS.testSuite,
        config: {
          testSuite: this.testSuite,
          color: this.color,
        },
      });
    },

    showAlarmListModal() {
      const widget = generateWidgetByType(WIDGET_TYPES.alarmList);

      const testSuiteFilter = {
        title: this.testSuite.name,
        filter: { $and: [{ 'entity.impact': this.testSuite._id }] },
      };

      widget.parameters = {
        ...widget.parameters,
        mainFilter: testSuiteFilter,
        viewFilters: [testSuiteFilter],
      };

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
        },
      });
    },
  },
};
</script>