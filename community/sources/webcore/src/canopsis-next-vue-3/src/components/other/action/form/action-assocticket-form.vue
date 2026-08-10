<template>
  <c-layout column>
    <action-note-form
      :model-value="modelValue"
      :name="name"
      :depth="depth"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <c-form-block-row :label="t('common.ticket')" :depth="depth">
      <c-text-field v-model="ticket" :name="`${name}.ticket`" :label="t('common.ticket')" />
    </c-form-block-row>

    <c-form-block-row :label="t('common.url')" :depth="depth">
      <c-text-field v-model="ticketUrl" :name="`${name}.ticket_url`" :label="t('common.url')" />
    </c-form-block-row>

    <c-form-block-row :label="t('common.title')" :depth="depth">
      <c-text-field
        v-model="ticketUrlTitle"
        :name="`${name}.ticket_url_title`"
        :label="t('common.title')"
      />
    </c-form-block-row>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import ActionNoteForm from './action-note-form.vue';

export default {
  components: { ActionNoteForm },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const model = field =>
      computed({
        get: () => props.modelValue[field] ?? '',
        set: value =>
          emit('update:modelValue', {
            ...props.modelValue,
            [field]: value,
          }),
      });

    return {
      t,
      ticket: model('ticket'),
      ticketUrl: model('ticket_url'),
      ticketUrlTitle: model('ticket_url_title'),
    };
  },
};
</script>
