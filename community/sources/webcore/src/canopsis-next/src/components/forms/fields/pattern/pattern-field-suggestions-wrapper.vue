<template>
  <v-layout
    class="pattern-field-suggestions-wrapper gap-3"
    column
  >
    <div
      :class="{ 'pattern-field-suggestions-wrapper__content--has-suggestions': !!suggestions.length }"
      class="pattern-field-suggestions-wrapper__content"
    >
      <v-layout
        v-if="suggestions.length"
        class="mb-3"
        align-center
      >
        <v-chip class="pattern-field-suggestions-wrapper__current" small>
          <strong>{{ $t('common.current') }}</strong>
        </v-chip>
      </v-layout>
      <slot />
    </div>
    <pattern-suggestions
      v-if="suggestions.length"
      :suggestions="suggestions"
      :patterns="patterns"
      :entity-attributes="entityAttributes"
      :optimized-fields-regexps="optimizedFieldsRegexps"
      @apply:suggestion="applySuggestion"
      @reject:all="rejectAllSuggestions"
      @show:entities-comparison="showEntitiesComparisonModal"
    />
  </v-layout>
</template>

<script>
import PatternSuggestions from './pattern-suggestions.vue';

export default {
  components: {
    PatternSuggestions,
  },
  props: {
    suggestions: {
      type: Array,
      default: () => [],
    },
    patterns: {
      type: Object,
      required: true,
    },
    entityAttributes: {
      type: Array,
      default: () => [],
    },
    optimizedFieldsRegexps: {
      type: Array,
      default: () => [],
    },
    entitiesCount: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const applySuggestion = index => emit('apply:suggestion', index);
    const rejectAllSuggestions = () => emit('reject:all');
    const showEntitiesComparisonModal = suggestion => emit('show:entities-comparison', suggestion);

    return {
      applySuggestion,
      rejectAllSuggestions,
      showEntitiesComparisonModal,
    };
  },
};
</script>

<style lang="scss">
:root {
  --pattern-suggestions-border-color: #808080;
}

.pattern-field-suggestions-wrapper {
  &__content {
    padding: 12px;

    &--has-suggestions {
      gap: 8px;
      border: 2px solid var(--pattern-suggestions-border-color);
      border-radius: 10px;
    }
  }

  &__current {
    background-color: var(--pattern-suggestions-border-color) !important;
    color: var(--v-text-dark-primary, #FFFFFF) !important;
  }
}
</style>
