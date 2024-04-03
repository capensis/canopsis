<template>
  <div
    :style="{ background: color }"
    :class="{ 'cursor-pointer': !!$listeners.click }"
    class="c-engine-chip white--text font-weight-bold"
    v-on="$listeners"
  >
    <slot />
  </div>
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss">
.c-engine-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 1px solid white;
  border-radius: 20px;
  height: 30px;
  padding: 3px 7px;
  font-size: 12px;
  min-width: 60px;
}
</style>
