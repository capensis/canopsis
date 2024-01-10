<template>
  <v-menu
    v-model="shown"
    :position-x="pageX"
    :position-y="pageY"
    :close-on-content-click="false"
    :disabled="disabled"
    ignore-click-upper-outside
    offset-overflow
    offset-x
    absolute
  >
    <v-list dense>
      <v-list-item
        v-for="action in availableActions"
        :key="action.text"
        @click="applyAction(action)"
      >
        <v-list-item-content>
          <v-list-item-title>{{ action.text }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-menu>
</template>

<script>
import { findRealParent } from 'vue2-leaflet';

export default {
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    markerItems: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      mapObject: undefined,
      latlng: undefined,
      marker: undefined,
      shown: false,
      pageX: 0,
      pageY: 0,
    };
  },
  computed: {
    availableActions() {
      if (this.marker) {
        return this.markerItems;
      }

      return this.items;
    },
  },
  mounted() {
    this.parentContainer = findRealParent(this.$parent);

    this.mapObject = this.parentContainer.mapObject;

    this.mapObject.on('contextmenu', this.open);
  },

  beforeDestroy() {
    this.mapObject.off('contextmenu', this.open);
  },

  methods: {
    open({ latlng, originalEvent, marker }) {
      if (this.disabled) {
        return;
      }

      if (this.shown) {
        this.shown = false;
        return;
      }

      this.latlng = latlng;
      this.marker = marker;
      this.pageX = originalEvent.pageX;
      this.pageY = originalEvent.pageY - window.scrollY;
      this.shown = true;
    },

    close() {
      this.shown = false;
    },

    applyAction({ action }) {
      action({
        latlng: this.latlng,
        marker: this.marker,
        pageX: this.pageX,
        pageY: this.pageY,
      });
    },
  },
};
</script>
