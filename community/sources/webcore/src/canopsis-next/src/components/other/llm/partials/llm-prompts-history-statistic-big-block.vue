<template>
  <v-card
    :class="blockClasses"
    class="llm-prompts-history-statistic-big-block"
    elevation="0"
    rounded="lg"
  >
    <v-row class="pa-4 ma-0" align="center">
      <v-col cols="auto">
        <div class="text-title-large mb-3">
          {{ title }}
        </div>
        <div :class="valueClasses">
          {{ value }}
        </div>
        <div class="text-center text-body-2">
          {{ period }}
        </div>
      </v-col>

      <v-col v-if="$slots.default" cols="auto">
        <slot />
      </v-col>
    </v-row>
  </v-card>
</template>

<script>
import { computed } from 'vue';
import { useTheme } from 'vuetify';

export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      required: true,
    },
    period: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'info',
    },
  },

  setup(props) {
    const theme = useTheme();

    const blockClasses = computed(() => ({
      'llm-prompts-history-statistic-big-block--dark': theme.current.value.dark,
      'llm-prompts-history-statistic-big-block--light': !theme.current.value.dark,
    }));

    const valueClasses = computed(() => [
      'text-h6',
      'font-weight-bold',
      'text-no-wrap',
      'text-center',
      'mb-2',
      `llm-prompts-history-statistic-big-block__value--${props.color}`,
    ]);

    return {
      blockClasses,
      valueClasses,
    };
  },
};
</script>

<style lang="scss" scoped>
.llm-prompts-history-statistic-big-block {
  width: fit-content;
  border: 1px solid;

  &--light {
    background-color: #f9f9f9;
    border-color: #e0e0e0;
  }

  &--dark {
    background-color: var(--v-background-lighten1);
    border-color: var(--v-background-lighten2);
  }

  &__value--info {
    color: var(--v-info-base);
  }

  &__value--warning {
    color: var(--v-warning-base);
  }
}
</style>
