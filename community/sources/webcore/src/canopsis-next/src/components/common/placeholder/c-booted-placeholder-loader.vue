<template>
  <div class="c-booted-placeholder-loader position-relative">
    <c-placeholder-loader v-if="!allBooted" class="position-absolute" />
    <div v-if="booted" v-show="allBooted">
      <slot />
    </div>
  </div>
</template>

<script>
import {
  inject,
  onBeforeMount,
  onBeforeUnmount,
  ref,
} from 'vue';

export default {
  props: {
    asyncBootingProvider: {
      type: String,
      default: '$asyncBooting',
    },
    eager: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const booted = ref(false);
    const allBooted = ref(false);
    const asyncBooting = inject(props.asyncBootingProvider, null);
    const asyncBootingKey = Symbol('asyncBootingKey');

    const setBooted = () => {
      booted.value = true;
    };
    const setAllBooted = () => {
      allBooted.value = true;
    };

    onBeforeMount(() => {
      if (props.eager || !asyncBooting || asyncBooting.wasCalled) {
        setBooted();
        setAllBooted();

        return;
      }

      asyncBooting.register(asyncBootingKey, setBooted, setAllBooted);
    });

    onBeforeUnmount(() => asyncBooting?.unregister(asyncBootingKey));

    return {
      allBooted,
      booted,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-booted-placeholder-loader {
  width: 100%;
  height: 100%;
}
</style>
