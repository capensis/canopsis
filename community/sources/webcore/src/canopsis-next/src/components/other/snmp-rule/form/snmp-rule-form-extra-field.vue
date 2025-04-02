<template>
  <v-layout class="gap-3" column>
    <snmp-rule-form-field :label="$tc('common.customField', 2)" />
    <v-card v-for="(item, index) in form" :key="item.key">
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
            @click="remove(index)"
          />
        </v-layout>
        <snmp-rule-form-module-mib-objects-fields
          v-field="form[index].value"
          :items="items"
          large
        />
      </v-card-text>
    </v-card>
    <v-flex>
      <v-btn
        color="primary"
        @click="add"
      >
        {{ $t('common.add') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { snmpRuleExtraItemToForm } from '@/helpers/entities/snmp-rule/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

import SnmpRuleFormField from './snmp-rule-form-field-title.vue';
import SnmpRuleFormModuleMibObjectsFields from './snmp-rule-form-module-mib-objects-fields.vue';

export default {
  inject: ['$validator'],
  components: { SnmpRuleFormModuleMibObjectsFields, SnmpRuleFormField },
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
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const add = () => addItemIntoArray(snmpRuleExtraItemToForm());

    return {
      add,

      remove: removeItemFromArray,
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
