<template lang="pug">
  div
    snmp-rule-form-field-title(label="state")
    v-layout(row, wrap)
      v-flex(xs12)
        v-switch(
          :input-value="form.type",
          :false-value="$constants.SNMP_STATE_TYPES.simple",
          :true-value="$constants.SNMP_STATE_TYPES.template",
          :label="$t('modals.createSnmpRule.fields.state.labels.toCustom')",
          @change="updateTypeField"
        )
    v-divider(light)
    template(v-if="isTemplate")
      snmp-rule-form-module-mib-objects-form(
        :form="form.stateoid",
        :items="items",
        :label="$t('modals.createSnmpRule.fields.state.labels.defineVar')",
        @input="updateField('stateoid', $event)"
      )
      v-layout(row, wrap)
        v-flex(xs12)
          v-layout(v-for="(state, key) in $constants.ENTITIES_STATES", :key="key", row, wrap, align-center)
            v-flex(xs2)
              v-chip(:style="{ backgroundColor: $constants.ENTITIES_STATES_STYLES[state].color }", label)
                strong.state-title {{ $t(`modals.createChangeStateEvent.states.${key}`) }}
            v-flex(xs10)
              v-text-field(
                :value="form[key]",
                :placeholder="$t('modals.createSnmpRule.fields.state.labels.writeTemplate')",
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
import { SNMP_STATE_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormFieldTitle from './snmp-rule-form-field-title.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';

export default {
  components: {
    StateCriticityField,
    SnmpRuleFormFieldTitle,
    SnmpRuleFormModuleMibObjectsForm,
  },
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
  computed: {
    isTemplate() {
      return this.form.type === SNMP_STATE_TYPES.template;
    },
  },
  methods: {
    updateTypeField(type) {
      const state = {
        type,
      };

      if (type === SNMP_STATE_TYPES.template) {
        state.stateoid = {};
      }

      this.updateModel(state);
    },
  },
};
</script>

<style lang="scss" scoped>
  .state-title {
    text-transform: uppercase;
  }
</style>
