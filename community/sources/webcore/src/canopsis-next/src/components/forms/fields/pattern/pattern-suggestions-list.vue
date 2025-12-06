<template>
  <div
    v-if="suggestions.length"
    :class="{ 'pattern-suggestions--active-difference': hasDifference }"
    class="pattern-suggestions"
  >
    <v-layout
      align-center
      wrap
    >
      <v-tabs
        v-model="activeSuggestionIndex"
        class="pattern-suggestions__tabs"
        hide-slider
      >
        <pattern-suggestion-tab
          v-for="(suggestion, index) in suggestions"
          :key="index"
          :number="index + 1"
          :active="index === activeSuggestionIndex"
          :suggestion="suggestion"
          @show:entities-comparison="showEntitiesComparisonModal"
        />
      </v-tabs>
    </v-layout>

    <v-tabs-items v-model="activeSuggestionIndex">
      <v-tab-item
        v-for="(suggestion, index) in suggestions"
        :key="index"
        :value="index"
      >
        <pattern-suggestion-content
          :suggestion="suggestion"
          :entity-attributes="entityAttributes"
          :optimized-fields-regexps="optimizedFieldsRegexps"
          @apply="apply(index)"
        />
      </v-tab-item>
    </v-tabs-items>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

import PatternSuggestionTab from './pattern-suggestion-tab.vue';
import PatternSuggestionContent from './pattern-suggestion-content.vue';

export default {
  components: {
    PatternSuggestionTab,
    PatternSuggestionContent,
  },
  props: {
    suggestions: {
      type: Array,
      default: () => [],
    },
    entityAttributes: {
      type: Array,
      default: () => [],
    },
    optimizedFieldsRegexps: {
      type: Array,
      default: () => [],
    },
  },
  setup(props, { emit }) {
    const activeSuggestionIndex = ref(0);

    const hasDifference = computed(() => props.suggestions[activeSuggestionIndex.value]?.difference);

    const apply = index => emit('apply', index);
    const showEntitiesComparisonModal = suggestion => emit('show:entities-comparison', suggestion);

    return {
      activeSuggestionIndex,
      hasDifference,
      apply,
      showEntitiesComparisonModal,
    };
  },
};
</script>

<style lang="scss">
.pattern-suggestions {
  --pattern-border-primary: var(--v-primary-base);
  --pattern-border-info: var(--v-info-base);

  .theme--dark & {
    --pattern-border-primary: var(--v-primary-base);
    --pattern-border-info: var(--v-info-base);
  }

  .theme--light & {
    &__content, &__tab {
      background-color: var(--v-background-darken1);
    }
  }

  .theme--dark & {
    &__content, &__tab {
      background-color: var(--v-background-lighten1);
    }
  }

  &__pattern {
    margin-top: 16px;
  }

  &__tabs  {
    position: relative;
    top: 2px;
    z-index: 2;

    .v-tabs-bar__content {
      column-gap: 12px;
    }
  }
  }
</style>
