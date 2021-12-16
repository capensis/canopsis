<template lang="pug">
  v-list-group(data-test="statsAnnotationLine")
    v-list-tile(slot="activator") {{ $t('settings.statsAnnotationLine.title') }}
    v-container(fluid)
      v-layout(row, wrap)
        v-flex(xs12)
          v-switch(
            v-field="annotationLine.enabled",
            :label="$t('settings.statsAnnotationLine.enabled')",
            color="primary",
            data-test="annotationEnabled"
          )
        v-flex(xs12)
          v-text-field(
            v-field="annotationLine.value",
            v-validate="'numeric'",
            :label="$t('settings.statsAnnotationLine.value')",
            :disabled="!annotationLine.enabled",
            :name="valueName",
            :error-messages="errors ? errors.collect(valueName) : []",
            type="number",
            data-test="annotationValue"
          )
        v-flex(xs12)
          v-text-field(
            v-field="annotationLine.label",
            :label="$t('settings.statsAnnotationLine.label')",
            :disabled="!annotationLine.enabled",
            data-test="annotationLabel"
          )
        v-flex(xs12)
          c-color-picker-field(
            v-field="annotationLine.lineColor",
            :label="$t('settings.statsAnnotationLine.pickLineColor')",
            :disabled="!annotationLine.enabled"
          )
          c-color-picker-field(
            v-field="annotationLine.labelColor",
            :label="$t('settings.statsAnnotationLine.pickLabelColor')",
            :disabled="!annotationLine.enabled"
          )
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'annotationLine',
    event: 'input',
  },
  props: {
    annotationLine: {
      type: Object,
      default: () => ({}),
    },
    valueName: {
      type: String,
      default: 'annotationLine.value',
    },
  },
};
</script>
