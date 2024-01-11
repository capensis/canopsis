<template>
  <v-layout column>
    <v-card
      :color="backgroundColor"
      :style="{ color, fontSize: `${fontSize}px` }"
      :dark="isDarkBackground"
    >
      <v-card-text>
        <v-layout
          justify-space-between
          align-center
        >
          <span>{{ $t('theme.exampleText') }}</span>
          <v-icon :color="color">
            help
          </v-icon>
        </v-layout>
      </v-card-text>
    </v-card>
    <v-messages
      class="mt-2"
      v-if="!isTableColorReadable"
      :value="[$t('theme.errors.notReadable')]"
      color="error"
    />
  </v-layout>
</template>

<script>
import { isReadableColor, isDarkColor } from '@/helpers/color';

export default {
  props: {
    backgroundColor: {
      type: String,
      required: true,
    },
    color: {
      type: String,
      required: true,
    },
    fontSize: {
      type: Number,
      required: false,
    },
  },
  computed: {
    isDarkBackground() {
      return isDarkColor(this.backgroundColor);
    },

    isTableColorReadable() {
      return isReadableColor(this.backgroundColor, this.color);
    },
  },
};
</script>
