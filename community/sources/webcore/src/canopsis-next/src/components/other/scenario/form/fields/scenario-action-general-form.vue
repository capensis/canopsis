<template lang="pug">
  div.mt-4
    component(v-model="value", :is="component", :name="name")
</template>

<script>
import { SCENARIO_ACTION_TYPES } from '@/constants';

import { formMixin, formValidationHeaderMixin } from '@/mixins/form';

import ScenarioActionAssocticketField from './scenario-action-assocticket-field.vue';
import ScenarioActionNoteField from './scenario-action-note-field.vue';
import ScenarioActionPbehaviorField from './scenario-action-pbehavior-field.vue';
import ScenarioActionSnoozeField from './scenario-action-snooze-field.vue';
import ScenarioActionWebhookField from './scenario-action-webhook-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ScenarioActionAssocticketField,
    ScenarioActionNoteField,
    ScenarioActionPbehaviorField,
    ScenarioActionSnoozeField,
    ScenarioActionWebhookField,
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
        [SCENARIO_ACTION_TYPES.changeState]: 'c-change-state-field',
        [SCENARIO_ACTION_TYPES.snooze]: 'scenario-action-snooze-field',
        [SCENARIO_ACTION_TYPES.pbehavior]: 'scenario-action-pbehavior-field',
        [SCENARIO_ACTION_TYPES.assocticket]: 'scenario-action-assocticket-field',
        [SCENARIO_ACTION_TYPES.ack]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.ackremove]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.declareticket]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.cancel]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.webhook]: 'scenario-action-webhook-field',
      }[this.action.type];
    },
  },
};
</script>
