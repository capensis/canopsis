<template>
  <div class="fab ma-2">
    <v-layout column>
      <v-flex class="mb-3">
        <c-refresh-btn @click="$emit('refresh')" />
      </v-flex>
      <slot name="additional" />
      <v-tooltip
        v-if="hasAccess && $listeners.create"
        left
      >
        <template #activator="{ on }">
          <v-btn
            v-on="on"
            :dark="dark"
            color="primary"
            fab
            @click.stop="$listeners.create"
          >
            <v-icon>add</v-icon>
          </v-btn>
        </template>
        <slot />
      </v-tooltip>
    </v-layout>
  </div>
</template>

<script>
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
  data() {
    return {
      fab: false,
    };
  },
};
</script>
