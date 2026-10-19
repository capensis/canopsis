<template>
  <div>
    <calendar-pagination-arrow-btn
      :label="prevLabel"
      icon="keyboard_arrow_left"
      @click="$emit('prev')"
    />
    <calendar-period-picker
      :type="type"
      :focus="focus"
      @input="selectPeriod"
    />
    <calendar-pagination-arrow-btn
      :label="nextLabel"
      icon="keyboard_arrow_right"
      @click="$emit('next')"
    />
  </div>
</template>

<script>
import { lowerCase } from 'lodash';

import CalendarPaginationArrowBtn from '@/components/common/calendar/partials/calendar-pagination-arrow-btn.vue';

import CalendarPeriodPicker from './calendar-period-picker.vue';

export default {
  components: { CalendarPaginationArrowBtn, CalendarPeriodPicker },
  props: {
    focus: {
      type: Date,
      required: false,
    },
    type: {
      type: String,
      required: false,
    },
  },
  computed: {
    typeString() {
      return lowerCase(this.$t(`calendar.${this.type}`));
    },

    prevLabel() {
      return [
        this.$t('common.previous'),
        this.typeString,
      ].join(' ');
    },

    nextLabel() {
      return [
        this.$t('common.next'),
        this.typeString,
      ].join(' ');
    },
  },
  methods: {
    selectPeriod(value) {
      this.$emit('update:focus', value);
    },
  },
};
</script>
