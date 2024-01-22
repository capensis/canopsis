<template>
  <div class="fab ma-2">
    <v-layout column>
      <v-flex class="mb-3">
        <c-refresh-btn @click="$emit('refresh')" />
      </v-flex>
      <c-speed-dial
        v-if="hasAccess"
        :transition="transition"
        direction="left"
      >
        <template #activator="{ bind }">
          <v-btn
            v-bind="bind"
            color="primary"
            fab
          >
            <slot name="icons">
              <v-icon>add</v-icon>
              <v-icon>close</v-icon>
            </slot>
          </v-btn>
        </template>
        <slot />
      </c-speed-dial>
    </v-layout>
  </div>
</template>

<script>
export default {
  props: {
    transition: {
      type: String,
      default: 'slide-y-reverse-transition',
    },
    hasAccess: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      fab: false,
    };
  },
};
</script>
