<template>
  <v-layout
    class="fill-height"
    align-center
  >
    <v-layout
      :class="{ 'pattern-suggestions__tabs__number--difference': hasDifference }"
      class="pattern-suggestions__tabs__number"
      align-center
    >
      {{ number }}
    </v-layout>
    <v-expand-x-transition>
      <div v-show="active" class="pattern-suggestions__tabs__found">
        <v-layout class="gap-2 px-2" align-center>
          <span v-html="$t('pattern.foundEntities', { count: entitiesCount })" class="font-weight-regular" />
          <v-btn
            v-if="hasDifference"
            plain
            small
          >
            {{ $t('pattern.seeRecordsComparison') }}
          </v-btn>
          <strong v-else class="primary--text">{{ $t('pattern.sameEntities') }}</strong>
        </v-layout>
      </div>
    </v-expand-x-transition>
  </v-layout>
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
  setup(props) {
    const entitiesCount = computed(() => props.suggestion.found_entities);
    const hasDifference = computed(() => !!props.suggestion.difference);

    return {
      entitiesCount,
      hasDifference,
    };
  },
};
</script>

<style lang="scss">
.pattern-suggestions__tabs {
  &__number {
    height: 100%;
    padding: 0 16px;
    background-color: var(--v-success-lighten1);
    color: var(--v-text-light-primary, rgba(0, 0, 0, 0.87));
    transition: background-color 0.3s ease, color 0.3s ease;

    &--difference {
      background-color: var(--v-info-lighten2);
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

  .v-tab--active &__number {
    background-color: var(--v-success-base);
    color: var(--v-text-dark-primary, #FFFFFF);

    &.pattern-suggestions__tabs__number--difference {
      background-color: var(--v-info-base);
    }
  }
}
</style>
