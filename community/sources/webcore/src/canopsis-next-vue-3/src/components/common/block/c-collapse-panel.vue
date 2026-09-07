<template>
  <v-expansion-panels
    v-model="localExpanded"
    :style="panelStyle"
    :disabled="disabled"
    class="c-collapse-panel"
  >
    <v-expansion-panel class="c-collapse-panel__panel">
      <v-expansion-panel-title :style="headerStyle" class="c-collapse-panel__title">
        <div class="c-collapse-panel__header">
          <slot name="header">
            <span class="c-collapse-panel__header-text">
              {{ title }}
            </span>
          </slot>
        </div>

        <template #actions>
          <slot name="actions">
            <v-icon :icon="icon" />
          </slot>
        </template>
      </v-expansion-panel-title>

      <v-expansion-panel-text eager>
        <v-card class="c-collapse-panel__card" elevation="0">
          <v-card-text>
            <slot />
          </v-card-text>
        </v-card>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script setup>
import { computed, ref, watch } from 'vue';

import { CSS_COLORS_VARS } from '@/config';

import { useValidationChildren } from '@/hooks/validator/validation-children';

const props = defineProps({
  title: {
    type: String,
    default: '',
  },
  color: {
    type: String,
    default: 'rgb(128, 128, 128)',
  },
  outlineColor: {
    type: String,
    default: undefined,
  },
  icon: {
    type: String,
    default: 'expand_more',
  },
  error: {
    type: Boolean,
    default: false,
  },
  expanded: {
    type: Boolean,
    default: false,
  },
  disabled: {
    type: Boolean,
    default: false,
  },
});

const localExpanded = ref(props.expanded ? 0 : undefined);

const { hasChildrenError } = useValidationChildren();

const hasError = computed(() => props.error || hasChildrenError.value);

const panelStyle = computed(() => ({
  outlineColor: hasError.value ? CSS_COLORS_VARS.error : props.outlineColor,
}));

const headerStyle = computed(() => ({
  backgroundColor: props.color,
}));

watch(
  () => props.expanded,
  (expanded) => {
    localExpanded.value = expanded ? 0 : undefined;
  },
);
</script>

<style lang="scss">
.c-collapse-panel {
  --c-collapse-panel-border-radius: 5px;

  outline: 3px solid transparent;
  border-radius: var(--c-collapse-panel-border-radius) !important;

  &__panel {
    overflow: hidden;

    .v-expansion-panel__shadow {
      overflow: hidden;
    }
  }

  &__card.v-card {
    border-radius: 0;
  }

  &__title {
    .v-expansion-panel-title__overlay {
      pointer-events: none;
    }

    .v-icon {
      color: white;
    }
  }

  &__header {
    position: relative;
    z-index: 1;
    flex: 1 1 auto;
    min-width: 0;
  }

  &__header-text {
    color: white;
  }
}
</style>
