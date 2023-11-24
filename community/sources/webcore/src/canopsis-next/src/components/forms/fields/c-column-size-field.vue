<template lang="pug">
  v-radio-group.c-column-size-field(
    v-field="value",
    :class="{ 'c-column-size-field--mobile': mobile }",
    :name="name",
    color="primary",
    hide-details,
    mandatory,
    row
  )
    v-radio.ma-0(
      v-for="item in availableItems",
      :key="item.src",
      :value="item.value",
      color="primary"
    )
      template(#label="")
        v-img.my-2(:src="item.src")
</template>

<script>
import { computed } from 'vue';

import oneColumnMobileSrc from '@/assets/images/column-mobile-1.svg';
import twoColumnMobileSrc from '@/assets/images/column-mobile-2.svg';
import oneColumnTabletSrc from '@/assets/images/column-tablet-1.svg';
import twoColumnTabletSrc from '@/assets/images/column-tablet-2.svg';
import threeColumnTabletSrc from '@/assets/images/column-tablet-3.svg';
import fourColumnTabletSrc from '@/assets/images/column-tablet-4.svg';
import oneColumnDesktopSrc from '@/assets/images/column-desktop-1.svg';
import twoColumnDesktopSrc from '@/assets/images/column-desktop-2.svg';
import threeColumnDesktopSrc from '@/assets/images/column-desktop-3.svg';
import fourColumnDesktopSrc from '@/assets/images/column-desktop-4.svg';
import sixColumnDesktopSrc from '@/assets/images/column-desktop-6.svg';
import twelveColumnDesktopSrc from '@/assets/images/column-desktop-12.svg';

export default {
  props: {
    value: {
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
  setup(props) {
    const mobileItems = computed(() => [
      { value: 1, src: oneColumnMobileSrc },
      { value: 2, src: twoColumnMobileSrc },
    ]);

    const tabletItems = computed(() => [
      { value: 1, src: oneColumnTabletSrc },
      { value: 2, src: twoColumnTabletSrc },
      { value: 3, src: threeColumnTabletSrc },
      { value: 4, src: fourColumnTabletSrc },
    ]);

    const desktopItems = computed(() => [
      { value: 1, src: oneColumnDesktopSrc },
      { value: 2, src: twoColumnDesktopSrc },
      { value: 3, src: threeColumnDesktopSrc },
      { value: 4, src: fourColumnDesktopSrc },
      { value: 6, src: sixColumnDesktopSrc },
      { value: 12, src: twelveColumnDesktopSrc },
    ]);

    const availableItems = computed(() => {
      if (props.mobile) {
        return mobileItems.value;
      }

      if (props.tablet) {
        return tabletItems.value;
      }

      return desktopItems.value;
    });

    return {
      availableItems,
    };
  },
};
</script>

<style lang="scss">
.c-column-size-field {
  .v-input__control {
    width: 100%;
  }

  .v-radio {
    flex-direction: column;
    width: 50%;
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
