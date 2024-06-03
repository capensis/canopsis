<template>
  <div class="position-relative">
    <c-placeholder-loader v-if="!allBooted" class="position-absolute" />
    <div v-if="booted" v-show="allBooted">
      <slot />
    </div>
  </div>
</template>

<script>
import { ref, inject, onBeforeUnmount } from 'vue';

export default {
  props: {
    asyncBootingProvider: {
      type: String,
      default: '$asyncBooting',
    },
  },
  setup(props) {
    const booted = ref(false);
    const allBooted = ref(false);
    const asyncBooting = inject(props.asyncBootingProvider);
    const asyncBootingKey = Symbol('asyncBootingKey');

    const setBooted = () => booted.value = true;
    const setAllBooted = () => allBooted.value = true;

    asyncBooting.register(asyncBootingKey, setBooted, setAllBooted);

    onBeforeUnmount(() => asyncBooting.unregister(asyncBootingKey));

    return {
      booted,
      allBooted,
    };
  },
};
</script>

<style lang="scss" scoped>
.position-relative {
  width: 100%;
  height: 100%;
}
</style>
