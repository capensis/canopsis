<template lang="pug">
  div
    snmp-rule-form-field(label="state")
      v-flex(xs12)
        v-switch(
        :input-value="form.type",
        false-value="simple",
        true-value="template",
        label="To custom",
        @change="updateTypeField"
        )
    v-divider(light)
    template(v-if="form.type === 'template'")
      snmp-rule-form-vars-field(
      label="Define matching snmp var",
      :value="form.stateoid",
      :items="items",
      @input="updateField('stateoid', $event)"
      )
      v-layout(row, wrap)
        v-flex(xs12)
          v-layout(v-for="(state, key) in $constants.ENTITIES_STATES", :key="key", row, wrap, align-center)
            v-flex(xs2)
              v-chip(:style="{ backgroundColor: $constants.ENTITIES_STATES_STYLES[state].color }", label)
                strong.tt-uppercase {{ $t(`modals.createChangeStateEvent.states.${key}`) }}
            v-flex(xs10)
              v-text-field(
              :value="form[key]",
              placeholder="write template",
              @input="updateField(key, $event)"
              )
    template(v-else)
      v-layout.mt-3(row, wrap)
        v-flex(xs12)
          state-criticity-field(
          :value="form.state",
          @input="updateField('state', $event)"
          )
</template>

<script>
import formMixin from '@/mixins/form';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormField from './snmp-rule-form-field.vue';
import SnmpRuleFormVarsField from './snmp-rule-form-vars-field.vue';

export default {
  components: { StateCriticityField, SnmpRuleFormField, SnmpRuleFormVarsField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    updateTypeField(type) {
      const state = {
        type,
      };

      if (type === 'template') {
        state.stateoid = {};
      }

      this.updateModel(state);
    },
  },
};
</script>
