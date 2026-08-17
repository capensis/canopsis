<template>
  <widget-settings-item :title="label || $t('settings.treeOfDependenciesSettings')">
    <c-layout>
      <c-radio-group
        :model-value="modelValue"
        class="mt-0"
        name="opened"
        hide-details
        mandatory
        @update:model-value="$emit('update:modelValue', $event)"
      >
        <v-radio
          v-for="type in types"
          :key="type.value"
          :label="type.label"
          :value="type.value"
          color="primary"
        />
      </c-radio-group>
    </c-layout>
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import { TREE_OF_DEPENDENCIES_SHOW_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Number,
      default: TREE_OF_DEPENDENCIES_SHOW_TYPES.custom,
    },
    label: {
      type: String,
      default: '',
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    const types = computed(() =>
      Object.values(TREE_OF_DEPENDENCIES_SHOW_TYPES).map(value => ({
        value,
        label: t(`entity.treeOfDependenciesShowTypes.${value}`),
      })),
    );

    return {
      types,
    };
  },
};
</script>
