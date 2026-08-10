<template>
  <c-combobox
    ref="field"
    :model-value="value"
    :items="variablesItems"
    :label="label || $t('common.payload')"
    :disabled="disabled"
    :return-object="false"
    :menu="variablesShown"
    :menu-props="menuProps"
    :clearable="clearable"
    :name="name"
    :required="required"
    class="c-payload-text-field"
    no-filter
    @blur="handleBlur"
    @click="onSelectionChange"
    @keyup="onSelectionChange"
    @update:menu="updateMenu"
    @update:model-value="updateModel"
  >
    <template v-if="$slots.append || $slots['append-inner']" #append-inner>
      <slot name="append" />
      <slot name="append-inner" />
    </template>
    <template v-if="$slots['append-outer']" #append>
      <slot name="append-outer" />
    </template>
    <template #prepend-item>
      <variables-list
        :value="variablesMenuValue"
        :items="variables"
        children-key="variables"
        show-value
        hide-empty-value
        clickable-parent
        @update:value="pasteVariable"
      />
    </template>
  </c-combobox>
</template>

<script>
import VariablesList from '@/components/common/text-editor/variables-list.vue';
import { payloadFieldMixin } from '@/mixins/payload/payload-field';

export default {
  components: { VariablesList },
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
  },

  computed: {
    variablesItems() {
      return [];
    },

    menuProps() {
      return {
        closeOnContentClick: false,
        minWidth: 200,
      };
    },
  },

  methods: {
    updateMenu(value) {
      if (!value) {
        this.hideVariablesMenu();
      }
    },
  },
};
</script>
