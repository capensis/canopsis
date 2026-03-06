<template>
  <portal-target
    :name="$constants.PORTALS_NAMES.massActionsPanel"
    class="the-mass-actions-panels__portal-target"
    multiple
  />
</template>

<script>
export default {};
</script>

<style lang="scss" scoped>
.the-mass-actions-panels__portal-target {
  --actions-grid-columns-count: 12;
  --actions-grid-column-gap: 20px;
  --actions-grid-padding: 0 20px;

  position: absolute;
  right: 0;
  top: 0;
  width: 100%;
  padding: var(--actions-grid-padding);
  column-gap: var(--actions-grid-column-gap);
  display: grid;
  grid-template-columns: repeat(var(--actions-grid-columns-count), 1fr);
  pointer-events: none;
}
</style>
