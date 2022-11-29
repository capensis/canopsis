<template lang="pug">
  v-layout(column)
    v-layout(justify-space-between, align-center)
      v-layout
        v-btn.mx-0(icon, @click="$emit('up')")
          v-icon arrow_upward
        v-btn.mx-0.ml-1(icon, @click="$emit('down')")
          v-icon arrow_downward
      v-btn.ma-0(icon, @click="$emit('remove')")
        v-icon(color="red") close
    v-text-field(
      v-field="column.label",
      v-validate="'required'",
      :placeholder="$t('common.label')",
      :error-messages="errors.collect(labelFieldName)",
      :name="labelFieldName"
    )
    v-text-field(
      v-field="column.value",
      v-validate="'required'",
      :placeholder="$t('common.value')",
      :error-messages="errors.collect(valueFieldName)",
      :name="valueFieldName"
    )
    v-layout(v-if="withTemplate", row)
      v-switch(
        :label="$t('settings.columns.withTemplate')",
        :input-value="!!column.template",
        color="primary",
        @change="enableTemplate($event)"
      )
      v-btn.primary(v-if="column.template", small, @click="showEditTemplateModal")
        span {{ $t('common.edit') }}
    v-switch(
      v-if="withHtml",
      v-field="column.isHtml",
      :label="$t('settings.columns.isHtml')",
      :disabled="!!column.template",
      color="primary"
    )
    v-switch(
      v-if="withColorIndicator",
      :label="$t('settings.colorIndicator.title')",
      :input-value="!!column.colorIndicator",
      :disabled="!!column.template",
      color="primary",
      @change="switchChangeColorIndicator($event)"
    )
    c-color-indicator-field(
      v-if="column.colorIndicator",
      v-field="column.colorIndicator",
      :value="column.colorIndicator",
      :disabled="!!column.template"
    )
</template>

<script>
import { COLOR_INDICATOR_TYPES, DEFAULT_COLUMN_TEMPLATE_VALUE, MODALS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formMixin,
  ],
  model: {
    prop: 'column',
    event: 'input',
  },
  props: {
    column: {
      type: Object,
      required: true,
    },
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
  },
  computed: {
    labelFieldName() {
      return `${this.name}.label`;
    },

    valueFieldName() {
      return `${this.name}.value`;
    },
  },
  methods: {
    enableTemplate(checked) {
      const value = checked
        ? DEFAULT_COLUMN_TEMPLATE_VALUE
        : null;

      return this.updateField('template', value);
    },

    showEditTemplateModal() {
      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          text: this.column.template,
          title: this.$t('settings.columns.withTemplate'),
          label: this.$t('common.template'),
          rules: {
            required: true,
          },
          action: value => this.updateField('template', value),
        },
      });
    },

    switchChangeColorIndicator(value) {
      return this.updateField('colorIndicator', value ? COLOR_INDICATOR_TYPES.state : null);
    },
  },
};
</script>
