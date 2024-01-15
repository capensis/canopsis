<template>
  <span class="c-circle-badge">
    <slot />
  </span>
</template>

<style lang="scss">
.c-circle-badge {
  --c-circle-badge-size: 20px;
  --c-circle-badge-font-size: calc(var(--c-circle-badge-size) * 0.6);

  color: white;

  padding: 4px;

  display: inline-flex;
  align-items: center;
  justify-content: center;

  font-size: var(--c-circle-badge-font-size);

  height: var(--c-circle-badge-size);
  min-width: var(--c-circle-badge-size);
  border-radius: calc(var(--c-circle-badge-size) / 2);
}
</style>
