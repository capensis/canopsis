<template>
  <component
    :is="inputComponent.is"
    v-validate="rules"
    class="c-mixed-input-field"
    v-bind="inputComponent.bind"
    v-on="inputComponent.on"
  />
</template>

<script>
import { computed } from 'vue';
import { isNull, pick } from 'lodash';

import { PATTERN_FIELD_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
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
  setup(props, { emit }) {
    const { errors } = useValidator();
    const { updateModel } = useModelField(props, emit);

    const rules = computed(() => ({ required: props.required }));
    const switchLabel = computed(() => String(props.value));

    const isInputTypeText = computed(() => [
      PATTERN_FIELD_TYPES.number,
      PATTERN_FIELD_TYPES.string,
    ].includes(props.inputType));

    const updateTextFieldValue = (value) => {
      let preparedValue = value;

      if (isNull(value) && props.inputType !== PATTERN_FIELD_TYPES.null) {
        preparedValue = '';
      }

      if (props.inputType === PATTERN_FIELD_TYPES.number) {
        preparedValue = Number(value);
      }

      updateModel(preparedValue);
    };

    const inputComponent = computed(() => {
      if (isInputTypeText.value) {
        const name = `${props.name}.combobox`;
        const additionalProps = props.items.length
          ? { ...pick(props, ['items', 'itemText', 'itemValue']), returnObject: false, forceSearching: true }
          : {};

        return {
          is: props.items.length ? 'v-combobox' : 'v-text-field',
          bind: {
            ...pick(props, [
              'value',
              'disabled',
              'hideDetails',
              'flat',
            ]),
            ...additionalProps,

            name,
            errorMessages: errors.collect(name),
            placeholder: props.label,
            class: 'c-mixed-input-field__text',
            type: props.inputType === PATTERN_FIELD_TYPES.number ? 'number' : 'text',
            singleLine: true,
          },
          on: {
            input: updateTextFieldValue,
            'update:search-input': updateTextFieldValue,
          },
        };
      }

      if (props.inputType === PATTERN_FIELD_TYPES.boolean) {
        const name = `${props.name}.switch`;

        return {
          is: 'v-switch',
          bind: {
            class: 'ma-0 c-mixed-input-field__switch',
            name,
            inputValue: props.value,
            label: switchLabel.value,
            disabled: props.disabled,
            color: 'primary',
            hideDetails: true,
          },
          on: {
            change: updateModel,
          },
        };
      }

      if (props.inputType === PATTERN_FIELD_TYPES.stringArray) {
        const name = `${props.name}.array_text_field`;

        return {
          is: 'c-array-text-field',
          bind: {
            name,
            values: props.value,
            disabled: props.disabled,
            errorMessages: errors.collect(name),
          },
          on: {
            change: updateModel,
          },
        };
      }

      const name = `${props.name}.text_field`;

      return {
        is: 'v-text-field',
        bind: {
          name,
          errorMessages: errors.collect(name),
          value: 'null',
          disabled: true,
        },
      };
    });

    return {
      rules,
      inputComponent,
    };
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
