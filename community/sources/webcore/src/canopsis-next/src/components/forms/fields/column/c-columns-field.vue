<template lang="pug">
  c-movable-card-iterator-field(v-field="columns", @add="add")
    template(#item="{ item, index }")
      column-field(
        v-field="columns[index]",
        :name="item.key",
        :type="type",
        :with-html="withHtml",
        :with-template="withTemplate",
        :with-color-indicator="withColorIndicator"
      )
</template>

<script>
import {
  MODALS,
  ENTITIES_TYPES,
  COLOR_INDICATOR_TYPES,
} from '@/constants';

import { widgetColumnToForm } from '@/helpers/forms/shared/widget-column';

import { formArrayMixin, formValidationHeaderMixin } from '@/mixins/form';

import ColumnField from './partials/column-field.vue';

export default {
  inject: ['$validator'],
  components: { ColumnField },
  mixins: [
    formArrayMixin,
    formValidationHeaderMixin,
  ],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    type: {
      type: String,
      default: ENTITIES_TYPES.alarm,
    },
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
      this.addItemIntoArray(widgetColumnToForm());
    },
  },
};
</script>
