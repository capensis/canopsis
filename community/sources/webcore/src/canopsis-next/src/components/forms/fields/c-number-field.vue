<template>
  <v-text-field
    v-field.number="value"
    v-validate="rules"
    v-bind="$attrs"
    :label="label"
    :error-messages="errorMessages"
    :disabled="disabled"
    :hide-details="hideDetails"
    :name="name"
    :min="min"
    :max="max"
    :step="step"
    :loading="loading"
    type="number"
    @paste="$emit('paste', $event)"
    @click="$emit('click', $event)"
    @input.native="checkBadInput"
  >
    <template #append="">
      <slot name="append" />
    </template>
    <template #append-outer="">
      <slot name="append-outer" />
    </template>
  </v-text-field>
</template>

<script>
import { ref, computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Number, String],
      required: false,
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
    required: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: undefined,
    },
    max: {
      type: Number,
      default: undefined,
    },
    step: {
      type: [Number, String],
      default: undefined,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { errors } = useValidator();

    const textFieldRef = ref(null);
    const hasBadInput = ref(false);

    const rules = computed(() => ({
      required: props.required,
      decimal: !!props.step,
      min_value: props.min ?? false,
      max_value: props.max ?? false,
    }));

    const checkBadInput = event => hasBadInput.value = event.target.validity.badInput;

    const errorMessages = computed(() => {
      if (hasBadInput.value) {
        return [t('errors.invalidNumberFormat')];
      }

      return errors.collect(props.name);
    });

    return {
      textFieldRef,
      rules,
      checkBadInput,
      errorMessages,
    };
  },
};
</script>
