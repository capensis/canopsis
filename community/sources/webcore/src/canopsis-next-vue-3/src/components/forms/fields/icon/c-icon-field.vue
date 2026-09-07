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
    <template v-if="modelValue" #prepend-inner>
      <span class="c-icon-field__icon">
        <v-icon size="18">
          {{ getItemIcon({ value: modelValue }) }}
        </v-icon>
      </span>
    </template>

    <template #item="{ props: itemProps, item }">
      <v-divider v-if="isDividerItem(item)" />
      <v-list-item v-else v-bind="itemProps">
        <template #prepend>
          <span class="c-icon-field__icon">
            <v-icon size="18">
              {{ getItemIcon(item) }}
            </v-icon>
          </span>
        </template>
        <template #title>
          {{ getItemText(item) }}
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
import { useIcon } from '@/hooks/store/modules/icon';

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
    const { registeredIconsById } = useIcon();

    const registeredIconsItems = computed(() => Object.values(registeredIconsById.value ?? {}).map(({ title }) => ({
        text: title,
        value: `$vuetify.icon.${title}`,
      })),
    );
    const materialIconsItems = computed(() => materialIcons.map(icon => ({
        text: icon,
        value: icon,
      })),
    );
    const items = computed(() => {
      if (!registeredIconsItems.value.length) {
        return materialIconsItems.value;
      }

      return [
        ...registeredIconsItems.value,
        { type: 'divider', divider: true, text: '', value: '__icon-divider__' },
        ...materialIconsItems.value,
      ];
    });

    const isDividerItem = item => item?.divider || item?.type === 'divider';
    const getItemValue = item => item?.value ?? item?.raw?.value ?? '';
    const getItemText = item => item?.text ?? item?.title ?? item?.raw?.text ?? getItemValue(item);
    const getItemIcon = (item) => {
      const value = getItemValue(item);

      if (!value || isDividerItem(item)) {
        return '';
      }

      if (value.startsWith('$vuetify.icon.')) {
        return `canopsis:${value.slice('$vuetify.icon.'.length)}`;
      }

      return value;
    };

    return {
      t,
      items,
      isDividerItem,
      getItemValue,
      getItemText,
      getItemIcon,
    };
  },
};
</script>

<style lang="scss">
.c-icon-field__icon {
  display: inline-flex;
  max-width: 24px;
  max-height: 24px;
  overflow: hidden;
}
</style>
