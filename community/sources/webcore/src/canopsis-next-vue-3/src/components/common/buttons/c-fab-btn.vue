<template>
  <div class="fab ma-2">
    <c-layout gap="3" column>
      <c-refresh-btn @click="$emit('refresh')" />

      <slot name="additional" />

      <v-tooltip v-if="hasAccess && hasCreateListener" location="left">
        <template #activator="{ props: activatorProps }">
          <v-btn
            :theme="dark ? 'dark' : undefined"
            color="primary"
            size="x-large"
            icon
            v-bind="activatorProps"
            @click.stop="$emit('create')"
          >
            <c-icon>add</c-icon>
          </v-btn>
        </template>
        <slot />
      </v-tooltip>
    </c-layout>
  </div>
</template>

<script>
export default {
  props: {
    hasAccess: {
      type: Boolean,
      default: true,
    },
    hasCreateListener: {
      type: Boolean,
      default: true,
    },
    dark: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['refresh', 'create'],
};
</script>
