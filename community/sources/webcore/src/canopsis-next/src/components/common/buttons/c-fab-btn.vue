<template>
  <div class="fab ma-2">
    <c-layout gap="3" column>
      <c-refresh-btn @click="$emit('refresh')" />

      <slot name="additional" />

      <c-tooltip v-if="hasAccess && hasCreateListener" location="left">
        <template #activator="{ props: activatorProps }">
          <v-btn
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
      </c-tooltip>
    </c-layout>
  </div>
</template>

<script>
import { computed, getCurrentInstance, onBeforeUpdate, ref } from 'vue';

export default {
  props: {
    hasAccess: {
      type: Boolean,
      default: true,
    },
    dark: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['refresh', 'create'],

  setup() {
    const instance = getCurrentInstance();
    const updateTick = ref(0);

    onBeforeUpdate(() => {
      updateTick.value += 1;
    });

    const hasCreateListener = computed(() => Boolean(updateTick.value >= 0 && instance?.vnode.props?.onCreate),
    );

    return {
      hasCreateListener,
    };
  },
};
</script>
