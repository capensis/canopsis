<template>
  <card-with-see-alarms-btn
    :style="{ backgroundColor: color }"
    :show-button="hasAlarmsListAccess"
    class="text-white cursor-pointer ma-1"
    tile
    @click="showTestSuiteInformationModal"
    @show:alarms="showAlarmListModal"
  >
    <c-layout class="fill-height" fill>
      <div class="pa-2">
        <h3 class="text-center">
          {{ testSuite.name }}
        </h3>
        <v-divider color="white" />
        <v-row class="pt-1">
          <v-col cols="6">
            <c-layout class="fill-height" justify="space-between" column>
              <c-mini-bar-chart
                :history="testSuite.mini_chart"
                :unit="$constants.TIME_UNITS.second"
              />
              <span class="pre-wrap">
                {{ convertDateToString(testSuite.timestamp, 'testSuiteFormat') }}
              </span>
            </c-layout>
          </v-col>
          <v-col cols="6">
            <test-suite-statistics :test-suite="testSuite" />
          </v-col>
        </v-row>
      </div>
    </c-layout>
  </card-with-see-alarms-btn>
</template>

<script>
import { mapActions } from 'pinia';

import { MODALS, TEST_SUITE_COLORS, USER_PERMISSIONS } from '@/constants';

import { useAlarmStore } from '@/stores/modules/alarm';

import { convertDateToString } from '@/helpers/date/date';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { authMixin } from '@/mixins/auth';

import CardWithSeeAlarmsBtn from '@/components/common/card/card-with-see-alarms-btn.vue';

import TestSuiteStatistics from './test-suite-statistics.vue';

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
      return this.checkAccess(USER_PERMISSIONS.business.testingWeather.actions.alarmsList);
    },

    color() {
      return TEST_SUITE_COLORS[this.testSuite.state];
    },
  },
  methods: {
    ...mapActions(useAlarmStore, {
      fetchComponentAlarmsListWithoutStore: 'fetchComponentAlarmsListWithoutStore',
    }),

    convertDateToString,

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
