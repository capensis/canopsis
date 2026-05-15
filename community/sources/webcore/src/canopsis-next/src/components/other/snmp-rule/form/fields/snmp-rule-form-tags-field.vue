<template>
  <c-form-block-array-field
    v-field="form"
    :item-to-form="snmpRuleTagToForm"
    :label="$tc('common.tag', 2)"
  >
    <template #item="{ index, remove }">
      <v-layout align-center>
        <snmp-rule-form-module-mib-objects-form
          v-field="form[index]"
          :items="items"
          large
        />
        <c-action-btn
          class="mr-0 ml-0"
          type="delete"
          @click="remove"
        />
      </v-layout>
    </template>
  </c-form-block-array-field>
</template>

<script>
import { snmpRuleTagToForm } from '@/helpers/entities/snmp-rule/form';

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
      snmpRuleTagToForm,
    };
  },
};
</script>

<style lang="scss" scoped>
.v-btn.active {
  &:hover:before {
    opacity: .16;
  }

  &:before {
    background-color: currentColor;
  }
}

.vars-input ::v-deep .v-input__slot {
  height: 56px;
}
</style>
