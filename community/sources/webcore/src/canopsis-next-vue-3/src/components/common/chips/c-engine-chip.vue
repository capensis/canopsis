<template>
  <div
    :style="{ background: color }"
    :class="{ 'cursor-pointer': clickable }"
    class="c-engine-chip text-white font-weight-bold"
  >
    <slot />
  </div>
</template>

<script>
export default {
  name: 'CEngineChip',

  props: {
    color: {
      type: String,
      default: undefined,
    },
    clickable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.c-engine-chip {
  align-items: center;
  border: 1px solid white;
  border-radius: 20px;
  display: inline-flex;
  font-size: 12px;
  height: 30px;
  justify-content: center;
  min-width: 60px;
  padding: 3px 7px;
}
</style>
