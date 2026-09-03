<template>
  <c-entity-field
    :value="value"
    :entity-types="entityTypes"
    :label="$t('common.connectorName')"
    :placeholder="!value.length ? $t('common.all') : ''"
    :persistent-placeholder="!value.length"
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
    <template #item="{ item, props }">
      <v-list-item v-bind="getItemProps(props)">
        <template #prepend="{ isSelected }">
          <v-checkbox-btn :model-value="isSelected || value.includes(item.name)" color="primary" />
        </template>
        <v-list-item-title>
          {{ item.name }}
        </v-list-item-title>
      </v-list-item>
    </template>
    <template #selection="{ item, index, props }">
      <v-chip
        v-if="index < showCount"
        v-bind="props"
        size="small"
        closable
      >
        <span class="text-truncate">
          {{ getItemName(item) }}
        </span>
      </v-chip>
      <span v-else-if="index === showCount">
        +{{ value.length - showCount }} {{ $t('common.more') }}
      </span>
      <span v-else />
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

  data() {
    return {
      entityTypes: [BASIC_ENTITY_TYPES.connector],
    };
  },

  methods: {
    getItemName(item) {
      return item?.name ?? item;
    },

    getItemProps(props) {
      const { title, ...listItemProps } = props;

      return listItemProps;
    },
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
