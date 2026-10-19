<template>
  <widget-settings-item>
    <template #title>
      <c-layout align="center" gap="2">
        <span>{{ title }}</span>
        <v-icon :icon="sizeIcon" />
      </c-layout>
    </template>

    <c-column-size-field
      :model-value="modelValue"
      :name="fieldName"
      :mobile="mobile"
      :tablet="tablet"
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Number,
      required: false,
    },
    title: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      default: '',
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

  emits: ['update:modelValue'],

  computed: {
    sizeIcon() {
      if (this.mobile) {
        return 'phone_android';
      }

      if (this.tablet) {
        return 'tablet';
      }

      return 'laptop';
    },

    fieldName() {
      if (this.name) {
        return this.name;
      }

      if (this.mobile) {
        return 'columnMobile';
      }

      if (this.tablet) {
        return 'columnTablet';
      }

      return 'columnDesktop';
    },
  },
};
</script>
