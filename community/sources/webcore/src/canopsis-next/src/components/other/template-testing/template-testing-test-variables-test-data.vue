<template>
  <v-card class="test-data-card">
    <v-card-text>
      <v-layout column>
        <span class="text-subtitle-2">{{ $t('templateTesting.testData') }}</span>
        <div>
          <v-select
            v-field="items"
            :items="items"
            :label="$t('templateTesting.testData')"
            name="testData"
          />
        </div>
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

<style lang="scss" scoped>
.theme--light .test-data-card {
  background-color: var(--v-application-background-darken1);
}

.theme--dark .test-data-card {
  background-color: var(--v-application-background-lighten2);
}
</style>
