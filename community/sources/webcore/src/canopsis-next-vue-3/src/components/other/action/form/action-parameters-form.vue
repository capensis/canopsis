<template>
  <component
    :is="componentName"
    :model-value="modelValue"
    :name="name"
    :depth="depth"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { ACTION_TYPES } from '@/constants';

import ActionAssocticketForm from './action-assocticket-form.vue';
import ActionChangeStateForm from './action-change-state-form.vue';
import ActionNoteForm from './action-note-form.vue';
import ActionSnoozeForm from './action-snooze-form.vue';

const ACTION_COMPONENTS = {
  [ACTION_TYPES.snooze]: 'action-snooze-form',
  [ACTION_TYPES.ack]: 'action-note-form',
  [ACTION_TYPES.ackremove]: 'action-note-form',
  [ACTION_TYPES.cancel]: 'action-note-form',
  [ACTION_TYPES.assocticket]: 'action-assocticket-form',
  [ACTION_TYPES.changeState]: 'action-change-state-form',
  [ACTION_TYPES.pbehavior]: 'action-note-form',
};

export default {
  components: {
    ActionAssocticketForm,
    ActionChangeStateForm,
    ActionNoteForm,
    ActionSnoozeForm,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: String,
      default: ACTION_TYPES.ack,
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

  setup(props) {
    const componentName = computed(() => ACTION_COMPONENTS[props.type] || 'action-note-form');

    return {
      componentName,
    };
  },
};
</script>
