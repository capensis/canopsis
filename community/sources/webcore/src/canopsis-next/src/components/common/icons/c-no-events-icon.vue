<template lang="pug">
  v-tooltip(v-if="value", :max-width="maxWidth", :top="top", :right="right", :bottom="bottom", :left="left")
    v-icon(slot="activator", :color="color") sync_problem
    span {{ message }}
</template>

<script>
import moment from 'moment-timezone';

import { TIME_UNITS } from '@/constants';

import { convertTimestampToMoment } from '@/helpers/date/date';

export default {
  props: {
    value: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'error',
    },
    maxWidth: {
      type: [String, Number],
      required: false,
    },
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
    size: {
      type: String,
      required: false,
    },
  },
  computed: {
    message() {
      const diff = moment().diff(convertTimestampToMoment(this.value), TIME_UNITS.second);
      const duration = this.$options.filters.duration(diff);

      return this.$t('icons.noEvents', { duration });
    },
  },
};
</script>
