<template lang="pug">
  v-card.white--text.cursor-pointer.weather-item.ma-1(
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
              span.pre-wrap {{ testSuite.timestamp | date('testSuiteFormat') }}
          v-flex(xs6)
            test-suite-statistics(:test-suite="testSuite")
    v-btn.see-alarms-btn(
      v-if="hasAlarmsListAccess",
      flat,
      @click.stop="showAlarmListModal"
    ) {{ $t('serviceWeather.seeAlarms') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import VRuntimeTemplate from 'v-runtime-template';

import {
  MODALS,
  USERS_PERMISSIONS,
  TEST_SUITE_COLORS,
} from '@/constants';

import { authMixin } from '@/mixins/auth';

import TestSuiteStatistics from './test-suite-statistics.vue';

const { mapActions } = createNamespacedHelpers('alarm');

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
    ...mapActions({
      fetchComponentAlarmsListWithoutStore: 'fetchComponentAlarmsListWithoutStore',
    }),

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
      try {
        this.$modals.show({
          name: MODALS.alarmsList,
          config: {
            title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.testSuite.name }),
            fetchList: params => this.fetchComponentAlarmsListWithoutStore({
              params: { ...params, _id: this.testSuite.entity_id },
            }),
          },
        });
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
