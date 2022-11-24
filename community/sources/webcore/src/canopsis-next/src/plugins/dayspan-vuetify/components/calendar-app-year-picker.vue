<template lang="pug">
  v-autocomplete(v-field="value", :items="years", :menu-props="menuProps")
</template>

<script>
import { range } from 'lodash';

export default {
  props: {
    value: {
      type: Number,
      required: true,
    },
    currentYear: {
      type: Number,
      default: () => new Date().getFullYear(),
    },
    previousYearsOffset: {
      type: Number,
      default: 10,
    },
    nextYearsOffset: {
      type: Number,
      default: 50,
    },
  },
  computed: {
    menuProps() {
      return {
        auto: true,
      };
    },

    years() {
      const start = this.currentYear - this.previousYearsOffset;
      const end = this.currentYear + this.nextYearsOffset;

      return range(start, end);
    },
  },
};
</script>
