<template>
  <v-expansion-panels
    class="c-collapse-panel"
    accordion
    :style="panelStyle"
  >
    <v-expansion-panel>
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
        <v-card>
          <v-card-text>
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
      default: 'grey',
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
  border-radius: 5px;
  overflow: hidden;
  outline: 3px solid transparent;
}
</style>
