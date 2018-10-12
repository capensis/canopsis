<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.statsSelect')}}
    v-container
      v-card.mb-2
        v-container.pt-0(fluid)
          v-select(
            v-model="form.stat"
            :items="statsTypes",
            @change="statChange",
          )
          v-switch(:label="$t('common.trend')", v-model="form.trend", @change="trendChange")
</template>

<script>
import { STATS_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      form: {
        stat: 'alarms_created',
        trend: 'true',
      },
    };
  },
  computed: {
    /**
     * Get stats different types from constant, and return an object with stat's value and stat's translated title
     */
    statsTypes() {
      return Object.values(STATS_TYPES)
        .map(item => ({ value: item.value, text: this.$t(`stats.types.${item.value}`), options: item.options }));
    },
  },
  mounted() {
    this.form = { ...this.value };
  },
  methods: {
    statChange(event) {
      this.$emit('input', { ...this.form, stat: event });
    },
    trendChange(event) {
      this.$emit('input', { ...this.form, trend: event });
    },
  },
};
</script>
