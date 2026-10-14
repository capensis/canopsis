<template>
  <v-list class="kpi-collection-settings-metrics-form">
    <v-list-subheader>
      {{ $t('kpi.collectionSetting.optionalMetrics') }}
    </v-list-subheader>
    <v-list-item>
      <template #prepend>
        <v-checkbox
          :model-value="value.enabled_instructions"
          color="primary"
          hide-details
          @update:model-value="updateField('enabled_instructions', $event)"
        />
      </template>
      <v-list-item-title class="kpi-collection-settings-metrics-form__label">
        {{ $t('kpi.collectionSetting.manualInstructions') }}
      </v-list-item-title>
    </v-list-item>
    <v-list-item>
      <template #prepend>
        <v-checkbox
          :model-value="value.enabled_not_acked_metrics"
          color="primary"
          hide-details
          @update:model-value="updateField('enabled_not_acked_metrics', $event)"
        />
      </template>
      <v-list-item-title class="kpi-collection-settings-metrics-form__label">
        {{ $t('kpi.collectionSetting.notAckedMetrics') }}
      </v-list-item-title>
    </v-list-item>
    <v-list-item>
      <template #prepend>
        <v-checkbox
          :model-value="value.enabled_sli_metrics"
          color="primary"
          hide-details
          @update:model-value="updateField('enabled_sli_metrics', $event)"
        />
      </template>
      <v-list-item-title class="kpi-collection-settings-metrics-form__label">
        {{ $t('kpi.collectionSetting.sliDuration') }}
      </v-list-item-title>
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

  emits: ['update:value'],

  methods: {
    updateField(field, value) {
      this.$emit('update:value', {
        ...this.value,
        [field]: value,
      });
    },
  },
};
</script>

<style lang="scss">
.kpi-collection-settings-metrics-form {
  &__label {
    color: rgb(var(--v-theme-on-surface));
  }
}
</style>
