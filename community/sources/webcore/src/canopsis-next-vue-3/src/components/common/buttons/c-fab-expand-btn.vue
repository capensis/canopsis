<script setup>
const emit = defineEmits(['refresh']);

defineProps({
  transition: {
    type: String,
    default: 'slide-y-reverse-transition',
  },
  hasAccess: {
    type: Boolean,
    default: true,
  },
});
</script>

<template>
  <div class="fab ma-2">
    <c-layout column>
      <div class="mb-3">
        <c-refresh-btn @click="emit('refresh')" />
      </div>

      <v-speed-dial
        v-if="hasAccess"
        location="left center"
        persistent
        attach
      >
        <template #activator="{ props: activatorProps, isActive }">
          <v-btn
            v-bind="activatorProps"
            :class="{ 'c-fab--active': isActive }"
            color="primary"
            class="c-fab"
            icon
          >
            <slot
              :is-active="isActive"
              name="icons"
            >
              <v-icon size="large" icon="add" />
            </slot>
          </v-btn>
        </template>

        <slot />
      </v-speed-dial>
    </c-layout>
  </div>
</template>

<style lang="scss" scoped>
.c-fab {
  --v-btn-height: 56px;

  width: 56px;
  height: 56px;
  border-radius: 50%;

  .v-icon {
    transition: transform 0.2s ease;
  }

  &--active {
    .v-icon {
      transform: rotate(45deg);
    }
  }
}
</style>
