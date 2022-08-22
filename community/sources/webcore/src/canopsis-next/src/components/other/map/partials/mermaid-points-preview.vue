<template lang="pug">
  div
    mermaid-point-marker(
      v-for="point in points",
      :key="point._id",
      :x="point.x",
      :y="point.y",
      :entity="point.entity",
      :size="markerSize",
      :color-indicator="colorIndicator",
      :pbehavior-enabled="pbehaviorEnabled",
      ref="points",
      @click="openPopup(point, $event)"
    )
    v-menu(
      v-if="activePoint",
      :value="true",
      :position-x="positionX",
      :position-y="positionY",
      :close-on-content-click="false",
      ignore-click-outside,
      offset-overflow,
      offset-x,
      absolute,
      top
    )
      point-popup(
        v-click-outside="clickOutsideDirective",
        :point="activePoint",
        :template="popupTemplate",
        :color-indicator="colorIndicator",
        :actions="popupActions",
        @show:alarms="showAlarms",
        @show:map="showLinkedMap",
        @close="closePopup"
      )
</template>

<script>
import { entitiesServiceEntityMixin } from '@/mixins/entities/service-entity';

import MermaidPointMarker from './mermaid-point-marker.vue';
import PointPopup from './point-popup.vue';

export default {
  components: { MermaidPointMarker, PointPopup },
  mixins: [entitiesServiceEntityMixin],
  props: {
    points: {
      type: Array,
      required: true,
    },
    markerSize: {
      type: Number,
      default: 24,
    },
    popupTemplate: {
      type: String,
      required: false,
    },
    popupActions: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      positionX: 0,
      positionY: 0,
      activePoint: undefined,
    };
  },
  computed: {
    clickOutsideDirective() {
      return {
        handler: this.closePopup,
        include: () => this.$refs.points.map(({ $el }) => $el),
        closeConditional: () => true,
      };
    },
  },
  methods: {
    openPopup(point, event) {
      const { top, left, width } = event.target.getBoundingClientRect();

      this.positionY = top;
      this.positionX = left + width / 2;
      this.activePoint = point;
    },

    closePopup() {
      this.activePoint = undefined;
    },

    showLinkedMap() {
      this.$emit('show:map', this.activePoint.map);
      this.closePopup();
    },

    showAlarms() {
      this.$emit('show:alarms', this.activePoint);
      this.closePopup();
    },
  },
};
</script>
