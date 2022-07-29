<template lang="pug">
  div(@contextmenu.prevent="handleClickContextMenu")
    slot
    v-menu(
      v-model="shown",
      :position-x="position.x",
      :position-y="position.y",
      :close-on-content-click="false",
      ignore-click-upper-outside,
      offset-overflow,
      offset-x,
      absolute
    )
      slot(name="menu", :position="position")
</template>

<script>
export default {
  provide() {
    return {
      $openContextmenu: this.openContextmenu,
    };
  },
  data() {
    return {
      shown: false,
      position: {
        x: 0,
        y: 0,
      },
    };
  },
  methods: {
    openContextmenu({ x, y }) {
      this.position.x = x;
      this.position.y = y;

      this.shown = true;
    },

    handleClickContextMenu(event) {
      this.openContextmenu({ x: event.pageX, y: event.pageY });
    },
  },
};
</script>
