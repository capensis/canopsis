<template>
  <v-layout>
    <v-flex xs12>
      <c-duration-field
        v-field="timebased.time_interval"
        :label="$t('metaAlarmRule.timeInterval')"
        required
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
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
