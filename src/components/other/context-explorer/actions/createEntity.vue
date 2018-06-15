<template lang="pug">
v-speed-dial(
  :direction="top",
  :open-on-hover="true",
  transition="scale-transition")
  v-btn(slot='activator', v-model='fab', color='blue darken-2', dark, fab)
    v-icon add
    v-icon close
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, small, color='indigo')
      v-icon watch
    span Watcher
  v-tooltip(left)
    v-btn(slot="activator", fab, dark, small, color='red')
      v-icon perm_identity
    span Entities
</template>

<script>
export default {
};
</script>
