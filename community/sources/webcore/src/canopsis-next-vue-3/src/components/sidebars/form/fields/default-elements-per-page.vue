<template>
  <widget-settings-item
    :title="$t('settings.defaultNumberOfElementsPerPage')"
    :sub-title="subTitle"
  >
    <c-items-per-page-field
      :model-value="modelValue"
      :label="$t('settings.elementsPerPage')"
      class="default-elements-per-page__field"
      hide-details
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { PAGINATION_LIMIT } from '@/config';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    subTitle: {
      type: String,
      default: '',
    },
  },

  emits: ['update:modelValue'],
};
</script>

<style lang="scss" scoped>
.default-elements-per-page__field {
  width: 100%;
  max-width: none;
}
</style>
