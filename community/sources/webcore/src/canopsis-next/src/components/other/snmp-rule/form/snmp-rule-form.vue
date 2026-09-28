<template>
  <v-layout class="gap-3" column>
    <c-enabled-field
      v-field="form.enabled"
      hide-details
      with-background
    />

    <snmp-rule-form-module-form v-field="form.oid" class="mb-2" />

    <v-alert
      :value="!!selectedModuleMib?.description"
      class="mt-3"
      color="grey darken-1"
    >
      {{ selectedModuleMib?.description }}
    </v-alert>

    <c-form-block>
      <c-form-block-row :label="$t('snmpRule.output')">
        <snmp-rule-form-module-mib-objects-form
          v-field="form.output"
          :items="selectedModuleMibObjects"
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('snmpRule.component')">
        <snmp-rule-form-module-mib-objects-form
          v-field="form.component"
          :items="selectedModuleMibObjects"
          required
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('snmpRule.resource')">
        <snmp-rule-form-module-mib-objects-form
          v-field="form.resource"
          :items="selectedModuleMibObjects"
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('snmpRule.connectorName')">
        <snmp-rule-form-module-mib-objects-form
          v-field="form.connector_name"
          :items="selectedModuleMibObjects"
          required
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.state')" indented>
        <snmp-rule-form-state-form
          v-field="form.state"
          :items="selectedModuleMibObjects"
        />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.tag', 2)" indented>
        <snmp-rule-form-tags-field
          v-field="form.tags"
          :items="selectedModuleMibObjects"
        />
      </c-form-block-row>

      <c-form-block-row :label="$tc('common.customField', 2)" indented>
        <snmp-rule-form-extra-field
          v-field="form.extra"
          :items="selectedModuleMibObjects"
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';
import SnmpRuleFormTagsField from './fields/snmp-rule-form-tags-field.vue';
import SnmpRuleFormExtraField from './fields/snmp-rule-form-extra-field.vue';

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
  setup(props) {
    const selectedModuleMib = computed(() => props.form.oid?.mib);

    const selectedModuleMibObjects = computed(() => {
      const mib = selectedModuleMib.value;

      return mib?.objects ? Object.keys(mib.objects) : [];
    });

    return {
      selectedModuleMib,
      selectedModuleMibObjects,
    };
  },
};
</script>
