<template lang="pug">
  v-container
    v-card.my-2(
      v-for="(column, index) in columns",
      :key="`settings-column-${index}`"
    )
      v-layout.pt-2(justify-space-between)
        v-flex(xs3)
          v-layout.text-xs-center.pl-2(justify-space-between)
            v-flex(xs1)
              v-btn(icon, @click.prevent="up(index)")
                v-icon arrow_upward
            v-flex(xs5)
              v-btn(icon, @click.prevent="down(index)")
                v-icon arrow_downward
        v-flex.d-flex(xs3)
          div.text-xs-right.pr-2
            v-btn(icon, @click.prevent="removeItemFromArray(index)")
              v-icon(color="error") close
      v-layout(justify-center, wrap)
        v-flex(xs11)
          v-text-field(
            v-validate="'required'",
            :placeholder="$t('common.label')",
            :error-messages="errors.collect(`label[${index}]`)",
            :name="`label[${index}]`",
            :value="column.label",
            @input="updateFieldInArrayItem(index, 'label', $event)"
          )
        v-flex(xs11)
          v-text-field(
            v-validate="'required'",
            :placeholder="$t('common.value')",
            :error-messages="errors.collect(`value[${index}]`)",
            :value="column.value",
            :name="`value[${index}]`",
            @input="updateFieldInArrayItem(index, 'value', $event)"
          )
        v-flex(v-if="withTemplate", xs11)
          v-layout(row)
            v-switch(
              :label="$t('settings.columns.withTemplate')",
              :input-value="!!column.template",
              color="primary",
              @change="enableTemplate(index, $event)"
            )
            v-btn.primary(v-if="column.template", small, @click="showEditTemplateModal(index)")
              span {{ $t('common.edit') }}
        v-flex(v-if="withHtml", xs11)
          v-switch(
            :label="$t('settings.columns.isHtml')",
            :input-value="column.isHtml",
            :disabled="!!column.template",
            color="primary",
            @change="updateFieldInArrayItem(index, 'isHtml', $event)"
          )
        v-flex(v-if="withColorIndicator", xs11)
          v-switch(
            :label="$t('settings.colorIndicator.title')",
            :input-value="!!column.colorIndicator",
            :disabled="!!column.template",
            color="primary",
            @change="switchChangeColorIndicator(index, $event)"
          )
          v-layout(v-if="column.colorIndicator", row)
            c-color-indicator-field(
              :value="column.colorIndicator",
              :disabled="!!column.template",
              @input="updateFieldInArrayItem(index, 'colorIndicator', $event)"
            )
    v-btn.ml-0(color="primary", @click.prevent="add") {{ $t('common.add') }}
</template>

<script>
import { COLOR_INDICATOR_TYPES, DEFAULT_COLUMN_TEMPLATE_VALUE, MODALS } from '@/constants';

import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    columns: {
      type: [Array, Object],
      default: () => [],
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
  },
  methods: {
    enableTemplate(index, checked) {
      const value = checked
        ? DEFAULT_COLUMN_TEMPLATE_VALUE
        : null;

      return this.updateFieldInArrayItem(index, 'template', value);
    },

    showEditTemplateModal(index) {
      const column = this.columns[index];

      this.$modals.show({
        name: MODALS.textEditor,
        config: {
          text: column.template,
          title: this.$t('settings.columns.withTemplate'),
          label: this.$t('common.template'),
          rules: {
            required: true,
          },
          action: value => this.updateFieldInArrayItem(index, 'template', value),
        },
      });
    },

    switchChangeColorIndicator(index, value) {
      return this.updateFieldInArrayItem(index, 'colorIndicator', value ? COLOR_INDICATOR_TYPES.state : null);
    },

    add() {
      const column = { label: '', value: '' };

      if (this.withHtml) {
        column.isHtml = false;
      }

      this.addItemIntoArray(column);
    },
    up(index) {
      if (index > 0) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index - 1];
        columns[index - 1] = temp;

        this.$emit('input', columns);
      }
    },
    down(index) {
      if (index < this.columns.length - 1) {
        const columns = [...this.columns];
        const temp = columns[index];

        columns[index] = columns[index + 1];
        columns[index + 1] = temp;

        this.$emit('input', columns);
      }
    },
  },
};
</script>
