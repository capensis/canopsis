<template>
  <v-expansion-panels
    v-model="localExpanded"
    :style="panelStyle"
    class="c-collapse-panel elevation-2"
    accordion
  >
    <v-expansion-panel class="c-collapse-panel__panel">
      <v-expansion-panel-header :color="headerColor">
        <slot name="header">
          <span class="white--text">{{ title }}</span>
        </slot>
        <template #actions="">
          <slot name="actions">
            <v-icon color="white">
              {{ icon }}
            </v-icon>
          </slot>
        </template>
      </v-expansion-panel-header>
      <v-expansion-panel-content
        :style="panelContentStyle"
        eager
      >
        <v-card
          class="c-collapse-panel__card"
          flat
        >
          <v-card-text>
            <slot />
          </v-card-text>
        </v-card>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
import { ref, computed } from 'vue';
import { Validator } from 'vee-validate';

import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  props: {
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
      required: false,
    },
    icon: {
      type: String,
      default: '$vuetify.icons.expand',
    },
    error: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const localExpanded = ref(props.expanded ? 0 : null);

    const { hasChildrenError } = useValidationChildren();

    const panelStyle = computed(() => ({ outlineColor: props.outlineColor }));

    const panelContentStyle = computed(() => ({ backgroundColor: props.color }));

    const hasError = computed(() => props.error || hasChildrenError.value);

    const headerColor = computed(() => (hasError.value ? 'error' : props.color));

    return {
      localExpanded,
      panelStyle,
      headerColor,
      panelContentStyle,
      hasError,
    };
  },
};
</script>

<style lang="scss">
.c-collapse-panel {
  --c-collapse-panel-border-radius: 5px;

  outline: 3px solid transparent;
  border-radius: var(--c-collapse-panel-border-radius) !important;

  &__panel {
    overflow: hidden;
  }

  &__card.v-card {
    border-radius: 0;
  }
}
</style>
