<template lang="pug">
  v-layout
    v-layout(v-if="$mq === 'laptop'")
      v-btn(icon @click.stop)
        v-icon(color="pink") delete
      v-btn(icon @click.stop)
        v-icon(color="green") edit
      slot
    v-menu(bottom, left, @click.native.stop)
      v-btn(icon, slot="activator")
        v-icon more_vert
      v-list(class="pa-3")
        v-list-tile
        v-list-tile-title(@click="handleMoreInfosClick") More infos
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');

export default {
  name: 'ActionsPanel',
  props: {
    alarm: {
      type: Object,
    },
  },
  methods: {
    ...modalMapActions({ openMoreInfosModal: 'show' }),

    handleMoreInfosClick() {
      this.openMoreInfosModal({ name: 'more-infos', config: { alarm: this.alarm } });
    },
  },

};
</script>

<style scoped>
</style>
