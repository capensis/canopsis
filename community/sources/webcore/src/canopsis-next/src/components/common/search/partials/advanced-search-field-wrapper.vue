<template>
  <div
    :class="[wrapperClass, themeClasses]"
    class="v-input v-input--hide-details v-text-field v-text-field--single-line
     v-text-field--is-booted v-select primary--text"
  >
    <div class="v-input__control">
      <slot name="activator" />
      <slot name="items" />
    </div>
  </div>
</template>

<script>
import { computed } from 'vue';
import Themeable from 'vuetify/lib/mixins/themeable';

export default {
  mixins: [Themeable],
  props: {
    isMenuActive: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const wrapperClass = computed(() => ({
      'v-input--is-focused': props.isMenuActive,
    }));

    return {
      wrapperClass,
    };
  },
};
</script>
