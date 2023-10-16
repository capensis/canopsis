<template>
  <div class="fab">
    <v-layout column>
      <c-refresh-btn @click="$emit('refresh')" />
      <v-speed-dial
        v-if="hasAccess"
        v-model="fab"
        :transition="transition"
        direction="left"
      >
        <template #activator="">
          <v-btn
            v-model="fab"
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
      </v-speed-dial>
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
