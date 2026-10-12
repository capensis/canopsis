<template>
  <component
    v-if="bindProps.is"
    v-bind="bindProps"
    :is="bindProps.is"
    v-field="value"
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

const ACTION_COMPONENT_BY_TYPE = {
  [ACTION_TYPES.changeState]: 'c-change-state-field',
  [ACTION_TYPES.snooze]: 'action-snooze-form',
  [ACTION_TYPES.unsnooze]: 'action-note-form',
  [ACTION_TYPES.pbehavior]: 'action-pbehavior-form',
  [ACTION_TYPES.assocticket]: 'action-assocticket-form',
  [ACTION_TYPES.ack]: 'action-note-form',
  [ACTION_TYPES.ackremove]: 'action-note-form',
  [ACTION_TYPES.cancel]: 'action-note-form',
  [ACTION_TYPES.webhook]: 'action-webhook-form',
};

export default {
  inject: ['$validator'],
  components: {
    ActionAssocticketForm,
    ActionNoteForm,
    ActionPbehaviorForm,
    ActionSnoozeForm,
    ActionWebhookForm,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
    hasPreviousWebhook: {
      type: Boolean,
      default: false,
    },
    templateVars: {
      type: Object,
      default: () => ({}),
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
  setup(props) {
    const bindProps = computed(() => {
      const childProps = {
        is: ACTION_COMPONENT_BY_TYPE[props.type],
        name: props.name,
        depth: props.depth,
      };

      if (props.type === ACTION_TYPES.webhook) {
        childProps.hasPrevious = props.hasPreviousWebhook;
      }

      if (props.type === ACTION_TYPES.changeState) {
        childProps.variables = props.templateVars?.output;
      } else {
        childProps.templateVars = props.templateVars;
      }

      return childProps;
    });

    return {
      bindProps,
    };
  },
};
</script>
