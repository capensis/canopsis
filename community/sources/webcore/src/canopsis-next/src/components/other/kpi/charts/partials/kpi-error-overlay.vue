<template>
  <v-layout
    class="overlay"
    align-center
    justify-center
  >
    <span class="error--text">{{ $t('kpi.noData') }}</span>
  </v-layout>
</template>

<style lang="scss" scoped>
.overlay {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  right: 0;
  background: rgba(255, 255, 255, 0.7);
}
</style>
