<template>
  <div>
    <c-layout class="pbehavior-recurrence-rule-periods__header" align="center" justify="center">
      <span class="text-h6">
        {{ t('common.periods') }}
      </span>

      <div class="pbehavior-recurrence-rule-periods__actions">
        <c-action-btn
          :tooltip="t('pbehavior.periodsCalendar')"
          icon="event_note"
          @click="showPbehaviorRecurrenceRuleModal"
        />
      </div>
    </c-layout>

    <c-select
      :model-value="selectedRange"
      :items="availableRanges"
      :label="t('common.range')"
      class="mt-2"
      item-title="text"
      item-value="value"
      @update:model-value="updateSelectedRange"
    />

    <v-progress-linear
      v-show="pending"
      :height="3"
      class="ma-0"
      color="primary"
      indeterminate
    />

    <v-data-table
      :headers="headers"
      :items="timespans"
      :items-per-page="itemsPerPage"
      :items-per-page-options="itemsPerPageOptions"
      :loading="pending"
      :no-data-text="t('common.noData')"
      class="mt-2"
    >
      <template #item.from="{ item }">
        {{ dateToString(item.from) }}
      </template>

      <template #item.to="{ item }">
        {{ dateToString(item.to) }}
      </template>
    </v-data-table>
  </div>
</template>

<script>
import { MODALS, PBEHAVIOR_RRULE_PERIODS_RANGES, TIME_UNITS } from '@/constants';

import {
  addUnitToDate,
  convertDateToEndOfUnitTimestamp,
  convertDateToStartOfUnitTimestamp,
  convertDateToString,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useInfo } from '@/hooks/store/modules/info';
import { usePbehaviorTimespan } from '@/hooks/store/modules/pbehavior-timespan';

export default {
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
      itemsPerPage: 5,
      itemsPerPageOptions: [5, 10, 25],
      selectedRange: PBEHAVIOR_RRULE_PERIODS_RANGES.thisWeek,
    };
  },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { timezone } = useInfo();
    const { fetchPbehaviorTimespansListWithoutStore } = usePbehaviorTimespan();

    return {
      t,
      modals,
      popups,
      timezone,
      fetchPbehaviorTimespansListWithoutStore,
    };
  },

  computed: {
    availableRanges() {
      return Object.values(PBEHAVIOR_RRULE_PERIODS_RANGES).map(value => ({
        value,
        text: this.t(`recurrenceRule.periodsRanges.${value}`),
      }));
    },

    headers() {
      return [
        { title: this.t('common.begin'), key: 'from', sortable: false },
        { title: this.t('common.end'), key: 'to', sortable: false },
      ];
    },
  },

  watch: {
    pbehavior() {
      return this.fetchList();
    },
  },

  mounted() {
    return this.fetchList();
  },

  methods: {
    /**
     * Get object with start and end fields by selected rrules periods range key.
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

    updateSelectedRange(range) {
      this.selectedRange = range;

      return this.fetchList();
    },

    showPbehaviorRecurrenceRuleModal() {
      this.modals.show({
        name: MODALS.pbehaviorRecurrenceRule,
        config: {
          pbehavior: this.pbehavior,
        },
      });
    },

    dateToString(date) {
      return convertDateToString(date);
    },

    /**
     * Fetch timespans list.
     *
     * @returns {Promise<void>}
     */
    async fetchList() {
      try {
        this.pending = true;

        const { pbehavior } = this;
        const { start, end } = this.getStartAndStopForSelectedRange();

        const from = convertDateToTimestampByTimezone(start, this.timezone);
        const to = convertDateToTimestampByTimezone(end, this.timezone);

        const data = pbehaviorToTimespanRequest({
          pbehavior,
          from,
          to,
        });

        this.timespans = await this.fetchPbehaviorTimespansListWithoutStore({ data });
      } catch (err) {
        console.error(err);

        this.popups.error({ text: err.message || this.t('errors.default') });
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.pbehavior-recurrence-rule-periods {
  &__header {
    position: relative;
  }

  &__actions {
    position: absolute;
    right: 0;
  }
}
</style>
