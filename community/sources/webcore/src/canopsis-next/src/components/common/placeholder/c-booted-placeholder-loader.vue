<template>
  <div class="position-relative">
    <c-placeholder-loader v-if="!localReady" class="position-absolute" />
    <div v-if="booted" v-show="localReady">
      <slot />
    </div>
  </div>
</template>

<script>
import { ref, watch, nextTick } from 'vue';

export default {
  props: {
    booted: {
      type: Boolean,
      default: false,
    },
    ready: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const localReady = ref();

    const unwatch = watch(() => props.ready, (ready) => {
      if (ready) {
        localReady.value = ready;

        nextTick(() => unwatch());
      }
    }, { immediate: true });

    return { localReady };
  },
};
</script>

<style lang="scss" scoped>
.position-relative {
  width: 100%;
  height: 100%;
}
</style>
