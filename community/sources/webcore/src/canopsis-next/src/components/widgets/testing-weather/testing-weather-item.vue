<template>
  <card-with-see-alarms-btn
    class="white--text cursor-pointer ma-1"
    :style="{ backgroundColor: color }"
    :show-button="hasAlarmsListAccess"
    tile
    @click="showTestSuiteInformationModal"
    @show:alarms="showAlarmListModal"
  >
    <v-layout
      class="fill-height"
    >
      <v-flex class="pa-2">
        <h3 class="text-md-center">
          {{ testSuite.name }}
        </h3>
        <v-divider
          class="white"
          light
        />
        <v-layout
          class="pt-1"
          justify-start
        >
          <v-flex xs6>
            <v-layout
              class="fill-height"
              column
              justify-space-between
            >
              <c-mini-bar-chart
                :history="testSuite.mini_chart"
                :unit="$constants.TIME_UNITS.second"
              />
              <span class="pre-wrap">{{ testSuite.timestamp | date('testSuiteFormat') }}</span>
            </v-layout>
          </v-flex>
          <v-flex xs6>
            <test-suite-statistics :test-suite="testSuite" />
          </v-flex>
        </v-layout>
      </v-flex>
    </v-layout>
  </card-with-see-alarms-btn>
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
