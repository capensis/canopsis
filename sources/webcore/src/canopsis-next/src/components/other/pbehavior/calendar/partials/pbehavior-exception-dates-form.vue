<template lang="pug">
  div
    h3.my-3.grey--text Exception dates
    v-layout(row)
      v-flex(xs12)
        pbehavior-exception-date-field(
          v-for="date in dates",
          :key="date.key",
          :value="date"
        )
    v-layout(row)
      v-flex
        v-btn Add an exception date
      v-flex
        v-btn Chose list of exceptions
</template>

<script>
import PbehaviorExceptionDateField from './pbehavior-exception-date-field.vue';

export default {
  components: { PbehaviorExceptionDateField },
  component: { PbehaviorExceptionDateField },
  props: {
    dates: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
