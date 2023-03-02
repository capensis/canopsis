<template lang="pug">
  v-expansion-panel.c-collapse-panel(:style="panelStyle")
    v-expansion-panel-content(
      :class="panelContentClass",
      :style="panelContentStyle",
      :lazy="lazy"
    )
      template(#actions="")
        v-icon(color="white") {{ icon }}
      template(#header="")
        slot(name="header")
          span.white--text {{ title }}
      v-card
        v-card-text
          slot
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
    lazy: {
      type: Boolean,
      default: false,
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
      return { error: this.hasError };
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
