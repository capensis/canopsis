<template lang="pug">
  div
    div(v-if="$mq === 'laptop'")
        v-btn(icon @click.stop)
          v-icon(color="pink") delete
        v-btn(icon @click.stop)
          v-icon(color="green") edit
        slot
          v-menu(bottom, left, @click.native.stop)
           v-btn(icon, slot="activator")
            v-icon more_vert
           v-list
             v-list-tile
             v-list-tile-title BONJOUR
    v-menu(bottom, left, @click.native.stop, v-if="$mq === 'mobile' || $mq === 'tablet'")
        v-btn(icon slot="activator")
          v-icon more_vert
        v-list
          v-list-tile
          v-list-tile-title BONJOUR
</template>

<script>
export default {
  name: 'ActionsPanel',
};
</script>

<style scoped>
</style>
