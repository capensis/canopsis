<template>
  <slot :has-error="hasError(name)" :messages="messages" />
</template>

<script setup>
import { computed, toRef, watch } from 'vue';

import { useFormErrorsMap } from '@/hooks/form';
import { useValidationAttachRequiredForField } from '@/hooks/validator/validation-attach-required';

const props = defineProps({
  name: {
    type: String,
    required: true,
  },
  valid: {
    type: Boolean,
    default: false,
  },
  message: {
    type: String,
    required: true,
  },
});

const { hasError } = useFormErrorsMap();
const { asyncValidateRequiredRule } = useValidationAttachRequiredForField(
  toRef(props, 'name'),
  () => props.valid,
  false,
);

const messages = computed(() => (hasError(props.name) ? [props.message] : []));

watch(
  () => props.valid,
  () => {
    if (hasError(props.name)) {
      asyncValidateRequiredRule();
    }
  },
);
</script>
