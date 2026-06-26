<template>
  <c-layout column>
    <action-note-form
      :model-value="modelValue"
      :name="name"
      :depth="depth"
      @update:model-value="$emit('update:modelValue', $event)"
    />

    <c-form-block-row
      :label="t('common.state')"
      :depth="depth"
    >
      <c-select
        v-model="state"
        :items="states"
        :name="`${name}.state`"
        :label="t('common.state')"
        item-title="text"
        item-value="value"
        required
      />
    </c-form-block-row>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STATES } from '@/constants';

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

    const states = computed(() => Object.values(ALARM_STATES).map(value => ({
      value,
      text: t(`common.stateTypes.${value}`),
    })));

    const state = computed({
      get: () => props.modelValue.state ?? ALARM_STATES.critical,
      set: value => emit('update:modelValue', {
        ...props.modelValue,
        state: value,
      }),
    });

    return {
      t,
      states,
      state,
    };
  },
};
</script>
