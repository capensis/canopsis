<template lang="pug">
  div Duration
</template>

<script>
export default {
  props: {
    value: {
      type: String,
    },
  },
  data() {
    return {
      units: {
        hour: {
          text: 'Hour',
          value: 'h',
        },
        day: {
          text: 'Day',
          value: 'd',
        },
        week: {
          text: 'Week',
          value: 'w',
        },
        month: {
          text: 'Month',
          value: 'm',
        },
      },
    };
  },
  computed: {
    durationValue() {
      // TODO : Return duration value using the value's string
      return null;
    },
    durationUnit() {
      // TODO : Return duration unit using the value's string
      return null;
    },
  },
};
</script>
