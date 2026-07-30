<template>
  <c-select
    :model-value="value"
    :items="preparedTypes"
    item-title="text"
    item-value="value"
    class="c-input-type-field"
    hide-details
    v-bind="$attrs"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #selection="{ item }">
      <v-icon :icon="getItemIcon(item)" class="c-input-type-field__icon" size="16" />
    </template>
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template #prepend>
          <v-icon :icon="getItemIcon(item)" class="c-input-type-field__icon" size="16" />
        </template>
        <template #title>
          {{ getItemText(item) }}
          <span v-if="getItemRaw(item)?.defined" class="ml-2 text-medium-emphasis">
            ({{ $t('common.defined') }})
          </span>
        </template>
      </v-list-item>
    </template>
  </c-select>
</template>

<script>
import { PATTERN_FIELD_TYPES } from '@/constants';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: String,
      default: '',
    },
    types: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    preparedTypes() {
      return this.types.map((type) =>
        type.text ? type : { ...type, text: this.$t(`common.mixedField.types.${type.value}`) },
      );
    },
  },
  methods: {
    getItemRaw(item) {
      return item?.raw ?? item;
    },

    getItemText(item) {
      return this.getItemRaw(item)?.text ?? '';
    },

    getItemIcon(item) {
      return this.getInputTypeIcon(this.getItemRaw(item)?.value);
    },

    getInputTypeIcon(type) {
      return (
        {
          [PATTERN_FIELD_TYPES.string]: 'title',
          [PATTERN_FIELD_TYPES.number]: 'exposure_plus_1',
          [PATTERN_FIELD_TYPES.boolean]: 'toggle_on',
          [PATTERN_FIELD_TYPES.null]: 'space_bar',
          [PATTERN_FIELD_TYPES.stringArray]: 'view_array',
          [PATTERN_FIELD_TYPES.timestamp]: 'calendar_today',
        }[type] ?? 'help_outline'
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.c-input-type-field {
  &__icon {
    color: inherit !important;
    opacity: 0.6;
  }
}
</style>
