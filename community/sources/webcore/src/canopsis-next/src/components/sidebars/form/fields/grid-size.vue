<template>
  <widget-settings-item>
    <template #title="">
      <v-layout align-center>
        <span>{{ title }}</span>
        <v-icon class="ml-2">
          {{ sizeIcon }}
        </v-icon>
      </v-layout>
    </template>
    <c-column-size-field
      v-field="value"
      :mobile="mobile"
      :tablet="tablet"
    />
  </widget-settings-item>
</template>

<script>
import { computed } from 'vue';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Number,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    mobile: {
      type: Boolean,
      default: false,
    },
    tablet: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const sizeIcon = computed(() => {
      if (props.mobile) {
        return 'phone_android';
      }

      if (props.tablet) {
        return 'tablet';
      }

      return 'laptop';
    });

    return {
      sizeIcon,
    };
  },
};
</script>
