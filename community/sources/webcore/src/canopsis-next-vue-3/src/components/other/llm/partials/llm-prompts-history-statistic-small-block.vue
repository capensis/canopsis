<template>
  <v-sheet
    :class="blockClasses"
    class="llm-prompts-history-statistic-small-block pa-3"
    elevation="0"
    rounded="lg"
  >
    <div class="text-body-1 font-weight-medium mb-3">
      {{ title }}
    </div>
    <v-row align="center" justify="space-between" no-gutters>
      <v-col cols="auto">
        <span class="llm-prompts-history-statistic-small-block__value text-title-large font-weight-medium text-no-wrap">
          {{ value }}
        </span>
      </v-col>
      <v-col v-if="percent || percent === 0" class="text-right" cols="auto">
        <span class="llm-prompts-history-statistic-small-block__percent text-h6 font-weight-medium text-no-wrap">
          {{ percent }}%
        </span>
      </v-col>
    </v-row>
  </v-sheet>
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
    percent: {
      type: [String, Number],
      default: '',
    },
  },

  setup() {
    const theme = useTheme();

    const blockClasses = computed(() => ({
      'llm-prompts-history-statistic-small-block--dark': theme.current.value.dark,
      'llm-prompts-history-statistic-small-block--light': !theme.current.value.dark,
    }));

    return {
      blockClasses,
    };
  },
};
</script>

<style lang="scss" scoped>
.llm-prompts-history-statistic-small-block {
  flex: 0 0 200px;
  width: 200px;
  border: 1px solid;

  &--light {
    background-color: #d5eaf9;
    border-color: #e0e0e0;

    .llm-prompts-history-statistic-small-block__percent {
      color: var(--v-state-pause-base);
    }
  }

  &--dark {
    background-color: var(--v-background-lighten2);
    border-color: var(--v-background-lighten3);

    .llm-prompts-history-statistic-small-block__percent {
      color: var(--v-state-pause-lighten2);
    }
  }
}

.llm-prompts-history-statistic-small-block__value {
  color: var(--v-info-base);
}
</style>
