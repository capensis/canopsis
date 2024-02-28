<template>
  <div class="impact-state-indicator">
    <span class="impact-state-text white--text">{{ value }}</span>
    <div class="pointer-wrapper">
      <div
        :style="pointerStyle"
        class="pointer"
      />
      <div class="section-wrapper">
        <div
          v-for="color in $config.COLORS.impactState"
          :key="color"
          :style="{ backgroundColor: color }"
          class="section"
        />
      </div>
    </div>
  </div>
</template>

<script>
/**
 * Half of width of triangle -1px for including sections wrapper border
 *
 * @type {number}
 */
const POINTER_SHIFT = 4;

export default {
  props: {
    value: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    pointerStyle() {
      return { left: `${this.value - POINTER_SHIFT}px` };
    },
  },
};
</script>

<style lang="scss" scoped>
$fontSize: 12px;
$fontWeight: 500;

$borderColor: white;
$borderHeight: 5px;

$sectionWidth: 1px;
$sectionHeight: 10px;

.impact-state-indicator {
  position: relative;
  display: flex;
  align-items: center;
  padding-top: $borderHeight;

  .impact-state-text {
    margin-right: 3px;
    font-size: $fontSize;
    line-height: $fontSize;
    font-weight: $fontWeight;
  }
}

.pointer {
  position: absolute;
  bottom: 100%;
  display: block;
  width: 0;
  height: 0;
  border-left: $borderHeight solid transparent;
  border-right: $borderHeight solid transparent;
  border-top: $borderHeight solid $borderColor;

  &-wrapper {
    position: relative;
  }
}

.section {
  margin-bottom: -1px;
  width: $sectionWidth;
  height: $sectionHeight;

  &-wrapper {
    display: flex;
    border-radius: 3px;
    overflow: hidden;
    border: 2px solid $borderColor;
  }
}
</style>
