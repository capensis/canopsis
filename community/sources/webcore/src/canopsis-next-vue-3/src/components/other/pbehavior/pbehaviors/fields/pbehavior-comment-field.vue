<template>
  <c-layout gap="2" align="center">
    <c-textarea
      :model-value="comment.message"
      :disabled="!!comment._id"
      :label="$t('common.message')"
      :name="messageFieldName"
      :max-length="max"
      required
      @update:model-value="updateField('message', $event)"
    />

    <c-action-btn type="delete" @click="$emit('remove')" />
  </c-layout>
</template>

<script setup>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';

const props = defineProps({
  comment: {
    type: Object,
    required: true,
  },
  max: {
    type: Number,
    default: 255,
  },
});

const emit = defineEmits(['update:comment', 'remove']);

const { updateField } = useModelField(props, emit);

const messageFieldName = computed(() => `${props.comment._id || props.comment.key}_message`);
</script>
