<template lang="pug">
  div
    snmp-rule-form-module-form(
      v-field="form.oid",
      :module-mibs.sync="moduleMibs"
    )
    v-layout(v-if="selectedModuleMib", row)
      v-flex(xs12)
        v-alert.mt-3(
          :value="selectedModuleMib.description",
          color="grey darken-1"
        ) {{ selectedModuleMib.description }}
    snmp-rule-form-module-mib-objects-form(
      v-field="form.output",
      :items="selectedModuleMibObjects",
      :label="$t('snmpRule.output')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.component",
      :items="selectedModuleMibObjects",
      :label="$t('snmpRule.component')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.resource",
      :items="selectedModuleMibObjects",
      :label="$t('snmpRule.resource')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.connector_name",
      :items="selectedModuleMibObjects",
      :label="$t('snmpRule.connectorName')",
      large
    )
    snmp-rule-form-state-form(
      v-field="form.state",
      :items="selectedModuleMibObjects"
    )
</template>

<script>
import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';

export default {
  components: {
    SnmpRuleFormModuleForm,
    SnmpRuleFormModuleMibObjectsForm,
    SnmpRuleFormStateForm,
  },
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
