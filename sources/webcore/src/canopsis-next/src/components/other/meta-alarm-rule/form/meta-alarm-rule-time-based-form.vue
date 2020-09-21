<template lang="pug">
  v-layout
    v-flex(xs12)
      time-interval-field(
        v-field="timebased.time_interval",
        :intervalLabel="$t('metaAlarmRule.fields.timeInterval')"
      )
</template>

<script>
import TimeIntervalField from '@/components/forms/fields/time-interval.vue';

export default {
  inject: ['$validator'],
  components: {
    TimeIntervalField,
  },
  model: {
    prop: 'timebased',
    event: 'input',
  },
  props: {
    timebased: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

