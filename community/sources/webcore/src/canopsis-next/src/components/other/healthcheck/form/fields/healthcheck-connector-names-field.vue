<template>
  <c-entity-field
    :value="value"
    :entity-types="entityTypes"
    :label="$t('common.connectorName')"
    :placeholder="!value.length ? $t('common.all') : ''"
    :persistent-placeholder="!value.length"
    :show-count="showCount"
    item-text="name"
    item-value="name"
    class="healthcheck-connector-names-field"
    density="compact"
    multiple
    clearable
    autocomplete
    hide-details
    @input="$emit('input', $event)"
  >
    <template #item="{ item, props, search }">
      <v-list-item v-bind="getItemProps(props)">
        <template #prepend="{ isSelected }">
          <v-checkbox-btn :model-value="isSelected || value.includes(item.name)" color="primary" />
        </template>
        <v-list-item-title>
          <v-list-item-mask :text="item.name" :mask="search" />
        </v-list-item-title>
      </v-list-item>
    </template>
  </c-entity-field>
</template>

<script>
import { BASIC_ENTITY_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: Array,
      required: true,
    },
    showCount: {
      type: Number,
      default: 2,
    },
  },

  emits: ['input'],

  setup() {
    const entityTypes = [BASIC_ENTITY_TYPES.connector];

    const getItemProps = (itemProps) => {
      const { title, ...listItemProps } = itemProps;

      return listItemProps;
    };

    return {
      entityTypes,
      getItemProps,
    };
  },
};
</script>

<style lang="scss">
.healthcheck-connector-names-field {
  & ::placeholder {
    color: currentColor !important;
    opacity: 1;
  }
}
</style>
