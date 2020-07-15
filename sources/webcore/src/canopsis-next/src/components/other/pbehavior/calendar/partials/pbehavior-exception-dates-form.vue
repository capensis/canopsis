<template lang="pug">
  div
    h3.my-3.grey--text Exception dates
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
        v-btn.ml-0(outline, @click="addExceptionDate") Add an exception date
      v-flex
        v-btn.mr-0(outline, @click="showSelectExceptionDatesModal") Chose list of exceptions
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
