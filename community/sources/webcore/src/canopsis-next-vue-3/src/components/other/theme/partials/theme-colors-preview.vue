<template>
  <c-layout column>
    <v-card :color="backgroundColor" :style="cardStyle">
      <v-card-text>
        <c-layout justify="space-between" align="center">
          <span>{{ preparedText }}</span>
        </c-layout>
      </v-card-text>
    </v-card>
    <v-messages
      v-if="!isTableColorReadable"
      :messages="[t('theme.errors.notReadable')]"
      class="mt-2"
      color="error"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { isReadableColor } from '@/helpers/color';

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

    const isTableColorReadable = computed(() =>
      isReadableColor(props.backgroundColor, props.color),
    );
    const preparedText = computed(() => props.text || t('theme.exampleText'));
    const cardStyle = computed(() => ({
      color: props.color,
      fontSize: props.fontSize ? `${props.fontSize}px` : undefined,
    }));

    return {
      t,
      isTableColorReadable,
      preparedText,
      cardStyle,
    };
  },
};
</script>
