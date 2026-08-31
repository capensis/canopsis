<template>
  <c-autocomplete
    :model-value="modelValue"
    :items="items"
    :label="label || t('common.icon')"
    :placeholder="placeholder"
    :name="name"
    :required="required"
    :disabled="disabled"
    item-title="text"
    item-value="value"
    hide-no-data
    clearable
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #selection="{ item }">
      <v-icon v-if="getItemValue(item)" class="mr-2" size="18">
        {{ getItemValue(item) }}
      </v-icon>
      <span>{{ getItemText(item) }}</span>
    </template>

    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <template #prepend>
          <v-icon size="18">
            {{ getItemValue(item) }}
          </v-icon>
        </template>
      </v-list-item>
    </template>

    <template v-if="$slots['append-outer']" #append-outer>
      <slot name="append-outer" />
    </template>

    <template v-if="$slots['no-data']" #no-data>
      <slot name="no-data" />
    </template>
  </c-autocomplete>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import materialIcons from '@/assets/material-icons/MaterialIcons.json';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon_name',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const items = computed(() => materialIcons.map(icon => ({
        text: icon,
        value: icon,
      })),
    );

    const getItemValue = item => item?.value ?? '';
    const getItemText = item => item?.text ?? getItemValue(item);

    return {
      t,
      items,
      getItemValue,
      getItemText,
    };
  },
};
</script>
