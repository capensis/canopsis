<template>
  <div @contextmenu.prevent="handleClickContextMenu">
    <slot />
    <v-menu
      v-model="shown"
      :position-x="position.x"
      :position-y="position.y"
      :close-on-content-click="false"
      ignore-click-upper-outside
      offset-overflow
      offset-x
      absolute
    >
      <slot
        :position="position"
        :data="data"
        name="menu"
      />
    </v-menu>
  </div>
</template>

<script>
export default {
  provide() {
    return {
      $contextmenu: {
        open: this.openContextmenu,
      },
    };
  },
  data() {
    return {
      shown: false,
      position: {
        x: 0,
        y: 0,
      },
      data: undefined,
    };
  },
  methods: {
    openContextmenu({ x, y, data }) {
      this.position.x = x;
      this.position.y = y;
      this.data = data;

      this.shown = true;
    },

    handleClickContextMenu(event) {
      this.openContextmenu({ x: event.pageX, y: event.pageY });
    },
  },
};
</script>
