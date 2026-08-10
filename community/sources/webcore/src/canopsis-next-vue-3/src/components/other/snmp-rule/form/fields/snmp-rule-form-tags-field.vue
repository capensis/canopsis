<template>
  <c-form-block-array-field
    :form="form"
    :item-to-form="snmpRuleTagToForm"
    :label="t('common.tag', 2)"
    :add-button-label="t('common.addTag')"
    @update:form="updateModel"
  >
    <template #item="{ index, remove }">
      <c-layout align="center">
        <snmp-rule-form-module-mib-objects-form
          :form="form[index]"
          :items="items"
          :name="`tags.${index}`"
          large
          @update:form="updateItemInArray(index, $event)"
        />
        <c-action-btn class="mr-0 ml-0" type="delete" @click="remove" />
      </c-layout>
    </template>
  </c-form-block-array-field>
</template>

<script>
import { snmpRuleTagToForm } from '@/helpers/entities/snmp-rule/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

import SnmpRuleFormModuleMibObjectsForm from '../snmp-rule-form-module-mib-objects-form.vue';

export default {
  components: { SnmpRuleFormModuleMibObjectsForm },

  props: {
    form: {
      type: Array,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel, updateItemInArray } = useArrayModelField(props, emit);

    return {
      t,
      snmpRuleTagToForm,
      updateModel,
      updateItemInArray,
    };
  },
};
</script>
