<template>
  <c-form-block-array-field
    v-field="form"
    :item-to-form="snmpRuleExtraItemToForm"
    :label="$tc('common.customField', 2)"
    :add-button-label="$t('common.addCustomField')"
  >
    <template #item="{ item, index, remove }">
      <v-card>
        <v-card-text>
          <v-layout align-center>
            <v-text-field
              v-field="form[index].name"
              v-validate="'required'"
              :label="$t('common.name')"
              :name="`${item.key}.name`"
              :error-messages="errors.collect(`${item.key}.name`)"
            />

            <c-action-btn
              class="mr-0 ml-0"
              type="delete"
              @click="remove"
            />
          </v-layout>
          <snmp-rule-form-module-mib-objects-form
            v-field="form[index].value"
            :items="items"
            large
          />
        </v-card-text>
      </v-card>
    </template>
  </c-form-block-array-field>
</template>

<script>
import { snmpRuleExtraItemToForm } from '@/helpers/entities/snmp-rule/form';

import SnmpRuleFormModuleMibObjectsForm from '../snmp-rule-form-module-mib-objects-form.vue';

export default {
  inject: ['$validator'],
  components: { SnmpRuleFormModuleMibObjectsForm },
  model: {
    prop: 'form',
    event: 'input',
  },
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
  setup() {
    return {
      snmpRuleExtraItemToForm,
    };
  },
};
</script>
