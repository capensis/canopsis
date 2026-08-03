<template>
  <c-form-block-array-field
    :form="form"
    :item-to-form="snmpRuleExtraItemToForm"
    :label="t('common.customField', 2)"
    :add-button-label="t('common.addCustomField')"
    @update:form="updateModel"
  >
    <template #item="{ item, index, remove }">
      <v-card>
        <v-card-text>
          <c-layout align="center">
            <c-text-field
              :model-value="form[index].name"
              :label="t('common.name')"
              :name="`${item.key}.name`"
              required
              @update:model-value="updateFieldInArrayItem(index, 'name', $event)"
            />

            <c-action-btn
              class="mr-0 ml-0"
              type="delete"
              @click="remove"
            />
          </c-layout>
          <snmp-rule-form-module-mib-objects-form
            :form="form[index].value"
            :items="items"
            :name="`${item.key}.value`"
            large
            @update:form="updateFieldInArrayItem(index, 'value', $event)"
          />
        </v-card-text>
      </v-card>
    </template>
  </c-form-block-array-field>
</template>

<script>
import { snmpRuleExtraItemToForm } from '@/helpers/entities/snmp-rule/form';

import { useI18n } from '@/hooks/i18n';
import { useArrayModelField } from '@/hooks/form/array-model-field';

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
    const {
      updateModel,
      updateFieldInArrayItem,
    } = useArrayModelField(props, emit);

    return {
      t,
      snmpRuleExtraItemToForm,
      updateModel,
      updateFieldInArrayItem,
    };
  },
};
</script>
