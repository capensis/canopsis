<template>
  <widget-settings-flat-item :title="title">
    <template #actions>
      <v-switch
        :model-value="modelValue"
        class="ma-0 pa-0"
        color="primary"
        hide-details
        @update:model-value="$emit('update:modelValue', $event)"
      />
    </template>
  </widget-settings-flat-item>
</template>

<script>
import WidgetSettingsFlatItem from '@/components/sidebars/partials/widget-settings-flat-item.vue';

export default {
  components: { WidgetSettingsFlatItem },

  props: {
    title: {
      type: String,
      default: '',
    },
    modelValue: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
