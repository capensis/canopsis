<template>
  <c-layout gap="2" align="center" class="c-btn-with-error">
    <v-btn v-bind="$attrs" :color="error ? 'error' : 'primary'" class="c-btn-with-error__button">
      <slot />
    </v-btn>
    <v-fade-transition>
      <span v-show="error" class="c-btn-with-error__message text-error">
        {{ error }}
      </span>
    </v-fade-transition>
  </c-layout>
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

<style lang="scss" scoped>
.c-btn-with-error {
  min-width: 0;

  &__button {
    flex: 0 0 auto;
  }

  &__message {
    min-width: 0;
  }
}
</style>
