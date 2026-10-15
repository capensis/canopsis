<template>
  <v-row class="c-infos-attribute-field" density="compact">
    <v-col v-bind="fieldColProps">
      <c-combobox
        v-if="combobox"
        :model-value="value.dictionary"
        :items="items"
        :name="dictionaryName"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :loading="pending"
        :hide-details="row"
        item-title="value"
        item-value="value"
        return-object
        required
        @update:model-value="updateInfosDictionary"
      />
      <c-text-field
        v-else
        :model-value="value.dictionary"
        :name="dictionaryName"
        :disabled="disabled"
        :label="label || $t('common.dictionary')"
        :hide-details="row"
        required
        @update:model-value="updateInfosDictionary"
      />
    </v-col>
    <v-col v-bind="fieldColProps">
      <c-select
        :model-value="value.field"
        :items="fieldItems"
        :name="fieldName"
        :disabled="disabled || !value.dictionary"
        :label="label || $t('common.field')"
        :hide-details="row"
        item-title="text"
        item-value="value"
        required
        @update:model-value="updateInfosField"
      />
    </v-col>
  </v-row>
</template>

<script>
import { keyBy } from 'lodash';

import { PATTERN_FIELD_TYPES, PATTERN_RULE_INFOS_FIELDS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  emits: ['update:value'],
  mixins: [formMixin],
  props: {
    value: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'infos',
    },
    divider: {
      type: String,
      default: '.',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    row: {
      type: Boolean,
      default: false,
    },
    column: {
      type: Boolean,
      required: false,
    },
    pending: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    dictionaryName() {
      return `${this.name}${this.divider}dictionary`;
    },

    fieldName() {
      return `${this.name}${this.divider}field`;
    },

    fieldItems() {
      return [
        {
          text: this.$t('common.name'),
          value: PATTERN_RULE_INFOS_FIELDS.name,
        },
        {
          text: this.$t('common.value'),
          value: PATTERN_RULE_INFOS_FIELDS.value,
        },
      ];
    },

    fieldColProps() {
      if (this.row) {
        return { cols: 6 };
      }

      if (this.column) {
        return { cols: 12 };
      }

      return {};
    },

    itemsByValue() {
      return keyBy(this.items, 'value');
    },
  },
  methods: {
    /**
     * Updates `value.dictionary` and related type fields based on the provided infos.
     *
     * When a combobox item is selected, `infos` is an object with `value` and optional
     * `definedType`. When a plain text is entered, `infos` is a string.
     *
     * @param {Object|string} [infos={}]
     * @param {string} [infos.value] Dictionary name when an object is provided
     * @param {string} [infos.definedType] Field defined type coming from the dictionary item
     * @returns {Object}
     */
    updateInfosDictionary(infos = {}) {
      return this.updateModel({
        ...this.value,

        dictionary: infos?.value ?? infos,
        fieldType: infos?.definedType ?? PATTERN_FIELD_TYPES.string,
      });
    },

    /**
     * Updates `value.field` and keeps `definedType` only when field is `value`.
     *
     * @param {'name'|'value'} field Field key from `PATTERN_RULE_INFOS_FIELDS`
     * @returns {Object}
     */
    updateInfosField(field) {
      return this.updateModel({
        ...this.value,

        field,
        fieldType:
          field === PATTERN_RULE_INFOS_FIELDS.value
            ? (this.itemsByValue[this.value.dictionary]?.definedType ?? PATTERN_FIELD_TYPES.string)
            : PATTERN_FIELD_TYPES.string,
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.c-infos-attribute-field {
  width: 100%;

  :deep(.v-col) {
    min-width: 0;
  }
}
</style>
