<template>
  <v-expansion-panels
    class="c-collapse-panel"
    accordion
    :style="panelStyle"
  >
    <v-expansion-panel class="c-collapse-panel__panel elevation-2">
      <v-expansion-panel-header :color="color">
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
        :class="panelContentClass"
        :style="panelContentStyle"
      >
        <v-card
          class="c-collapse-panel__card"
          flat
        >
          <v-card-text
            class="c-collapse-panel__card"
          >
            <slot />
          </v-card-text>
        </v-card>
      </v-expansion-panel-content>
    </v-expansion-panel>
  </v-expansion-panels>
</template>

<script>
import { Validator } from 'vee-validate';

import { validationChildrenMixin } from '@/mixins/form';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  mixins: [validationChildrenMixin],
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
  },
  computed: {
    panelStyle() {
      return { outlineColor: this.outlineColor };
    },

    panelContentClass() {
      return ['c-collapse-panel__content', { error: this.hasError }];
    },

    panelContentStyle() {
      return { backgroundColor: this.color };
    },

    hasError() {
      return this.error || this.hasChildrenError;
    },
  },
};
</script>

<style lang="scss">
.c-collapse-panel {
  --c-collapse-panel-border-radius: 5px;

  outline: 3px solid transparent;

  &__panel {
    border-radius: var(--c-collapse-panel-border-radius);
    overflow: hidden;
  }
}
</style>
