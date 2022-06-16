<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    circle(
      :cx="centerX",
      :cy="centerY",
      :r="radius",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :rect="rect",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      @start:resize="startResize"
    )
</template>

<script>
import { DIRECTIONS } from '../../constants';

import RectSelection from '../common/rect-selection.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectSelection },
  props: {
    circle: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    padding: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'blue',
    },
    pointerEvents: {
      type: String,
      default: 'all',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
  },
  data() {
    return {
      direction: undefined,
    };
  },
  computed: {
    radius() {
      return this.circle.diameter / 2;
    },

    centerX() {
      return this.circle.x + this.radius;
    },

    centerY() {
      return this.circle.y + this.radius;
    },

    rect() {
      return {
        x: this.circle.x,
        y: this.circle.y,
        width: this.circle.diameter,
        height: this.circle.diameter,
      };
    },
  },
  methods: {
    startResize(direction) {
      this.$mouseMove.register(this.onResize);
      this.$mouseUp.register(this.finishResize);
      this.direction = direction;
    },

    finishResize() {
      this.$mouseMove.unregister(this.onResize);
      this.$mouseUp.unregister(this.finishResize);
    },

    /**
     * TODO: Should be moved in utils and reused with square
     */
    onResize({ x: cursorX, y: cursorY }) {
      const newCircle = {
        x: this.circle.x,
        y: this.circle.y,
        size: this.circle.diameter,
      };

      switch (this.direction) {
        case DIRECTIONS.west:
        case DIRECTIONS.southWest: {
          const newSize = this.circle.size + this.circle.x - cursorX;

          if (newSize >= 0) {
            newCircle.size = newSize;
            newCircle.x = cursorX;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newCircle.x += this.circle.diameter;
            newCircle.y -= absoluteNewSize;
            newCircle.size = absoluteNewSize;

            this.direction = DIRECTIONS.northEast;
          }

          break;
        }
        case DIRECTIONS.east:
        case DIRECTIONS.northEast: {
          const newSize = cursorX - this.circle.x;

          if (newSize >= 0) {
            newCircle.y += this.circle.diameter - newSize;
            newCircle.size = newSize;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newCircle.y += this.circle.diameter;
            newCircle.x -= absoluteNewSize;
            newCircle.size = absoluteNewSize;

            this.direction = DIRECTIONS.southWest;
          }

          break;
        }
        case DIRECTIONS.south:
        case DIRECTIONS.southEast: {
          const newSize = cursorY - this.circle.y;

          if (newSize >= 0) {
            newCircle.size = newSize;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newCircle.size = absoluteNewSize;
            newCircle.x -= absoluteNewSize;
            newCircle.y -= absoluteNewSize;

            this.direction = DIRECTIONS.northWest;
          }

          break;
        }
        case DIRECTIONS.north:
        case DIRECTIONS.northWest: {
          const newSize = this.circle.y + this.circle.diameter - cursorY;

          if (newSize >= 0) {
            const diffBetweenSizes = this.circle.diameter - newSize;

            newCircle.size = newSize;
            newCircle.y += diffBetweenSizes;
            newCircle.x += diffBetweenSizes;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newCircle.y += this.circle.diameter;
            newCircle.x += this.circle.diameter;
            newCircle.size = absoluteNewSize;

            this.direction = DIRECTIONS.southEast;
          }

          break;
        }
      }

      this.$emit('resize', newCircle, this.direction);
    },
  },
};
</script>
