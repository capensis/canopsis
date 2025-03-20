<template>
  <v-layout column>
    <external-data-table-general-info-form v-if="!form._id" :form="externalDataTable" />
    <c-id-field
      v-else
      :value="form._id"
      disabled
    />

    <v-text-field
      v-for="rule in externalDataTable.linked_rules"
      v-validate="'required'"
      v-field="form[rule._id]"
      :key="rule._id"
      :label="rule.name"
      :name="rule._id"
      :error-messages="errors.collect(rule._id)"
    />
  </v-layout>
</template>

<script>
import ExternalDataTableGeneralInfoForm
  from '@/components/other/external-data-table/form/external-data-table-general-info-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataTableGeneralInfoForm },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    externalDataTable: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
