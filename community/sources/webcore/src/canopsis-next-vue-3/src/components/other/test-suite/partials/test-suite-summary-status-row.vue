<template>
  <c-information-block-row :label="label">
    <v-row>
      <v-col cols="2">
        {{ count }}
      </v-col>
      <v-col v-if="testsStatusPercent || testsStatusPercent === 0">
        {{ convertNumberToRoundedPercentString(testsStatusPercent, precision) }}
      </v-col>
    </v-row>
  </c-information-block-row>
</template>

<script>
import { convertNumberToRoundedPercentString } from '@/helpers/string';

export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      required: false,
    },
    total: {
      type: Number,
      required: false,
    },
    precision: {
      type: Number,
      default: 2,
    },
  },
  computed: {
    testsStatusPercent() {
      if (this.count === 0) {
        return this.count;
      }

      return this.count / this.total;
    },
  },
  methods: {
    convertNumberToRoundedPercentString,
  },
};
</script>
