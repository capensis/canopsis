<template lang="pug">
  card-with-see-alarms-btn.white--text.cursor-pointer.ma-1(
    :style="{ backgroundColor: color }",
    :show-button="hasAlarmsListAccess",
    tile,
    @click="showTestSuiteInformationModal",
    @show:alarms="showAlarmListModal"
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
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, USERS_PERMISSIONS, TEST_SUITE_COLORS } from '@/constants';

import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { authMixin } from '@/mixins/auth';

import CardWithSeeAlarmsBtn from '@/components/common/card/card-with-see-alarms-btn.vue';

import TestSuiteStatistics from './test-suite-statistics.vue';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  components: {
    CardWithSeeAlarmsBtn,
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
            widget: generatePreparedDefaultAlarmListWidget(),
            title: this.$t('modals.alarmsList.prefixTitle', { prefix: this.testSuite.name }),
            fetchList: params => this.fetchComponentAlarmsListWithoutStore({
              params: { ...params, _id: this.testSuite.entity_id },
            }),
          },
        });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
</script>
