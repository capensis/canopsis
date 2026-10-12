<template>
  <v-row>
    <v-col cols="12">
      <test-suite-details :test-suite="testSuite" />
    </v-col>
  </v-row>
</template>

<script>
import TestSuiteDetails from './test-suite-details.vue';

export default {
  components: { TestSuiteDetails },
  props: {
    testSuite: {
      type: Object,
      required: true,
    },
  },
};
</script>
