<template>
  <v-text-field
    v-field="value"
    v-bind="$attrs"
    :error-messages="errorMessages"
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

import { useInjectValidator } from '@/hooks/validator/inject-validator';
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
  },
  setup(props) {
    const { t } = useI18n();
    const validator = useInjectValidator();
    const instance = useComponentInstance();

    const errorMessages = computed(() => (validator.errors.has(props.name) ? t('validation.messages.regex') : undefined));

    const attachRangeRule = () => {
      validator.attach({
        name: props.name,
        rules: 'required:true',
        getter: () => props.value.split(',').every((value) => {
          const preparedValue = Number(value);

          if (isNaN(preparedValue) || !isNumber(preparedValue)) {
            return false;
          }

          return isNumber(props.min) && isNumber(props.max)
            ? preparedValue >= props.min && preparedValue <= props.max
            : true;
        }),
        vm: instance,
      });
    };
    const validateRangeRule = () => validator.validate(props.name);
    const detachRangeRule = () => validator.detach(props.name);

    attachRangeRule();

    watch(() => props.value, () => {
      if (errorMessages.value) {
        validateRangeRule();
      }
    });
    onBeforeUnmount(detachRangeRule);

    return {
      errorMessages,
    };
  },
};
</script>
