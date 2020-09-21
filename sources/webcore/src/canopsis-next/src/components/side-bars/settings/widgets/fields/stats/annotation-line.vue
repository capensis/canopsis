<template lang="pug">
  v-list-group(data-test="statsAnnotationLine")
    v-list-tile(slot="activator") {{ $t('settings.statsAnnotationLine.title') }}
    v-container(fluid)
      v-layout(row, wrap)
        v-flex(xs12)
          v-switch(
            v-field="annotationLine.enabled",
            data-test="annotationEnabled",
            :label="$t('settings.statsAnnotationLine.enabled')",
            color="primary"
          )
        v-flex(xs12)
          v-text-field(
            v-field="annotationLine.value",
            v-validate="'numeric'",
            data-test="annotationValue",
            :label="$t('settings.statsAnnotationLine.value')",
            :disabled="!annotationLine.enabled",
            :name="valueName",
            :error-messages="errors ? errors.collect(valueName) : []",
            type="number"
          )
        v-flex(xs12)
          v-text-field(
            v-field="annotationLine.label",
            data-test="annotationLabel",
            :label="$t('settings.statsAnnotationLine.label')",
            :disabled="!annotationLine.enabled"
          )
        v-flex(xs12)
          v-btn(
            data-test="annotationLineColorButton",
            :style="{ backgroundColor: annotationLine.lineColor }",
            :disabled="!annotationLine.enabled",
            @click="showColorPickerModal('lineColor')"
          ) {{ $t('settings.statsAnnotationLine.pickLineColor') }}
          v-btn(
            data-test="annotationLabelColorButton",
            :style="{ backgroundColor: annotationLine.labelColor }",
            :disabled="!annotationLine.enabled",
            @click="showColorPickerModal('labelColor')"
          ) {{ $t('settings.statsAnnotationLine.pickLabelColor') }}
</template>

<script>
import { MODALS } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'annotationLine',
    event: 'input',
  },
  inject: ['$validator'],
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
  methods: {
    showColorPickerModal(key) {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.annotationLine[key],
          action: color => this.updateField(key, color),
        },
      });
    },
  },
};
</script>

