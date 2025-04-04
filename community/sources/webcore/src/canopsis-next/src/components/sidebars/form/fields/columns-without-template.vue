<template>
  <widget-settings-item :title="label" :eager="false">
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <c-columns-field
      v-field="columns"
      v-bind="$attrs"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  inheritAttrs: false,
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
