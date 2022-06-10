<template lang="pug">
  g
    slot
      rect(
        :x="square.x",
        :y="square.y",
        :width="square.size",
        :height="square.size",
        fill="transparent",
        cursor="move",
        pointer-events="all",
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
import RectSelection from '../common/rect-selection.vue';

import { DIRECTIONS } from '../../constants';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectSelection },
  props: {
    square: {
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
    rect() {
      return {
        x: this.square.x,
        y: this.square.y,
        width: this.square.size,
        height: this.square.size,
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

    onResize({ x: cursorX, y: cursorY }) {
      const newSquare = {
        x: this.square.x,
        y: this.square.y,
        size: this.square.size,
      };

      switch (this.direction) {
        case DIRECTIONS.west:
        case DIRECTIONS.southWest: {
          const newSize = this.square.size + this.square.x - cursorX;

          if (newSize >= 0) {
            newSquare.size = newSize;
            newSquare.x = cursorX;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newSquare.x += this.square.size;
            newSquare.y -= absoluteNewSize;
            newSquare.size = absoluteNewSize;

            this.direction = DIRECTIONS.northEast;
          }

          break;
        }
        case DIRECTIONS.east:
        case DIRECTIONS.northEast: {
          const newSize = cursorX - this.square.x;

          if (newSize >= 0) {
            newSquare.y += this.square.size - newSize;
            newSquare.size = newSize;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newSquare.y += this.square.size;
            newSquare.x -= absoluteNewSize;
            newSquare.size = absoluteNewSize;

            this.direction = DIRECTIONS.southWest;
          }

          break;
        }
        case DIRECTIONS.south:
        case DIRECTIONS.southEast: {
          const newSize = cursorY - this.square.y;

          if (newSize >= 0) {
            newSquare.size = newSize;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newSquare.size = absoluteNewSize;
            newSquare.x -= absoluteNewSize;
            newSquare.y -= absoluteNewSize;

            this.direction = DIRECTIONS.northWest;
          }

          break;
        }
        case DIRECTIONS.north:
        case DIRECTIONS.northWest: {
          const newSize = this.square.y + this.square.size - cursorY;

          if (newSize >= 0) {
            const diffBetweenSizes = this.square.size - newSize;

            newSquare.size = newSize;
            newSquare.y += diffBetweenSizes;
            newSquare.x += diffBetweenSizes;
          } else {
            const absoluteNewSize = Math.abs(newSize);

            newSquare.y += this.square.size;
            newSquare.x += this.square.size;
            newSquare.size = absoluteNewSize;

            this.direction = DIRECTIONS.southEast;
          }

          break;
        }
      }

      this.$emit('resize', newSquare, this.direction);
    },
  },
};
</script>
