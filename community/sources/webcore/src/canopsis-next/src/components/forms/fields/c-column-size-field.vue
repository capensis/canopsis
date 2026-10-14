<template>
  <c-radio-group
    :model-value="modelValue"
    :class="{ 'c-column-size-field--mobile': mobile }"
    :name="name"
    class="c-column-size-field"
    color="primary"
    hide-details
    inline
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="item in availableItems"
      :key="item.src"
      :value="item.value"
      color="primary"
    >
      <template #label>
        <v-img :src="item.src" class="my-2" />
      </template>
    </v-radio>
  </c-radio-group>
</template>

<script>
import oneColumnDesktopSrc from '@/assets/images/column-desktop-1.svg';
import twoColumnDesktopSrc from '@/assets/images/column-desktop-2.svg';
import threeColumnDesktopSrc from '@/assets/images/column-desktop-3.svg';
import fourColumnDesktopSrc from '@/assets/images/column-desktop-4.svg';
import sixColumnDesktopSrc from '@/assets/images/column-desktop-6.svg';
import twelveColumnDesktopSrc from '@/assets/images/column-desktop-12.svg';
import oneColumnMobileSrc from '@/assets/images/column-mobile-1.svg';
import twoColumnMobileSrc from '@/assets/images/column-mobile-2.svg';
import oneColumnTabletSrc from '@/assets/images/column-tablet-1.svg';
import twoColumnTabletSrc from '@/assets/images/column-tablet-2.svg';
import threeColumnTabletSrc from '@/assets/images/column-tablet-3.svg';
import fourColumnTabletSrc from '@/assets/images/column-tablet-4.svg';

export default {
  props: {
    modelValue: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'size',
    },
    mobile: {
      type: Boolean,
      default: false,
    },
    tablet: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    availableItems() {
      if (this.mobile) {
        return [
          { value: 1, src: oneColumnMobileSrc },
          { value: 2, src: twoColumnMobileSrc },
        ];
      }

      if (this.tablet) {
        return [
          { value: 1, src: oneColumnTabletSrc },
          { value: 2, src: twoColumnTabletSrc },
          { value: 3, src: threeColumnTabletSrc },
          { value: 4, src: fourColumnTabletSrc },
        ];
      }

      return [
        { value: 1, src: oneColumnDesktopSrc },
        { value: 2, src: twoColumnDesktopSrc },
        { value: 3, src: threeColumnDesktopSrc },
        { value: 4, src: fourColumnDesktopSrc },
        { value: 6, src: sixColumnDesktopSrc },
        { value: 12, src: twelveColumnDesktopSrc },
      ];
    },
  },
};
</script>

<style lang="scss">
.c-column-size-field {
  .v-input__control {
    width: 100%;
  }

  .v-selection-control {
    flex-direction: column;
    width: 50%;
    margin-right: 0 !important;
  }

  .v-label {
    width: 100%;
    max-width: 170px;
  }

  &--mobile {
    .v-label {
      max-width: 100px;
    }
  }
}
</style>
