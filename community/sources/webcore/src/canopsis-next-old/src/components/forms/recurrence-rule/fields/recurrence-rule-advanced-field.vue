<template>
  <v-text-field
    v-field="value"
    v-bind="$attrs"
    :error-messages="errorMessage"
    :name="name"
    persistent-hint
  >
    <template #append="">
      <c-help-icon
        v-if="helpText"
        :text="helpText"
        icon="help"
        max-width="250"
        left
      />
    </template>
  </v-text-field>
</template>

<script>
import { isNumber, isNaN } from 'lodash';
import { computed, onBeforeUnmount, watch } from 'vue';

import { useValidator } from '@/hooks/validator/validator';
import { useComponentInstance } from '@/hooks/vue';
import { useI18n } from '@/hooks/i18n';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'regex',
    },
    min: {
      type: Number,
      required: false,
    },
    max: {
      type: Number,
      required: false,
    },
    negative: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const validator = useValidator();
    const instance = useComponentInstance();

    const invalidMessage = computed(
      () => {
        const isMinZero = props.min === 0;

        if (!props.negative || isMinZero) {
          return t(
            'recurrenceRule.errors.invalidRange',
            {
              min: props.negative && isMinZero ? -props.max : props.min,
              max: props.max,
            },
          );
        }

        return t(
          'recurrenceRule.errors.invalidRangeNegative',
          { min: props.min, max: props.max },
        );
      },
    );
    const errorMessage = computed(() => (validator.errors.has(props.name) ? invalidMessage.value : undefined));

    const isValueValid = () => !props.value.length || props.value.split(',').every((value) => {
      const preparedValue = props.negative
        ? Math.abs(+value)
        : +value;

      if (isNaN(preparedValue) || !isNumber(preparedValue)) {
        return false;
      }

      return isNumber(props.min) && isNumber(props.max)
        ? preparedValue >= props.min && preparedValue <= props.max
        : true;
    });
    const attachRangeRule = () => {
      validator.attach({
        name: props.name,
        rules: 'required:true',
        getter: isValueValid,
        vm: instance,
      });
    };
    const validateRangeRule = () => validator.validate(props.name);
    const detachRangeRule = () => validator.detach(props.name);

    attachRangeRule();

    watch(() => props.value, () => {
      if (errorMessage.value) {
        validateRangeRule();
      }
    });
    onBeforeUnmount(detachRangeRule);

    return {
      errorMessage,
    };
  },
};
</script>
