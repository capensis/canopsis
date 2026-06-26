<template>
  <v-toolbar
    class="groups-top-bar"
    density="compact"
  >
    <v-toolbar-title class="text-body-2">
      Groups top bar
    </v-toolbar-title>
  </v-toolbar>
</template>

<script>
export default {
  name: 'GroupsTopBar',
}
</script>

<style lang="scss" scoped>
.groups-top-bar {
  background: color-mix(in srgb, rgb(var(--v-theme-primary)) 92%, black);
  color: white;
}
</style>
