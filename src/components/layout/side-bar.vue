<template lang="pug">
  v-navigation-drawer(
    v-model='isOpen'
    absolute
    app
    :clipped="windowSize.x < 1264 ? false : true"
  )
    v-list
      v-list-tile
        v-list-tile-action
          v-icon home
        v-list-tile-action
          router-link(to='/') {{$t('common.home')}}
      v-list-group(value='true')
        v-list-tile(
          slot='activator'
        )
          v-list-tile-title View Group 1
        v-list-tile
          v-list-tile-title SubView 1
        v-list-tile
          v-list-tile-title SubView 2
      v-list-group(value='true')
        v-list-tile(
          slot='activator'
        )
          v-list-tile-title View Group 2
        v-list-tile
          v-list-tile-title SubView 1
        v-list-tile
          v-list-tile-title SubView 2

</template>

<script>
export default {
  name: 'SideBar',
  props: {
    windowSize: {
      type: Object,
    },
  },
  computed: {
    isOpen: {
      get() {
        return this.$store.state.app.isSideBarOpen;
      },
      set(state) {
        if (state !== this.$store.state.app.isSideBarOpen) {
          this.$store.dispatch('app/toggleSideBar');
        }
      },
    },
  },
};
</script>

<style scoped>

</style>
