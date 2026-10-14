<template>
  <v-menu
    :model-value="value"
    :target="[positionX, positionY]"
    :close-on-content-click="false"
    @update:model-value="updateShown"
  >
    <v-list dense>
      <v-list-item v-for="item in items" :key="item.event" @click="applyEvent(item.event)">
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
    positionX: {
      type: Number,
      required: true,
    },
    positionY: {
      type: Number,
      required: true,
    },
    editing: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['add:point', 'edit:point', 'remove:point', 'close'],
  computed: {
    items() {
      if (this.editing) {
        return [
          {
            text: this.$t('map.editPoint'),
            event: 'edit:point',
          },
          {
            text: this.$t('map.removePoint'),
            event: 'remove:point',
          },
        ];
      }

      return [
        {
          text: this.$t('map.addPoint'),
          event: 'add:point',
        },
      ];
    },
  },
  methods: {
    updateShown(value) {
      if (!value) {
        this.$emit('close');
      }
    },

    applyEvent(event) {
      this.$emit(event);
    },
  },
};
</script>
