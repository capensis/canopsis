<template lang="pug">
  div
    v-layout(justify-space-between, align-center)
      v-layout(justify-center)
        span.title {{ $t('common.periods') }}
      c-action-btn(
        :tooltip="$t('pbehavior.periodsCalendar')",
        icon="event_note",
        @click="showPbehaviorRecurrenceRuleModal"
      )
    v-layout.mt-2
      v-select(
        v-model="selectedRange",
        :items="availableRanges",
        :label="$t('common.range')",
        @change="fetchList"
      )
    v-layout
      v-data-iterator.data-iterator(
        :items="timespans",
        :loading="pending"
      )
        template(#header="")
          v-flex
            v-fade-transition
              v-progress-linear.progress.ma-0(
                v-show="pending",
                :height="3",
                color="primary",
                indeterminate
              )
        template(#item="{ item }")
          v-flex
            v-card
              v-card-title {{ item.from | date }} — {{ item.to | date }}
        v-flex(#no-data="")
          v-card
            v-card-title {{ $t('common.noData') }}
</template>

<script>
import { MODALS, PBEHAVIOR_RRULE_PERIODS_RANGES, TIME_UNITS } from '@/constants';

import {
  addUnitToDate,
  convertDateToEndOfUnitTimestamp,
  convertDateToStartOfUnitTimestamp,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';

import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

export default {
  inject: ['$system'],
  mixins: [entitiesPbehaviorTimespansMixin],
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      timespans: [],
      selectedRange: PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek,
    };
  },
  computed: {
    availableRanges() {
      return Object.values(PBEHAVIOR_RRULE_PERIODS_RANGES)
        .map(value => ({ value, text: this.$t(`recurrenceRule.periodsRanges.${value}`) }));
    },
  },
  watch: {
    pbehavior() {
      this.fetchList();
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    /**
     * Get object with start and end fields by selected rrules periods range key
     *
     * @returns {{ start: number, end: number }}
     */
    getStartAndStopForSelectedRange() {
      switch (this.selectedRange) {
        case PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek:
          return {
            start: convertDateToStartOfUnitTimestamp(Date.now(), 'isoWeek'),
            end: convertDateToEndOfUnitTimestamp(Date.now(), 'isoWeek'),
          };

        case PBEHAVIOR_RRULE_PERIODS_RANGES.nextWeek: {
          const nextWeek = addUnitToDate(Date.now(), 1, TIME_UNITS.week);

          return {
            start: convertDateToStartOfUnitTimestamp(nextWeek, 'isoWeek'),
            end: convertDateToEndOfUnitTimestamp(nextWeek, 'isoWeek'),
          };
        }

        case PBEHAVIOR_RRULE_PERIODS_RANGES.next2Weeks:
          return {
            start: convertDateToStartOfUnitTimestamp(
              addUnitToDate(Date.now(), 1, TIME_UNITS.week),
              'isoWeek',
            ),
            end: convertDateToEndOfUnitTimestamp(
              addUnitToDate(Date.now(), 2, TIME_UNITS.week),
              'isoWeek',
            ),
          };

        case PBEHAVIOR_RRULE_PERIODS_RANGES.thisMonth:
          return {
            start: convertDateToStartOfUnitTimestamp(Date.now(), TIME_UNITS.month),
            end: convertDateToEndOfUnitTimestamp(Date.now(), TIME_UNITS.month),
          };

        case PBEHAVIOR_RRULE_PERIODS_RANGES.nextMonth:
          return {
            start: convertDateToStartOfUnitTimestamp(
              addUnitToDate(Date.now(), 1, TIME_UNITS.month),
              TIME_UNITS.month,
            ),
            end: convertDateToEndOfUnitTimestamp(
              addUnitToDate(Date.now(), 1, TIME_UNITS.month),
              TIME_UNITS.month,
            ),
          };
      }

      throw new Error('Incorrect range');
    },

    showPbehaviorRecurrenceRuleModal() {
      this.$modals.show({
        name: MODALS.pbehaviorRecurrenceRule,
        config: {
          pbehavior: this.pbehavior,
        },
      });
    },

    /**
     * Fetch timespans list
     *
     * @returns {Promise<void>}
     */
    async fetchList() {
      try {
        this.pending = true;

        const { pbehavior } = this;
        const { start, end } = this.getStartAndStopForSelectedRange();

        const from = convertDateToTimestampByTimezone(start, this.$system.timezone);
        const to = convertDateToTimestampByTimezone(end, this.$system.timezone);

        const data = pbehaviorToTimespanRequest({
          pbehavior,
          from,
          to,
        });

        this.timespans = await this.fetchTimespansListWithoutStore({
          data,
        });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || this.$t('errors.default') });
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss">
.data-iterator {
  position: relative;
  width: 100%;
  padding-top: 3px;

  .progress {
    position: absolute;
    top: 0;
  }
}
</style>
