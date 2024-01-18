<template>
  <v-layout column>
    <span>{{ $t('stateSetting.computeMethod') }}: {{ stateMethodName }}</span>
    <span>{{ stateMethodSummaryText }}</span>
  </v-layout>
</template>

<script>
export default {
  computed: {
    stateMethodName() {
      /**
       * TODO: Should be changed on real data
       */
      return '%rule name%';
    },

    stateMethodSummaryText() {
      /**
       * TODO: Should be changed on real data
       */
      return 'TODO: Should be implemented';
    },
  },
};
</script>
