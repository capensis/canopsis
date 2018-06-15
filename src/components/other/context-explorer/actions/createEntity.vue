<template lang="pug">
v-speed-dial(
  v-model="fab",
  :direction="direction",
  :open-on-hover="hover",
  :transition="transition")
  v-btn(slot='activator', v-model='fab', color='blue darken-2', dark='', fab='')
    v-icon account_circle
    v-icon close
  v-btn(fab='', dark='', small='', color='green')
    v-icon edit
  v-btn(fab='', dark='', small='', color='indigo')
    v-icon add
  v-btn(fab='', dark='', small='', color='red')
    v-icon delete
</template>

<script>
export default {
  data: () => ({
    direction: 'top',
    fab: false,
    fling: false,
    hover: true,
    tabs: null,
    top: false,
    right: false,
    bottom: true,
    left: true,
    transition: 'slide-y-reverse-transition',
  }),
};
</script>

<style scoped>

</style>
