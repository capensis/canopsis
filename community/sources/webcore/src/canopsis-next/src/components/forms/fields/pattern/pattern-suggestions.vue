<template>
  <v-layout
    class="c-pattern-suggestions gap-4"
    column
  >
    <pattern-suggestions-header @reject-all="rejectAll" />

    <pattern-suggestions-list
      :suggestions="suggestions"
      @apply="applySuggestion"
    />
  </v-layout>
</template>

<script>
import PatternSuggestionsHeader from './pattern-suggestions-header.vue';
import PatternSuggestionsList from './pattern-suggestions-list.vue';

export default {
  components: {
    PatternSuggestionsHeader,
    PatternSuggestionsList,
  },
  props: {
    suggestions: {
      type: Array,
      default: () => [],
    },
    type: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const rejectAll = () => emit('reject-all');

    const applySuggestion = (suggestion, index) => {
      if (suggestion) {
        emit('apply-suggestion', suggestion, index);
      }
    };

    return {
      rejectAll,
      applySuggestion,
    };
  },
};
</script>

<style lang="scss">
.c-pattern-suggestions {
  &__content {
    border: 2px solid #4caf50;
    border-radius: 4px;
    padding: 16px;
    background-color: #f5f5f5;
  }

  &__pattern {
    margin-top: 16px;
  }
}
</style>
