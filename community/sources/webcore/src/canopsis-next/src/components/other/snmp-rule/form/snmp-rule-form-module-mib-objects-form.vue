<template>
  <div>
    <snmp-rule-form-field :label="label" />
    <snmp-rule-form-module-mib-objects-fields
      v-field="form"
      :items="items"
      :large="large"
      :name="name"
      :required="required"
    />
  </div>
</template>

<script>
import SnmpRuleFormField from './snmp-rule-form-field-title.vue';
import SnmpRuleFormModuleMibObjectsFields from './snmp-rule-form-module-mib-objects-fields.vue';

export default {
  components: { SnmpRuleFormModuleMibObjectsFields, SnmpRuleFormField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    large: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
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
