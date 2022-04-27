<template lang="pug">
  div
    component(v-model="value", :is="component", :name="name")
</template>

<script>
import { ACTION_TYPES } from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import ActionChangeStateForm from './action-change-state-form.vue';
import ActionAssocticketForm from './action-assocticket-form.vue';
import ActionNoteForm from './action-note-form.vue';
import ActionPbehaviorForm from './action-pbehavior-form.vue';
import ActionSnoozeForm from './action-snooze-form.vue';
import ActionWebhookForm from './action-webhook-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ActionChangeStateForm,
    ActionAssocticketForm,
    ActionNoteForm,
    ActionPbehaviorForm,
    ActionSnoozeForm,
    ActionWebhookForm,
  },
  mixins: [formMixin, formValidationHeaderMixin],
  model: {
    prop: 'action',
    event: 'input',
  },
  props: {
    action: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'parameters',
    },
  },
  computed: {
    value: {
      get() {
        const { type, parameters } = this.action;

        return parameters[type];
      },
      set(value) {
        this.updateField(`parameters.${this.action.type}`, value);
      },
    },
    component() {
      return {
        [ACTION_TYPES.changeState]: 'action-change-state-form',
        [ACTION_TYPES.snooze]: 'action-snooze-form',
        [ACTION_TYPES.pbehavior]: 'action-pbehavior-form',
        [ACTION_TYPES.assocticket]: 'action-assocticket-form',
        [ACTION_TYPES.ack]: 'action-note-form',
        [ACTION_TYPES.ackremove]: 'action-note-form',
        [ACTION_TYPES.cancel]: 'action-note-form',
        [ACTION_TYPES.webhook]: 'action-webhook-form',
      }[this.action.type];
    },
  },
};
</script>
