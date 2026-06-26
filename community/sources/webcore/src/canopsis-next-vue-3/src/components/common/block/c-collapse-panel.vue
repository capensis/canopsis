<script setup>
import { computed, ref, watch } from 'vue';

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

// const { hasChildrenError } = useValidationChildren();
const hasChildrenError = ref(false);

const panelStyle = computed(() => ({
  outlineColor: props.outlineColor,
}));

const hasError = computed(() => props.error || hasChildrenError.value);

const headerColor = computed(() => (hasError.value ? 'error' : undefined));

const headerStyle = computed(() => {
  if (hasError.value) {
    return {};
  }

  return {
    backgroundColor: props.color,
  };
});

watch(
  () => props.expanded,
  expanded => {
    localExpanded.value = expanded ? 0 : undefined;
  },
);
</script>

<template>
  <v-expansion-panels
    v-model="localExpanded"
    :style="panelStyle"
    :disabled="disabled"
    class="c-collapse-panel"
  >
    <v-expansion-panel class="c-collapse-panel__panel">
      <v-expansion-panel-title
        :color="headerColor"
        :style="headerStyle"
      >
        <slot name="header">
          <span class="text-white">
            {{ title }}
          </span>
        </slot>

        <template #actions>
          <slot name="actions">
            <v-icon
              :icon="icon"
              color="white"
            />
          </slot>
        </template>
      </v-expansion-panel-title>

      <v-expansion-panel-text eager>
        <v-card
          class="c-collapse-panel__card"
          elevation="0"
        >
          <v-card-text>
            <slot />
          </v-card-text>
        </v-card>
      </v-expansion-panel-text>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<style lang="scss">
.c-collapse-panel {
  --c-collapse-panel-border-radius: 5px;

  outline: 3px solid transparent;
  border-radius: var(--c-collapse-panel-border-radius) !important;

  &__panel .v-expansion-panel__shadow {
    overflow: hidden;
  }

  &__card.v-card {
    border-radius: 0;
  }
}
</style>
