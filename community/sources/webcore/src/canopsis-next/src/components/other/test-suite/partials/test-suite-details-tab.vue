<template>
  <v-layout>
    <v-flex xs12>
      <test-suite-details :test-suite="testSuite" />
    </v-flex>
  </v-layout>
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
