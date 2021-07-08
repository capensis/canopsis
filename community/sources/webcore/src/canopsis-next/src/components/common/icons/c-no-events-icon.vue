<template lang="pug">
  v-tooltip(v-if="value")
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
