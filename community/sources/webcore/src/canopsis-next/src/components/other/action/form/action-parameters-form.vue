<template>
  <div>
    <component
      v-field="value"
      v-bind="props"
      :is="props.is"
    />
  </div>
</template>

<script>
import { ACTION_TYPES } from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import ActionAssocticketForm from './action-assocticket-form.vue';
import ActionNoteForm from './action-note-form.vue';
import ActionPbehaviorForm from './action-pbehavior-form.vue';
import ActionSnoozeForm from './action-snooze-form.vue';
import ActionWebhookForm from './action-webhook-form.vue';

export default {
  inject: ['$validator'],
  components: {
    ActionAssocticketForm,
    ActionNoteForm,
    ActionPbehaviorForm,
    ActionSnoozeForm,
    ActionWebhookForm,
  },
  mixins: [formMixin, formValidationHeaderMixin],
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
  },
  computed: {
    component() {
      return {
        [ACTION_TYPES.changeState]: 'c-change-state-field',
        [ACTION_TYPES.snooze]: 'action-snooze-form',
        [ACTION_TYPES.pbehavior]: 'action-pbehavior-form',
        [ACTION_TYPES.assocticket]: 'action-assocticket-form',
        [ACTION_TYPES.ack]: 'action-note-form',
        [ACTION_TYPES.ackremove]: 'action-note-form',
        [ACTION_TYPES.cancel]: 'action-note-form',
        [ACTION_TYPES.webhook]: 'action-webhook-form',
      }[this.type];
    },

    props() {
      const props = {
        is: this.component,
        name: this.name,
      };

      if (this.type === ACTION_TYPES.webhook) {
        props.hasPrevious = this.hasPreviousWebhook;
      }

      return props;
    },
  },
};
</script>
