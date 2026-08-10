<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabled" hide-details with-background />

    <snmp-rule-form-module-form v-model:form="oid" class="mb-2" />

    <v-alert v-if="selectedModuleMib?.description" class="mt-3" color="grey-darken-1">
      {{ selectedModuleMib?.description }}
    </v-alert>

    <c-form-block>
      <c-form-block-row :label="t('snmpRule.output')">
        <snmp-rule-form-module-mib-objects-form
          v-model:form="output"
          :items="selectedModuleMibObjects"
          name="output"
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="t('snmpRule.component')">
        <snmp-rule-form-module-mib-objects-form
          v-model:form="component"
          :items="selectedModuleMibObjects"
          name="component"
          required
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="t('snmpRule.resource')">
        <snmp-rule-form-module-mib-objects-form
          v-model:form="resource"
          :items="selectedModuleMibObjects"
          name="resource"
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="t('snmpRule.connectorName')">
        <snmp-rule-form-module-mib-objects-form
          v-model:form="connectorName"
          :items="selectedModuleMibObjects"
          name="connector_name"
          required
          large
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.state')" indented>
        <snmp-rule-form-state-form v-model:form="state" :items="selectedModuleMibObjects" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.tag', 2)" indented>
        <snmp-rule-form-tags-field v-model:form="tags" :items="selectedModuleMibObjects" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.customField', 2)" indented>
        <snmp-rule-form-extra-field v-model:form="extra" :items="selectedModuleMibObjects" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import SnmpRuleFormExtraField from './fields/snmp-rule-form-extra-field.vue';
import SnmpRuleFormTagsField from './fields/snmp-rule-form-tags-field.vue';
import SnmpRuleFormModuleForm from './snmp-rule-form-module-form.vue';
import SnmpRuleFormModuleMibObjectsForm from './snmp-rule-form-module-mib-objects-form.vue';
import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';

export default {
  components: {
    SnmpRuleFormModuleForm,
    SnmpRuleFormModuleMibObjectsForm,
    SnmpRuleFormStateForm,
    SnmpRuleFormTagsField,
    SnmpRuleFormExtraField,
  },

  setup() {
    const { t } = useI18n();
    const { enabled, oid, output, component, resource, connectorName, state, tags, extra } =
      useFormModels({
        enabled: 'enabled',
        oid: 'oid',
        output: 'output',
        component: 'component',
        resource: 'resource',
        connectorName: 'connector_name',
        state: 'state',
        tags: 'tags',
        extra: 'extra',
      });

    const selectedModuleMib = computed(() => oid.value?.mib);

    const selectedModuleMibObjects = computed(() => {
      const mib = selectedModuleMib.value;

      return mib?.objects ? Object.keys(mib.objects) : [];
    });

    return {
      t,
      enabled,
      oid,
      output,
      component,
      resource,
      connectorName,
      state,
      tags,
      extra,
      selectedModuleMib,
      selectedModuleMibObjects,
    };
  },
};
</script>
