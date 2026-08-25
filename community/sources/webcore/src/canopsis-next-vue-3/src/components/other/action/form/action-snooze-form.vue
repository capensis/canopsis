<template>
  <c-layout column>
    <action-note-form
      :model-value="modelValue"
      :name="name"
      :depth="depth"
      :template-vars="templateVars"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <c-form-block-row :label="t('common.duration')" :depth="depth">
      <c-duration-field
        v-model:duration="duration"
        :name="`${name}.duration`"
        :label="t('common.duration')"
        required
      />
    </c-form-block-row>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { TIME_UNITS } from '@/constants';

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
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const duration = computed({
      get: () => props.modelValue.duration ?? { value: 1, unit: TIME_UNITS.hour },
      set: value =>
        emit('update:modelValue', {
          ...props.modelValue,
          duration: value,
        }),
    });

    return {
      t,
      duration,
    };
  },
};
</script>
