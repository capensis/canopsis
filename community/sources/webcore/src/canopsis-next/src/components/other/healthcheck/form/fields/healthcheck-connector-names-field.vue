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
    multiple
    clearable
    @input="$emit('input', $event)"
  >
    <template #item="{ item, attrs, on, parent}">
      <v-list-item
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action>
          <v-checkbox
            :input-value="attrs.inputValue"
            :color="parent.color"
            hide-details
            small
          />
        </v-list-item-action>
        <v-list-item-content>
          {{ item.name }}
        </v-list-item-content>
      </v-list-item>
    </template>
    <template #selection="{ item, index, parent }">
      <v-chip
        v-if="index < showCount"
        small
        close
        @click:close="parent.onChipInput(item)"
      >
        <span class="text-truncate">{{ item }}</span>
      </v-chip>
      <span v-else-if="index === showCount">+{{ value.length - showCount }} {{ $t('common.more') }}</span>
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
  setup() {
    const entityTypes = [BASIC_ENTITY_TYPES.connector];

    return {
      entityTypes,
    };
  },
};
</script>

<style lang="scss">
.healthcheck-connector-names-field {
  & ::placeholder {
    .theme--dark & {
      color: var(--v-text-dark-primary) !important;
    }

    .theme--light & {
      color: var(--v-text-light-primary) !important;
    }
  }

  &.v-autocomplete input {
    padding: 8px 0 !important;
  }
}
</style>
