<template lang="pug">
  div
    snmp-rule-form-module-form(
    :form="form.oid",
    :moduleMibs.sync="moduleMibs",
    @input="updateField('oid', $event)"
    )
    v-layout(v-if="selectedModuleMib", row, wrap)
      v-alert(:value="selectedModuleMib", color="grey darken-1") {{ selectedModuleMib.description }}
    snmp-rule-form-vars-field(
    :value="form.output",
    :items="selectedModuleMibObjects",
    label="output",
    large,
    @input="updateField('output', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.resource",
    :items="selectedModuleMibObjects",
    label="resource",
    large,
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.component",
    :items="selectedModuleMibObjects",
    label="component",
    large,
    @input="updateField('component', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.connector_name",
    :items="selectedModuleMibObjects",
    label="connector_name",
    large,
    @input="updateField('connector_name', $event)"
    )
    snmp-rule-form-state-form(
    :form="form.state",
    :items="selectedModuleMibObjects",
    @input="updateField('state', $event)"
    )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';
import SnmpRuleFormField from './snmp-rule-form-field.vue';
import SnmpRuleFormVarsField from './snmp-rule-form-vars-field.vue';

export default {
  components: {
    SnmpRuleFormModuleForm, SnmpRuleFormStateForm, SnmpRuleFormField, SnmpRuleFormVarsField,
  },
  mixins: [formMixin, entitiesSnmpMibMixin],
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
      moduleMibs: [],
    };
  },
  computed: {
    selectedModuleMib() {
      return this.moduleMibs.find(({ name }) => name === this.form.oid.mibName);
    },

    selectedModuleMibObjects() {
      if (this.selectedModuleMib) {
        return Object.keys(this.selectedModuleMib.objects);
      }

      return [];
    },
  },
};
</script>

<style lang="scss" scoped>
  .tt-uppercase {
    text-transform: uppercase;
  }
</style>
