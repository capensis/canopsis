<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    rect(
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
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
      cursor="move",
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
    rect: {
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
      const directionArray = this.direction.split('');
      const newRect = {
        x: this.rect.x,
        y: this.rect.y,
        width: this.rect.width,
        height: this.rect.height,
      };

      directionArray.forEach((direction, index) => {
        switch (direction) {
          case DIRECTIONS.south: {
            const newHeight = cursorY - newRect.y;

            if (newHeight > 0) {
              newRect.height = newHeight;
            } else {
              newRect.height = Math.abs(newHeight);
              newRect.y -= newRect.height;

              directionArray[index] = DIRECTIONS.north;
            }

            break;
          }
          case DIRECTIONS.north: {
            const newHeight = newRect.height + newRect.y - cursorY;

            if (newHeight > 0) {
              newRect.height = newHeight;
              newRect.y = cursorY;
            } else {
              newRect.height = Math.abs(newHeight);
              newRect.y = cursorY - newRect.height;

              directionArray[index] = DIRECTIONS.south;
            }

            break;
          }
          case DIRECTIONS.east: {
            const newWidth = cursorX - newRect.x;

            if (newWidth > 0) {
              newRect.width = newWidth;
            } else {
              newRect.width = Math.abs(newWidth);
              newRect.x = cursorX;

              directionArray[index] = DIRECTIONS.west;
            }

            break;
          }
          case DIRECTIONS.west: {
            const newWidth = newRect.width + newRect.x - cursorX;

            if (newWidth > 0) {
              newRect.width = newWidth;
              newRect.x = cursorX;
            } else {
              newRect.width = Math.abs(newWidth);
              newRect.x = cursorX - newRect.width;

              directionArray[index] = DIRECTIONS.east;
            }

            break;
          }
        }
      });

      this.direction = directionArray.join('');

      this.$emit('resize', newRect, this.direction);
    },
  },
};
</script>
