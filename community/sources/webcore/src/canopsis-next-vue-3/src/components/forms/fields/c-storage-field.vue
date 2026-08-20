<template>
  <c-layout>
    <template v-if="modelValue">
      <c-text-field
        :model-value="modelValue"
        :disabled="disabled"
        :required="isRequired"
        :name="name"
        class="mt-0 pt-0"
        readonly
        hide-details
      />
      <c-action-btn
        :disabled="disabled"
        class="ml-2"
        type="edit"
        @click="$emit('edit', modelValue)"
      />
      <c-action-btn
        :disabled="disabled"
        type="delete"
        @click="$emit('remove')"
      />
    </template>
    <template v-else>
      <c-text-field
        :model-value="modelValue"
        :disabled="disabled"
        :required="isRequired"
        :name="name"
        class="d-none"
        hide-details
      />
      <c-btn-with-error
        :error="errorMessage"
        :disabled="disabled"
        class="ml-0"
        @click="$emit('add')"
      >
        {{ t('common.add') }}
      </c-btn-with-error>
    </template>
  </c-layout>
</template>

<script>
import { computed, toRef } from 'vue';

import { useFieldError } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'storage',
    },
  },

  emits: ['update:modelValue', 'edit', 'remove', 'add'],

  setup(props) {
    const { t } = useI18n();
    const fieldError = useFieldError(toRef(props, 'name'));
    const errorMessage = computed(() => fieldError.value || '');
    const isRequired = computed(() => props.required && !props.disabled);

    return {
      t,
      errorMessage,
      isRequired,
    };
  },
};
</script>
