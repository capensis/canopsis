<template lang="pug">
  div
    h3.my-3.grey--text {{ $t('pbehaviorExceptions.title') }}
    v-divider
    v-layout.mt-3(row)
      v-flex(xs12)
        pbehavior-exception-date-field(
          v-for="(date, index) in dates",
          v-field="dates[index]",
          :key="date.key",
          @delete="removeItemFromArray(index)"
        )
    v-layout(row)
      v-flex
        v-btn.ml-0(outline, @click="addExceptionDate") {{ $t('pbehaviorExceptions.create') }}
      v-flex
        v-btn.mr-0(outline, @click="showSelectExceptionDatesModal") {{ $t('pbehaviorExceptions.choose') }}
</template>

<script>
import moment from 'moment';

import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import formArrayMixin from '@/mixins/form/array';

import PbehaviorExceptionDateField from './pbehavior-exception-date-field.vue';

export default {
  components: { PbehaviorExceptionDateField },
  mixins: [formArrayMixin],
  model: {
    prop: 'dates',
    event: 'input',
  },
  props: {
    dates: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    showSelectExceptionDatesModal() {
      this.$modals.show({
        name: MODALS.selectExceptionsDatesLists,
        config: {
          action: (exceptions) => {
            const preparedExceptionDates = exceptions.reduce((acc, { exdates }) => {
              acc.push(...exdates);

              return acc;
            }, []);
            preparedExceptionDates.push(...this.dates);

            this.$emit('input', preparedExceptionDates);
          },
        },
      });
    },
    addExceptionDate() {
      const startOfTodayMoment = moment().startOf('day');

      this.addItemIntoArray({
        key: uid(),
        begin: startOfTodayMoment.toDate(),
        end: startOfTodayMoment.toDate(),
        type: '',
      });
    },
  },
};
</script>
