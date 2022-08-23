<template lang="pug">
  v-menu(
    :value="value",
    :position-x="positionX",
    :position-y="positionY",
    :close-on-content-click="false",
    absolute,
    @input="$listeners.close"
  )
    v-list(dense)
      v-list-tile(
        v-for="item in items",
        :key="item.text",
        @click="item.action",
        @contextmenu.prevent.stop=""
      )
        v-list-tile-content
          v-list-tile-title {{ item.text }}
</template>

<script>
export default {
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      default: () => ([]),
    },
    positionX: {
      type: Number,
      required: true,
    },
    positionY: {
      type: Number,
      required: true,
    },
  },
};
</script>
