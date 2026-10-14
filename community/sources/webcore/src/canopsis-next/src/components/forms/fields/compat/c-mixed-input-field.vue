<template>
  <component
    :is="inputComponent.is"
    class="c-mixed-input-field"
    v-bind="inputComponent.bind"
  />
</template>

<script>
import { isNull, pick } from 'lodash';
import { computed } from 'vue';

import { PATTERN_FIELD_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';

export default {
  emits: ['update:value'],

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
    required: {
      type: Boolean,
      default: false,
    },
  },

  setup(props, { emit }) {
    const { updateModel } = useModelField(props, emit);

    const switchLabel = computed(() => String(props.value));
    const isInputTypeText = computed(
      () => [PATTERN_FIELD_TYPES.number, PATTERN_FIELD_TYPES.string].includes(props.inputType),
    );

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
          ? {
              ...pick(props, ['items']),
              itemTitle: props.itemText,
              itemValue: props.itemValue,
              returnObject: false,
            }
          : {};

        return {
          is: props.items.length ? 'c-combobox' : 'c-text-field',
          bind: {
            ...pick(props, ['disabled', 'hideDetails']),
            ...additionalProps,
            modelValue: props.value,
            name,
            placeholder: props.label,
            class: 'c-mixed-input-field__text',
            type: props.inputType === PATTERN_FIELD_TYPES.number ? 'number' : 'text',
            required: props.required,
            'onUpdate:modelValue': updateTextFieldValue,
          },
        };
      }

      if (props.inputType === PATTERN_FIELD_TYPES.boolean) {
        const name = `${props.name}.switch`;

        return {
          is: 'c-enabled-field',
          bind: {
            class: 'ma-0 c-mixed-input-field__switch',
            name,
            modelValue: props.value,
            label: switchLabel.value,
            disabled: props.disabled,
            color: 'primary',
            hideDetails: true,
            'onUpdate:modelValue': updateModel,
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
            required: props.required,
            'onUpdate:values': updateModel,
          },
        };
      }

      const name = `${props.name}.text_field`;

      return {
        is: 'c-text-field',
        bind: {
          name,
          modelValue: 'null',
          disabled: true,
        },
      };
    });

    return {
      inputComponent,
    };
  },
};
</script>

<style lang="scss">
.c-mixed-input-field {
  &__switch {
    & .v-label {
      text-transform: capitalize;
    }
  }
}
</style>
