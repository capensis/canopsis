<template lang="pug">
  v-tooltip(
    v-if="value",
    :max-width="maxWidth",
    :top="top",
    :right="right",
    :bottom="bottom",
    :left="left",
    @input="showTooltipHandler"
  )
    v-icon(slot="activator", :color="color", :size="size") sync_problem
    span {{ message }}
</template>

<script>
import moment from 'moment-timezone';

import { TIME_UNITS } from '@/constants';

import { convertTimestampToMoment } from '@/helpers/date/date';

import { createPollingMixin } from '@/mixins/polling';

export default {
  mixins: [
    createPollingMixin({
      method: 'setMomentNow',
      delay: 1000,
    }),
  ],
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
  data() {
    return {
      momentNow: moment(),
    };
  },
  computed: {
    durationDiff() {
      return this.momentNow.diff(convertTimestampToMoment(this.value), TIME_UNITS.second);
    },

    message() {
      const duration = this.$options.filters.duration(this.durationDiff);

      return this.$t('icons.noEvents', { duration });
    },
  },
  methods: {
    showTooltipHandler(value) {
      this.setMomentNow();

      if (value) {
        this.startPolling();
      } else {
        this.startPolling();
      }
    },

    setMomentNow() {
      this.momentNow = moment();
    },
  },
};
</script>
