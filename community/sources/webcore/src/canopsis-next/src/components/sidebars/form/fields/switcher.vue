<template>
  <widget-settings-flat-item :title="title">
    <template #actions>
      <v-switch
        class="ma-0 pa-0"
        :input-value="value"
        color="primary"
        hide-details
        @change="$emit('input', $event)"
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
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
