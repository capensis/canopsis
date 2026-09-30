<template>
  <c-autocomplete
    :model-value="itemValue"
    :items="items"
    :disabled="disabled"
    :label="label || t('common.searchBy')"
    :name="name"
    :required="required"
    :item-props="getItemProps"
    :item-title="getItemTitle"
    item-value="value"
    @update:model-value="updateValue"
  >
    <template #item="{ props: itemProps, item }">
      <v-list-item v-bind="itemProps">
        <v-list-item-title>
          {{ getItemTitle(item) }}
        </v-list-item-title>

        <template v-if="item?.alias" #append>
          <c-tooltip :content="`infos.${item.originalValue}.value`" top>
            <template #activator="{ props: activatorProps }">
              <v-icon
                class="ml-1"
                color="primary"
                size="small"
                v-bind="activatorProps"
              >
                alternate_email
              </v-icon>
            </template>
          </c-tooltip>
        </template>
      </v-list-item>
    </template>
    <template v-if="selectedItem?.alias" #prepend-inner>
      <v-icon class="mt-1" color="primary" size="small">
        alternate_email
      </v-icon>
    </template>
  </c-autocomplete>
</template>

<script>
import { isNil } from 'lodash';
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  emits: ['update:value'],
  props: {
    value: {
      type: [String, Object],
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'attribute',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const itemsByValue = computed(() => Object.fromEntries(props.items.map(item => [item.value, item])));

    const itemValue = computed(() => {
      const value = props.value?.value ?? props.value;

      if (isNil(value) || value === '') {
        return null;
      }

      return itemsByValue.value[value]?.value ?? value;
    });

    const selectedItem = computed(() => {
      if (isNil(itemValue.value) || itemValue.value === '') {
        return null;
      }

      return itemsByValue.value[itemValue.value] ?? null;
    });

    const getItemTitle = item => item?.text ?? '';

    const getItemProps = item => ({
      disabled: item?.options?.disabled,
    });

    const updateValue = (value) => {
      const selected = itemsByValue.value[value?.value ?? value] ?? value;

      if (!props.returnObject) {
        emit('update:value', selected?.value ?? selected);

        return;
      }

      emit('update:value', selected);
    };

    return {
      t,
      itemValue,
      selectedItem,
      getItemTitle,
      getItemProps,
      updateValue,
    };
  },
};
</script>
