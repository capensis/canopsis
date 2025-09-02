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
          <span>{{ preparedText }}</span>
        </v-layout>
      </v-card-text>
    </v-card>
    <v-messages
      v-if="!isTableColorReadable"
      :value="[$t('theme.errors.notReadable')]"
      class="mt-2"
      color="error"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { isReadableColor, isDarkColor } from '@/helpers/color';

import { useI18n } from '@/hooks/i18n';

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
    text: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const { t } = useI18n();

    const isDarkBackground = computed(() => isDarkColor(props.backgroundColor));
    const isTableColorReadable = computed(() => isReadableColor(props.backgroundColor, props.color));
    const preparedText = computed(() => props.text || t('theme.exampleText'));

    return {
      isDarkBackground,
      isTableColorReadable,
      preparedText,
    };
  },
};
</script>
