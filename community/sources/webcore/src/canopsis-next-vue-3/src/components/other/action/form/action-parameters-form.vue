<template>
  <component
    :is="componentName"
    v-bind="componentBind"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import { ACTION_TYPES } from '@/constants';

import ActionAssocticketForm from './action-assocticket-form.vue';
import ActionNoteForm from './action-note-form.vue';
import ActionPbehaviorForm from './action-pbehavior-form.vue';
import ActionSnoozeForm from './action-snooze-form.vue';
import ActionWebhookForm from './action-webhook-form.vue';

const ACTION_COMPONENTS = {
  [ACTION_TYPES.snooze]: 'action-snooze-form',
  [ACTION_TYPES.unsnooze]: 'action-note-form',
  [ACTION_TYPES.pbehavior]: 'action-pbehavior-form',
  [ACTION_TYPES.pbehaviorRemove]: 'action-note-form',
  [ACTION_TYPES.ack]: 'action-note-form',
  [ACTION_TYPES.ackremove]: 'action-note-form',
  [ACTION_TYPES.cancel]: 'action-note-form',
  [ACTION_TYPES.assocticket]: 'action-assocticket-form',
  [ACTION_TYPES.changeState]: 'c-change-state-field',
  [ACTION_TYPES.webhook]: 'action-webhook-form',
};

export default {
  components: {
    ActionAssocticketForm,
    ActionNoteForm,
    ActionPbehaviorForm,
    ActionSnoozeForm,
    ActionWebhookForm,
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
    hasPreviousWebhook: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:modelValue'],

  setup(props) {
    const componentName = computed(() => ACTION_COMPONENTS[props.type] || 'action-note-form');
    const componentBind = computed(() => {
      const bind = {
        modelValue: props.modelValue,
        name: props.name,
        depth: props.depth,
      };

      if (props.type === ACTION_TYPES.webhook) {
        bind.hasPrevious = props.hasPreviousWebhook;
      }

      if (props.type === ACTION_TYPES.changeState) {
        bind.variables = props.templateVars?.output;
      } else {
        bind.templateVars = props.templateVars;
      }

      return bind;
    });

    return {
      componentName,
      componentBind,
    };
  },
};
</script>
