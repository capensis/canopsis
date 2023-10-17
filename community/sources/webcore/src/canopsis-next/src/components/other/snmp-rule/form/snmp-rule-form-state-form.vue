<template>
  <div>
    <snmp-rule-form-field-title label="state" />
    <v-layout wrap>
      <v-flex xs12>
        <v-switch
          :input-value="form.type"
          :false-value="$constants.SNMP_STATE_TYPES.simple"
          :true-value="$constants.SNMP_STATE_TYPES.template"
          :label="$t('snmpRule.toCustom')"
          color="primary"
          @change="updateTypeField"
        />
      </v-flex>
    </v-layout>
    <v-divider light />
    <template v-if="isTemplate">
      <snmp-rule-form-module-mib-objects-form
        v-field="form.stateoid"
        :items="items"
        :label="$t('snmpRule.defineVar')"
      />
      <v-layout wrap>
        <v-flex xs12>
          <v-layout
            v-for="{ value, color, key, text } in availableStates"
            :key="value"
            wrap
            align-center
          >
            <v-flex xs2>
              <v-chip
                :style="{ backgroundColor: color }"
                label
              >
                <strong class="state-title">{{ text }}</strong>
              </v-chip>
            </v-flex>
            <v-flex xs10>
              <v-text-field
                v-field="form[key]"
                :placeholder="$t('snmpRule.writeTemplate')"
              />
            </v-flex>
          </v-layout>
        </v-flex>
      </v-layout>
    </template>
    <template v-else>
      <v-layout
        class="mt-3"
        wrap
      >
        <v-flex xs12>
          <state-criticity-field v-field="form.state" />
        </v-flex>
      </v-layout>
    </template>
  </div>
</template>

<script>
import { SNMP_STATE_TYPES, SNMP_TEMPLATE_STATE_STATES } from '@/constants';

import { getSnmpRuleStateColor } from '@/helpers/entities/snmp-rule/color';

import { formBaseMixin } from '@/mixins/form';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormFieldTitle from './snmp-rule-form-field-title.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';

export default {
  components: {
    StateCriticityField,
    SnmpRuleFormFieldTitle,
    SnmpRuleFormModuleMibObjectsForm,
  },
  mixins: [formBaseMixin],
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
    stateValues: {
      type: Object,
      default: () => SNMP_TEMPLATE_STATE_STATES,
    },
  },
  computed: {
    isTemplate() {
      return this.form.type === SNMP_STATE_TYPES.template;
    },

    availableStates() {
      return Object.entries(this.stateValues).map(([key, state]) => ({
        key,
        text: this.$t(`snmpRule.states.${key}`),
        value: state,
        color: getSnmpRuleStateColor(state),
      }));
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
