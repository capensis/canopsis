<template>
  <label>
    <c-layout
      align="center"
      class="images-viewer bg-white"
    >
      <v-window v-model="activeImageIndex">
        <v-window-item
          v-for="image in images"
          :key="image[itemKey]"
          :reverse-transition="transition"
          :transition="transition"
          class="images-viewer__img-wrapper"
        >
          <div
            v-if="image.name"
            class="images-viewer__img-title text-white ml-1"
          >
            {{ image.name }}
          </div>
          <img
            :src="image.src"
            :alt="image.name"
            class="images-viewer__img"
          >
        </v-window-item>
      </v-window>

      <c-layout
        v-if="images.length"
        justify="space-between"
        class="images-viewer__actions"
      >
        <v-btn
          color="white"
          variant="text"
          icon
          @click="prevImage"
        >
          <c-icon :size="32">
            chevron_left
          </c-icon>
        </v-btn>

        <v-btn
          color="white"
          variant="text"
          icon
          @click="nextImage"
        >
          <c-icon :size="32">
            chevron_right
          </c-icon>
        </v-btn>
      </c-layout>

      <c-layout
        v-if="images.length"
        justify="space-between"
        class="images-viewer__image-count"
      >
        <span class="text-white pl-1">
          {{ activeImageIndex + 1 }}/{{ images.length }}
        </span>
      </c-layout>
    </c-layout>
  </label>
</template>

<script>
import { throttle } from 'lodash';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

/**
 * Example:
 *   c-images-viewer(
 *     :images="[{
 *       name: 'Image 1',
 *       src: 'https://src.com/image1'
 *     }, {
 *       name: 'Image 2',
 *       src: 'https://src.com/image2'
 *     }]",
 *     active="https://src.com/image1"
 *   )
 */
export default {
  props: {
    images: {
      type: Array,
      default: () => [],
    },
    active: {
      type: String,
      required: false,
    },
    itemKey: {
      type: String,
      default: 'src',
    },
  },

  data() {
    return {
      directionToNext: true,
      activeImageIndex: this.active
        ? this.images.findIndex(({ src }) => src === this.active)
        : 0,
    };
  },

  computed: {
    transition() {
      return this.directionToNext ? 'v-window-x-transition' : 'v-window-x-reverse-transition';
    },
  },

  mounted() {
    document.addEventListener('keydown', this.keyDownListener);
  },

  beforeUnmount() {
    document.removeEventListener('keydown', this.keyDownListener);
  },

  methods: {
    keyDownListener(event) {
      switch (event.key) {
        case 'ArrowRight':
          this.nextImage();
          break;
        case 'ArrowLeft':
          this.prevImage();
          break;
      }
    },

    prevImage() {
      this.setImageIndex(this.activeImageIndex - 1, false);
    },

    nextImage() {
      this.setImageIndex(this.activeImageIndex + 1, true);
    },

    setImageIndex: throttle(function setImageIndex(nextIndex, directionToNext) {
      this.directionToNext = directionToNext;

      this.activeImageIndex = (this.images.length + nextIndex) % this.images.length;
    }, VUETIFY_ANIMATION_DELAY * 2),
  },
};
</script>

<style lang="scss">
.images-viewer {
  align-items: center;

  .v-window__container {
    display: flex;
    align-items: center;
  }

  &__actions {
    position: absolute;
    left: 0;
    right: 0;
  }

  &__img-wrapper {
    top: unset !important;
  }

  &__img {
    display: block;
    max-height: 90vh;
    max-width: 90vw;
  }

  &__img-title {
    position: absolute;
    top: 0;
    right: 0;
    left: 0;
  }

  &__image-count {
    pointer-events: none;
    position: absolute;
    bottom: 0;
    left: 0;
  }
}
</style>
