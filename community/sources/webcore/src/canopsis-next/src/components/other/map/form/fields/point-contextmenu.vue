<template>
  <v-menu
    :value="value"
    :position-x="positionX"
    :position-y="positionY"
    :close-on-content-click="false"
    absolute
    @input="$emit('close')"
  >
    <v-list dense>
      <v-list-item
        v-for="item in items"
        :key="item.event"
        @click="applyEvent(item.event)"
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
    applyEvent(event) {
      this.$emit(event);
    },
  },
};
</script>
