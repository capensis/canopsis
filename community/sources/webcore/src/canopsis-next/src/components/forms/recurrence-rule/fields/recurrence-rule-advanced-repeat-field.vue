<template>
  <v-select
    v-model="selectedHint"
    :label="$t('recurrenceRule.repeatOn')"
    :items="hints"
  />
</template>

<script>
import { RRule } from 'rrule';

import { isSeveralEqual } from '@/helpers/collection';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    start: {
      type: Date,
      required: false,
    },
  },
  computed: {
    selectedHint: {
      set(value) {
        this.updateModel({
          ...this.value,
          ...value,
        });
      },
      get() {
        return this.hints.find(({ value }) => isSeveralEqual(value, this.value, Object.keys(value)));
      },
    },

    hints() {
      if (!this.start) {
        return [];
      }

      const hints = [];

      const dayOfMonth = this.start.getDate();

      if (dayOfMonth < 28) {
        hints.push({
          text: this.$t('recurrenceRule.dayOfMonth', { day: dayOfMonth }),
          value: {
            byyearday: '',
            byweekno: '',
            byhour: '',
            bysetpos: '',
            byweekday: [],
            bymonthday: String(dayOfMonth),
          },
        });
      }

      /**
       * 7 - day in week
       * @type {number}
       */
      const weekNumber = Math.ceil(this.start.getDate() / 7);
      const weekDayIndex = this.start.getDay();

      const weekDay = [
        this.$t('common.weekDays.monday'),
        this.$t('common.weekDays.tuesday'),
        this.$t('common.weekDays.wednesday'),
        this.$t('common.weekDays.thursday'),
        this.$t('common.weekDays.friday'),
        this.$t('common.weekDays.saturday'),
        this.$t('common.weekDays.sunday'),
      ][weekDayIndex];
      const weekDayValue = [
        RRule.MO.weekday,
        RRule.TU.weekday,
        RRule.WE.weekday,
        RRule.TH.weekday,
        RRule.FR.weekday,
        RRule.SA.weekday,
        RRule.SU.weekday,
      ][weekDayIndex];

      const weekNumberString = [
        this.$t('common.ordinals.first'),
        this.$t('common.ordinals.second'),
        this.$t('common.ordinals.third'),
        this.$t('common.ordinals.fourth'),
        this.$t('common.ordinals.fifth'),
      ][weekNumber - 1];

      hints.push({
        text: this.$t('recurrenceRule.weekDayOfMonth', { weekNumber: weekNumberString, weekDay }),
        value: {
          byyearday: '',
          bymonthday: '',
          byweekno: '',
          byhour: '',
          bysetpos: String(weekDayValue),
          byweekday: [RRule.WE.weekday],
        },
      });

      return hints;
    },
  },
};
</script>
