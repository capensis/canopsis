<template>
  <v-layout column>
    <external-data-table-general-info-form v-if="!form._id" :form="externalDataTable" />
    <c-id-field
      v-else
      :value="form._id"
      disabled
    />

    <v-text-field
      v-for="column in externalDataTable.column_configs"
      v-validate="'required'"
      v-field="form[column.name]"
      :key="column.name"
      :label="column.name"
      :name="column.name"
      :error-messages="errors.collect(column.name)"
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
