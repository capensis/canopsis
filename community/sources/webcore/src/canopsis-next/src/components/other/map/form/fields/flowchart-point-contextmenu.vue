<template>
  <v-menu
    :value="value"
    :position-x="positionX"
    :position-y="positionY"
    :close-on-content-click="false"
    absolute
    @input="$listeners.close"
  >
    <v-list dense>
      <v-list-item
        v-for="item in items"
        :key="item.text"
        @click="item.action"
        @contextmenu.prevent.stop=""
      >
        <v-list-item-content>
          <v-list-item-title>{{ item.text }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
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
