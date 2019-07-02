<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsAnnotationLine.title') }}
    v-container(fluid)
      v-layout(row, wrap)
        v-flex(xs12)
          v-switch(
          :input-value="annotationLine.enabled",
          :label="$t('settings.statsAnnotationLine.enabled')",
          @change="updateField('enabled', $event)"
          )
        v-flex(xs12)
          v-text-field(
          :value="annotationLine.value",
          :label="$t('settings.statsAnnotationLine.value')",
          :disabled="!annotationLine.enabled",
          type="number",
          :name="nameField"
          :error-messages="getCollectedErrorMessages(nameField)",
          v-validate="'numeric'"
          @input="updateField('value', $event)"
          )
        v-flex(xs12)
          v-text-field(
          :value="annotationLine.label",
          :label="$t('settings.statsAnnotationLine.label')",
          :disabled="!annotationLine.enabled",
          @input="updateField('label', $event)"
          )
        v-flex(xs12)
          v-btn(
          :style="{ backgroundColor: annotationLine.lineColor }",
          :disabled="!annotationLine.enabled",
          @click="showColorPickerModal('lineColor')"
          ) {{ $t('settings.statsAnnotationLine.pickLineColor') }}
          v-btn(
          :style="{ backgroundColor: annotationLine.labelColor }",
          :disabled="!annotationLine.enabled",
          @click="showColorPickerModal('labelColor')"
          ) {{ $t('settings.statsAnnotationLine.pickLabelColor') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formMixin from '@/mixins/form';

export default {
  mixins: [modalMixin, formMixin],
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
  },
  data() {
    return {
      nameField: 'Annotation_Line_Value',
    };
  },
  computed: {
    getCollectedErrorMessages() {
      return (name) => {
        if (this.errors) {
          return this.errors.collect(name);
        }

        return [];
      };
    },
  },
  methods: {
    showColorPickerModal(key) {
      this.showModal({
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

