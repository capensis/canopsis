<template>
  <v-layout class="gap-2">
    <v-flex xs6>
      <external-data-table-database-field
        :value="form.type"
        disabled
      />
    </v-flex>
    <v-flex xs6>
      <v-text-field
        :value="form.name"
        :label="$t('common.name')"
        disabled
      />
    </v-flex>
  </v-layout>
</template>

<script>
import ExternalDataTableDatabaseField
  from '@/components/other/external-data-table/form/fields/exterrnal-data-table-database-field.vue';

export default {
  components: { ExternalDataTableDatabaseField },
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
};
</script>
