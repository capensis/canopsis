<template>
  <v-list-group
    :disabled="disabled"
    class="widget-settings-group"
    color=""
    eager
  >
    <template #activator="{ props: activatorProps }">
      <v-list-item v-bind="activatorProps" :disabled="disabled">
        <template #title>
          <span :class="validationHeaderClass">
            {{ title }}
          </span>
        </template>
      </v-list-item>
    </template>

    <v-list ref="contentElement" class="widget-settings-group__list px-2 py-0">
      <slot />
    </v-list>
  </v-list-group>
</template>

<script>
import { ref } from 'vue';

import { useValidationHeader } from '@/hooks/validator/validation-header';

export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const contentElement = ref(null);
    const { validationHeaderClass } = useValidationHeader(contentElement);

    return {
      contentElement,
      validationHeaderClass,
    };
  },
};
</script>

<style lang="scss" scoped>
.widget-settings-group {
  &__list {
    background-color: rgba(var(--v-theme-on-surface), 0.04);
  }
}
</style>
