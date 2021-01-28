<template lang="pug">
  component(v-model="value", :is="component")
</template>

<script>
import { SCENARIO_ACTION_TYPES } from '@/constants';

import formMixin from '@/mixins/form/object';

import ScenarioActionAssocticketField from './scenario-action-assocticket-field.vue';
import ScenarioActionNoteField from './scenario-action-note-field.vue';
import ScenarioActionChangeStateField from './scenario-action-change-state-field.vue';
import ScenarioActionPbehaviorField from './scenario-action-pbehavior-field.vue';
import ScenarioActionSnoozeField from './scenario-action-snooze-field.vue';

export default {
  inject: ['$validator'],
  components: {
    ScenarioActionAssocticketField,
    ScenarioActionNoteField,
    ScenarioActionChangeStateField,
    ScenarioActionPbehaviorField,
    ScenarioActionSnoozeField,
  },
  mixins: [formMixin],
  model: {
    prop: 'action',
    event: 'input',
  },
  props: {
    action: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
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
        [SCENARIO_ACTION_TYPES.snooze]: 'scenario-action-snooze-field',
        [SCENARIO_ACTION_TYPES.pbehavior]: 'scenario-action-pbehavior-field',
        [SCENARIO_ACTION_TYPES.changeState]: 'scenario-action-change-state-field',
        [SCENARIO_ACTION_TYPES.assocticket]: 'scenario-action-assocticket-field',
        [SCENARIO_ACTION_TYPES.ack]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.ackremove]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.declareticket]: 'scenario-action-note-field',
        [SCENARIO_ACTION_TYPES.cancel]: 'scenario-action-note-field',
      }[this.action.type];
    },
  },
};
</script>
