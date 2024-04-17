<template>
  <v-layout class="gap-2" align-center>
    <v-btn
      :color="error ? 'error' : 'primary'"
      v-bind="$attrs"
      @click="$emit('click', $event)"
    >
      <slot>
        {{ $t('common.add') }}
      </slot>
    </v-btn>
    <v-fade-transition>
      <span
        v-show="error"
        class="error--text"
      >
        {{ error }}
      </span>
    </v-fade-transition>
  </v-layout>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    error: {
      type: String,
      default: '',
    },
  },
};
</script>
