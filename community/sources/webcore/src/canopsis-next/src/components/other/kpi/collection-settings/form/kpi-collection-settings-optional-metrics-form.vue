<template>
  <v-list
    class="kpi-collection-settings-metrics-form"
    subheader
  >
    <v-subheader>{{ $t('kpi.collectionSetting.optionalMetrics') }}</v-subheader>
    <v-list-item>
      <v-list-item-action>
        <v-checkbox
          v-field="value.enabled_instructions"
          color="primary"
        >
          <template #label="">
            <span class="kpi-collection-settings-metrics-form__label ml-4">
              {{ $t('kpi.collectionSetting.manualInstructions') }}
            </span>
          </template>
        </v-checkbox>
      </v-list-item-action>
    </v-list-item>
    <v-list-item>
      <v-list-item-action>
        <v-checkbox
          v-field="value.enabled_not_acked_metrics"
          color="primary"
        >
          <template #label="">
            <span class="kpi-collection-settings-metrics-form__label ml-4">
              {{ $t('kpi.collectionSetting.notAckedMetrics') }}
            </span>
          </template>
        </v-checkbox>
      </v-list-item-action>
    </v-list-item>
  </v-list>
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>

<style lang="scss">
.kpi-collection-settings-metrics-form {
  &__label {
    color: #fff;

    .theme--light & {
      color: #000;
    }
  }
}
</style>
