<template>
  <c-layout gap="3" column>
    <c-enabled-field v-model="enabled" with-background />

    <c-entity-field
      v-model:value="name"
      :label="$t('common.connectorName')"
      :entity-types="entityTypes"
      item-text="name"
      item-value="name"
      required
    />
  </c-layout>
</template>

<script>
import { BASIC_ENTITY_TYPES } from '@/constants';
import { useFormModels } from '@/hooks/form';

export default {
  setup() {
    const { name, enabled } = useFormModels(['name', 'enabled']);
    const entityTypes = [BASIC_ENTITY_TYPES.connector];

    return {
      name,
      enabled,
      entityTypes,
    };
  },
};
</script>
