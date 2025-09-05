<template>
  <v-card color="grey">
    <v-card-text>
      <v-layout column>
        <span class="text-subtitle-2">{{ $t('templateTesting.testData') }}</span>
        <v-select
          v-field="items"
          :items="items"
          :label="$t('templateTesting.testData')"
          name="testData"
        />
      </v-layout>
    </v-card-text>
  </v-card>
</template>

<script>
export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
