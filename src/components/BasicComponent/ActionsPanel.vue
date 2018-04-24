<template lang="pug">
  div
    v-btn(icon @click.stop)
      v-icon(color="pink") delete
    v-btn(icon @click.stop)
      v-icon(color="green") edit
    slot
      v-menu(bottom left @click.native.stop)
       v-btn(icon slot="activator")
        v-icon more_vert
       v-list
         v-list-tile
         v-list-tile-title Click Me
</template>

<script>
export default {
  name: 'ActionsPanel',
  methods: {
  },
};
</script>

<style scoped>

</style>
