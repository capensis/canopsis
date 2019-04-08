<template lang="pug">
  div
    v-layout(row, wrap)
      v-flex(xs12)
        .body-2 oid
      v-flex.pr-1(xs6)
        v-autocomplete.pt-0(
        placeholder="Select a mib module",
        :items="items",
        :loading="loading",
        hide-no-data,
        return-object
        )
      v-flex.pl-1(xs6)
        v-select.pt-0(
        :items="items"
        )
    snmp-rule-form-field(
    :value="form.output",
    label="output",
    @input="updateField('output', $event)"
    )
    snmp-rule-form-field(
    :value="form.resource",
    label="resource",
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-field(
    :value="form.component",
    label="component",
    @input="updateField('component', $event)"
    )
    snmp-rule-form-field(
    :value="form.connector_name",
    label="connector_name",
    @input="updateField('connector_name', $event)"
    )
    v-layout(row, wrap)
      v-flex(xs12)
        .body-2 state
      v-flex(xs12)
        v-switch(label="To custom")
      v-flex(xs12)
        state-criticity-field(:value="form.state.state")
      v-flex(xs12)
        v-text-field(placeholder="Snmp vars match field")
      v-flex(xs12)
        v-layout(row, wrap)
          v-flex(xs2)
            v-chip(color="green", label)
              strong CRITICAL
          v-flex(xs10)
            v-text-field(placeholder="write template")
</template>

<script>
import formMixin from '@/mixins/form';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormField from './snmp-rule-form-field.vue';

export default {
  components: { StateCriticityField, SnmpRuleFormField },
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
  },
  data() {
    return {
      items: [],
      loading: false,
    };
  },
};
</script>
