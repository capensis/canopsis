<template>
  <v-layout class="gap-3" column>
    <c-enabled-field v-field="form.enabled" with-background />
    <c-entity-field
      v-field="form.name"
      :label="$t('common.connectorName')"
      :entity-types="entityTypes"
      item-text="name"
      item-value="name"
      required
    />
  </v-layout>
</template>

<script>
import { BASIC_ENTITY_TYPES } from '@/constants';

export default {
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
  setup() {
    const entityTypes = [BASIC_ENTITY_TYPES.connector];

    return {
      entityTypes,
    };
  },
};
</script>
