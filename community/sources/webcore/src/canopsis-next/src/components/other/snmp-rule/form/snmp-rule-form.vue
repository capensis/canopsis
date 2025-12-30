<template>
  <v-layout class="gap-2" column>
    <c-enabled-field v-field="form.enabled" hide-details />
    <snmp-rule-form-module-form v-field="form.oid" />
    <v-layout>
      <v-flex xs12>
        <v-alert
          :value="!!selectedModuleMib?.description"
          class="mt-3"
          color="grey darken-1"
        >
          {{ selectedModuleMib?.description }}
        </v-alert>
      </v-flex>
    </v-layout>
    <snmp-rule-form-module-mib-objects-form
      v-field="form.output"
      :items="selectedModuleMibObjects"
      :label="$t('snmpRule.output')"
      large
    />
    <snmp-rule-form-module-mib-objects-form
      v-field="form.component"
      :items="selectedModuleMibObjects"
      :label="$t('snmpRule.component')"
      name="component"
      required
      large
    />
    <snmp-rule-form-module-mib-objects-form
      v-field="form.resource"
      :items="selectedModuleMibObjects"
      :label="$t('snmpRule.resource')"
      large
    />
    <snmp-rule-form-module-mib-objects-form
      v-field="form.connector_name"
      :items="selectedModuleMibObjects"
      :label="$t('snmpRule.connectorName')"
      name="connector_name"
      required
      large
    />
    <snmp-rule-form-state-form
      v-field="form.state"
      :items="selectedModuleMibObjects"
    />
    <snmp-rule-form-tags-field
      v-field="form.tags"
      :items="selectedModuleMibObjects"
    />
    <snmp-rule-form-extra-field
      v-field="form.extra"
      :items="selectedModuleMibObjects"
    />
  </v-layout>
</template>

<script>
import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';
import SnmpRuleFormTagsField from './snmp-rule-form-tags-field.vue';
import SnmpRuleFormExtraField from './snmp-rule-form-extra-field.vue';

export default {
  components: {
    SnmpRuleFormModuleForm,
    SnmpRuleFormModuleMibObjectsForm,
    SnmpRuleFormStateForm,
    SnmpRuleFormTagsField,
    SnmpRuleFormExtraField,
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
  computed: {
    selectedModuleMib() {
      return this.form.oid?.mib;
    },

    selectedModuleMibObjects() {
      return this.selectedModuleMib?.objects
        ? Object.keys(this.selectedModuleMib.objects)
        : [];
    },
  },
  methods: {
    addTag() {},
  },
};
</script>
