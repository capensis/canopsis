<template>
  <c-layout>
    <template v-if="modelValue">
      <c-text-field
        :model-value="modelValue"
        :disabled="disabled"
        :required="!disabled"
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
        :required="!disabled"
        :name="name"
        class="d-none"
      />
      <v-btn
        :disabled="disabled"
        color="primary"
        class="ml-0"
        @click="$emit('add')"
      >
        {{ $t('common.add') }}
      </v-btn>
    </template>
  </c-layout>
</template>

<script>
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
    name: {
      type: String,
      default: 'storage',
    },
  },

  emits: ['update:modelValue', 'edit', 'remove', 'add'],
};
</script>
