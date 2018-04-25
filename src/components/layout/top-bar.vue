<template lang="pug">
    v-toolbar(
      dense
      fixed
      clipped-left
      app
      color='blue darken-4'
    )
      v-toolbar-side-icon(
        class='white--text'
        @click.stop="handleClick"
      )
</template>

<script>
export default {
  name: 'TopBar',
  methods: {
    handleClick() {
      this.$store.dispatch('app/toggleSideBar');
    },
  },
};
</script>
