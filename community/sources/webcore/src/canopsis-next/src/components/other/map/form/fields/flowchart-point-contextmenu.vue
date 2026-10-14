<template>
  <v-menu
    :model-value="value"
    :target="[positionX, positionY]"
    :close-on-content-click="false"
    @update:model-value="updateShown"
  >
    <v-list dense>
      <v-list-item
        v-for="item in items"
        :key="item.text"
        @click="item.action"
        @contextmenu.prevent.stop=""
      >
        <v-list-item-title>{{ item.text }}</v-list-item-title>
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
      default: () => [],
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
  emits: ['close'],
  methods: {
    updateShown(value) {
      if (!value) {
        this.$emit('close');
      }
    },
  },
};
</script>
