<template lang="pug">
  component.c-mixed-input-field(
    :is="inputComponent.is",
    v-validate="rules",
    v-bind="inputComponent.bind",
    v-on="inputComponent.on"
  )
</template>

<script>
import { isNull, pick } from 'lodash';

import { PATTERN_FIELD_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: [String, Number, Boolean, Array],
      default: '',
    },
    inputType: {
      type: String,
      default: PATTERN_FIELD_TYPES.string,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'value',
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
      default: () => [],
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    errorMessages() {
      return this.errors.collect(this.name);
    },

    rules() {
      return {
        required: this.required,
      };
    },

    switchLabel() {
      return String(this.value);
    },

    isInputTypeText() {
      return [
        PATTERN_FIELD_TYPES.number,
        PATTERN_FIELD_TYPES.string,
      ].includes(this.inputType);
    },

    inputComponent() {
      if (this.isInputTypeText) {
        const additionalProps = this.items.length
          ? { ...pick(this, ['items', 'itemText', 'itemValue']), returnObject: false, forceSearching: true }
          : {};

        return {
          is: this.items.length ? 'v-combobox' : 'v-text-field',
          bind: {
            ...pick(this, [
              'value',
              'name',
              'disabled',
              'hideDetails',
              'flat',
              'errorMessages',
              'label',
            ]),
            ...additionalProps,

            class: 'c-mixed-input-field__text',
            type: this.inputType === PATTERN_FIELD_TYPES.number ? 'number' : 'text',
            singleLine: true,
            dense: true,
          },
          on: {
            input: this.updateTextFieldValue,
            'update:searchInput': this.updateTextFieldValue,
          },
        };
      }

      if (this.inputType === PATTERN_FIELD_TYPES.boolean) {
        return {
          is: 'v-switch',

          bind: {
            class: 'ma-0 c-mixed-input-field__switch',
            name: this.name,
            inputValue: this.value,
            label: this.switchLabel,
            disabled: this.disabled,
            color: 'primary',
            hideDetails: true,
          },
          on: {
            change: this.updateModel,
          },
        };
      }

      if (this.inputType === PATTERN_FIELD_TYPES.stringArray) {
        return {
          is: 'c-array-text-field',

          bind: {
            name: this.name,
            values: this.value,
            disabled: this.disabled,
            errorMessages: this.errorMessages,
          },
          on: {
            change: this.updateModel,
          },
        };
      }

      return {
        is: 'v-text-field',

        bind: {
          name: this.name,
          errorMessages: this.errorMessages,
          value: 'null',
          disabled: true,
        },
      };
    },
  },
  methods: {
    updateTextFieldValue(value) {
      let preparedValue = value;

      if (isNull(value) && this.inputType !== PATTERN_FIELD_TYPES.null) {
        preparedValue = '';
      }

      if (this.inputType === PATTERN_FIELD_TYPES.number) {
        preparedValue = Number(value);
      }

      this.updateModel(preparedValue);
    },
  },
};
</script>

<style lang="scss">
.c-mixed-input-field {
  &__switch {
    padding: 18px 0;

    & .v-label {
      text-transform: capitalize;
    }
  }
}
</style>
