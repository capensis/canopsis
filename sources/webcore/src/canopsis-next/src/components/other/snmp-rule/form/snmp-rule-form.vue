<template lang="pug">
  div
    snmp-rule-form-module-form(
    :form="form.oid",
    :moduleMibs.sync="moduleMibs",
    @input="updateField('oid', $event)"
    )
    v-layout(v-if="selectedModuleMib", row, wrap)
      v-alert.mt-3(
      :value="selectedModuleMib.description",
      color="grey darken-1"
      ) {{ selectedModuleMib.description }}
    snmp-rule-form-module-mib-objects-form(
    :form="form.output",
    :items="selectedModuleMibObjects",
    :label="$t('modals.createSnmpRule.fields.output.title')"
    large,
    @input="updateField('output', $event)"
    )
    snmp-rule-form-module-mib-objects-form(
    :form="form.component",
    :items="selectedModuleMibObjects",
    :label="$t('modals.createSnmpRule.fields.component.title')"
    large,
    @input="updateField('component', $event)"
    )
    snmp-rule-form-module-mib-objects-form(
    :form="form.resource",
    :items="selectedModuleMibObjects",
    :label="$t('modals.createSnmpRule.fields.resource.title')"
    large,
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-module-mib-objects-form(
    :form="form.connector_name",
    :items="selectedModuleMibObjects",
    :label="$t('modals.createSnmpRule.fields.connectorName.title')"
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

import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';

export default {
  components: {
    SnmpRuleFormModuleForm,
    SnmpRuleFormModuleMibObjectsForm,
    SnmpRuleFormStateForm,
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
