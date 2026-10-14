<template>
  <c-text-field
    ref="field"
    :model-value="value"
    :label="label || $t('common.payload')"
    :disabled="disabled"
    :clearable="clearable"
    :name="name"
    :required="required"
    :max-length="maxLength"
    class="c-payload-text-field"
    @blur="handleBlur"
    @click="onSelectionChange"
    @keyup="onSelectionChange"
    @update:model-value="updateModel"
  >
    <template v-if="$slots.append || $slots['append-inner'] || variables.length" #append-inner>
      <slot name="append" />
      <slot name="append-inner" />

      <variables-menu
        v-if="variables.length"
        :items="variables"
        :visible="variablesShown"
        :value="variablesMenuValue"
        :position-x="variablesMenuPosition.x"
        :position-y="variablesMenuPosition.y"
        :location="variablesMenuLocation"
        ignore-click-outside
        show-value
        hide-empty-value
        clickable-parent
        @update:value="pasteVariable"
        @close="hideVariablesMenu"
      />
    </template>
    <template v-if="$slots['append-outer']" #append>
      <slot name="append-outer" />
    </template>
  </c-text-field>
</template>

<script>
import { payloadFieldMixin } from '@/mixins/payload/payload-field';

import VariablesMenu from '@/components/common/text-editor/variables-menu.vue';

export default {
  components: { VariablesMenu },
  mixins: [payloadFieldMixin],

  props: {
    name: {
      type: String,
      default: 'payload',
    },
    label: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    maxLength: {
      type: Number,
      required: false,
    },
  },
  data() {
    return {
      variablesMenuLocation: 'bottom start',
    };
  },
  methods: {
    showVariablesMenu() {
      this.variablesShown = true;

      this.$nextTick(() => {
        const input = this.getInputElement();

        if (!input) {
          return;
        }

        const menuMaxHeight = 300;
        const offset = 8;
        const { top, bottom, left } = input.getBoundingClientRect();
        const spaceBelow = window.innerHeight - bottom;
        const spaceAbove = top;
        const showBelow = spaceBelow >= menuMaxHeight || spaceBelow >= spaceAbove;

        this.variablesMenuPosition.x = left;
        this.variablesMenuPosition.y = showBelow ? bottom + offset : top - offset;
        this.variablesMenuLocation = showBelow ? 'bottom start' : 'top start';
      });
    },
  },
};
</script>
