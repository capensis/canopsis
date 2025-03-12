<template>
  <v-layout column>
    <v-layout v-if="!form._id">
      <v-flex xs6>
        <external-data-table-database-field
          :value="externalDataTable.database"
          disabled
        />
      </v-flex>
      <v-flex xs6>
        <c-select-field
          :value="externalDataTable.name"
          :items="names"
          :label="$t('common.name')"
          disabled
        />
      </v-flex>
    </v-layout>
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
import { computed } from 'vue';

export default {
  inject: ['$validator'],
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
  setup(props) {
    const names = computed(() => [props.externalDataTable.name]);

    return {
      names,
    };
  },
};
</script>
