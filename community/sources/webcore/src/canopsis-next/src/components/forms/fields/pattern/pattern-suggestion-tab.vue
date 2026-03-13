<template>
  <v-tab
    :class="{ 'pattern-suggestions__tab--difference': hasDifference }"
    class="pattern-suggestions__tab"
  >
    <v-layout
      class="fill-height"
      align-center
    >
      <v-layout
        class="pattern-suggestions__tab__number"
        align-center
      >
        {{ number }}
      </v-layout>
      <v-expand-x-transition>
        <div v-show="active" class="pattern-suggestions__tab__found">
          <v-layout class="gap-2 px-2" align-center>
            <span v-html="$t('pattern.foundEntities', { count: entitiesCount })" class="font-weight-regular" />
            <v-btn
              v-if="hasDifference"
              plain
              small
              @click="showEntitiesComparisonModal"
            >
              {{ $t('pattern.seeRecordsComparison') }}
            </v-btn>
            <strong v-else class="primary--text">{{ $t('pattern.sameEntities') }}</strong>
          </v-layout>
        </div>
      </v-expand-x-transition>
    </v-layout>
  </v-tab>
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    number: {
      type: Number,
      required: true,
    },
    active: {
      type: Boolean,
      default: false,
    },
    suggestion: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const entitiesCount = computed(() => props.suggestion.found_entities);
    const hasDifference = computed(() => !!props.suggestion.difference);

    const showEntitiesComparisonModal = () => emit('show:entities-comparison', props.suggestion);

    return {
      entitiesCount,
      hasDifference,

      showEntitiesComparisonModal,
    };
  },
};
</script>

<style lang="scss">
.pattern-suggestions__tab {
  position: relative;
  min-width: 0;
  max-width: 465px;
  padding: 0;
  overflow: hidden;
  border-bottom: 2px solid var(--v-primary-base) !important;
  opacity: .6;

  .pattern-suggestions--active-difference & {
    border-bottom-color: var(--v-info-base) !important;
  }

  &, &:before {
    border-top-left-radius: 10px;
    border-top-right-radius: 10px;
  }

  &__found {
    overflow-wrap: normal;
    white-space: nowrap;
  }

  &__number {
    height: 100%;
    padding: 0 16px;
    background-color: var(--v-success-base);
    transition: background-color 0.3s ease, opacity 0.3s ease;
    color: var(--v-text-dark-primary, #FFFFFF);

    .pattern-suggestions__tab--difference & {
      background-color: var(--v-info-base);
    }
  }

  &__found {
    text-transform: lowercase;

    .theme--light & {
      color: var(--v-text-light-primary, rgba(0, 0, 0, 0.87));
    }

    .theme--dark & {
      color: var(--v-text-dark-primary, #FFFFFF);
    }
  }

  &.v-tab--active {
    border: 2px solid var(--v-primary-base);
    border-bottom: 0 !important;
    opacity: 1;

    .pattern-suggestions__tab__number {
      background-color: var(--v-primary-base);
      color: var(--v-text-dark-primary, #FFFFFF);
      margin-left: -2px;
    }
  }

  &--difference {
    color: var(--v-info-base) !important;
    border-color: var(--v-info-base) !important;

    &.v-tab--active .pattern-suggestions__tab__number {
      background-color: var(--v-info-base) !important;
      color: var(--v-text-dark-primary, #FFFFFF);
      margin-left: -2px;
    }
  }
}
</style>
