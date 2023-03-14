<template lang="pug">
  v-list.kpi-collection-settings-metrics-form(subheader)
    v-subheader {{ $t('kpi.collectionSetting.optionalMetrics') }}
    v-list-tile
      v-list-tile-action
        v-checkbox(v-field="value.enabled_manual_instructions", color="primary")
          template(#label="")
            span.kpi-collection-settings-metrics-form__label.ml-4 {{ $t('kpi.collectionSetting.manualInstructions') }}
    v-list-tile
      v-list-tile-action
        v-checkbox(v-field="value.enabled_not_acked_metrics", color="primary")
          template(#label="")
            span.kpi-collection-settings-metrics-form__label.ml-4 {{ $t('kpi.collectionSetting.notAckedMetrics') }}
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
