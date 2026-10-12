<template>
  <div class="c-mixed-field">
    <c-input-type-field
      :types="types"
      :value="inputType"
      :label="label"
      :disabled="disabled"
      :flat="flat"
      :name="name"
      :required="required"
      class="c-mixed-field__selector"
      @input="updateType"
    />
    <c-mixed-input-field
      v-field="value"
      :input-type="inputType"
      :name="name"
      :disabled="disabled"
      :flat="flat"
      :hide-details="hideDetails"
      :items="items"
      :item-text="itemText"
      :item-value="itemValue"
      :required="required"
      :types="types"
      class="ml-2 c-mixed-field__value"
    />
  </div>
</template>

<script>
import { PATTERN_FIELD_TYPES } from '@/constants';

import { convertValueByType, getFieldType } from '@/helpers/entities/pattern/form';

import { formBaseMixin } from '@/mixins/form';

export default {
  $_veeValidate: {
    name() {
      return this.name;
    },

    value() {
      return this.value;
    },
  },
  inject: ['$validator'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: [String, Number, Boolean, Array],
      default: '',
    },
    name: {
      type: String,
      default: 'value',
    },
    label: {
      type: String,
      default: null,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    flat: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    types: {
      type: Array,
      default: () => [
        { value: PATTERN_FIELD_TYPES.string },
        { value: PATTERN_FIELD_TYPES.number },
        { value: PATTERN_FIELD_TYPES.boolean },
        { value: PATTERN_FIELD_TYPES.null },
        { value: PATTERN_FIELD_TYPES.stringArray },
      ],
    },
  },
  computed: {
    inputType() {
      return getFieldType(this.value);
    },
  },
  watch: {
    types: {
      immediate: true,
      handler(types) {
        if (!types.some(({ value }) => value === this.inputType)) {
          const [type = {}] = types;

          this.updateType(type.value, type.defaultValue);
        }
      },
    },
  },
  methods: {
    updateType(type, defaultValue) {
      this.updateModel(
        convertValueByType(this.value, type, defaultValue),
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.c-mixed-field {
  display: flex;
  position: relative;

  &__selector {
    min-width: 45px;
    max-width: 45px;
  }

  &__value {
    width: 100%;
  }
}
</style>
