<template lang="pug">
  v-list-group(data-test="columnNames")
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ label }}
    v-container
      v-card.my-2(
        data-test="columnName",
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
                v-icon(color="red") close
        v-layout(justify-center, wrap)
          v-flex(xs11)
            v-text-field(
              data-test="columnNameLabel",
              v-validate="'required'",
              :placeholder="$t('common.label')",
              :error-messages="errors.collect(`label[${index}]`)",
              :name="`label[${index}]`",
              :value="column.label",
              @input="updateFieldInArrayItem(index, 'label', $event)"
            )
          v-flex(xs11)
            v-text-field(
              data-test="columnNameValue",
              v-validate="'required'",
              :placeholder="$t('common.value')",
              :error-messages="errors.collect(`value[${index}]`)",
              :value="column.value",
              :name="`value[${index}]`",
              @input="updateFieldInArrayItem(index, 'value', $event)"
            )
          v-flex(v-if="withHtml", xs11)
            v-switch(
              data-test="columnNameSwitch",
              :label="$t('settings.columns.isHtml')",
              :input-value="column.isHtml",
              color="primary",
              @change="updateFieldInArrayItem(index, 'isHtml', $event)"
            )
          v-flex(v-if="withColorIndicator", xs11)
            v-switch(
              :label="$t('settings.colorIndicator.title')",
              :input-value="!!column.colorIndicator",
              color="primary",
              @change="switchChangeColorIndicator(index, $event)"
            )
            v-layout(v-if="column.colorIndicator", row)
              color-indicator-field(
                :value="column.colorIndicator",
                @input="updateFieldInArrayItem(index, 'colorIndicator', $event)"
              )
      v-btn(
        data-test="columnNameAddButton",
        color="primary",
        @click.prevent="add"
      ) {{ $t('common.add') }}
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

import formArrayMixin from '@/mixins/form/array';
import formValidationHeaderMixin from '@/mixins/form/validation-header';

import ColorIndicatorField from '../partials/color-indicator-field.vue';

export default {
  inject: ['$validator'],
  components: { ColorIndicatorField },
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    columns: {
      type: [Array, Object],
      default: () => [],
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
