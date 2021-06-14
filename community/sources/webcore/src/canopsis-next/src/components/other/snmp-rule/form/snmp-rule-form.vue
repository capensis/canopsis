<template lang="pug">
  div
    snmp-rule-form-module-form(
      v-field="form.oid",
      :moduleMibs.sync="moduleMibs"
    )
    v-layout(v-if="selectedModuleMib", row, wrap)
      v-alert.mt-3(
        :value="selectedModuleMib.description",
        color="grey darken-1"
      ) {{ selectedModuleMib.description }}
    snmp-rule-form-module-mib-objects-form(
      v-field="form.output",
      :items="selectedModuleMibObjects",
      :label="$t('modals.createSnmpRule.fields.output.title')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.component",
      :items="selectedModuleMibObjects",
      :label="$t('modals.createSnmpRule.fields.component.title')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.resource",
      :items="selectedModuleMibObjects",
      :label="$t('modals.createSnmpRule.fields.resource.title')",
      large
    )
    snmp-rule-form-module-mib-objects-form(
      v-field="form.connector_name",
      :items="selectedModuleMibObjects",
      :label="$t('modals.createSnmpRule.fields.connectorName.title')",
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
