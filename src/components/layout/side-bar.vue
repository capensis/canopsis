<template lang="pug">
  v-navigation-drawer(
    :value='isOpen'
    absolute
    app
    hide-overlay
    temporary
    stateless
    clipped
  )
    v-list
      v-list-tile
        v-list-tile-action
          v-icon home
        v-list-tile-action
          router-link(to='/') Home
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
    isOpen: {
      type: Boolean,
      required: true,
    },
  },
};
</script>

<style scoped>

</style>
