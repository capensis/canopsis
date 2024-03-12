<template>
  <c-information-block-row :label="label">
    <v-layout>
      <v-flex xs2>
        {{ count }}
      </v-flex>
      <v-flex v-if="testsStatusPercent">
        {{ testsStatusPercent | percentage(precision) }}
      </v-flex>
    </v-layout>
  </c-information-block-row>
</template>

<script>
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
};
</script>
