<template>
  <v-treeview
    v-bind="$attrs"
    :items="items"
    :item-title="resolvedItemTitle"
    :item-value="resolvedItemValue"
    :item-children="itemChildren"
    :open-on-click="openOnClick"
    :opened="currentOpened"
    :activated="currentActivated"
    :selected="selected"
    :model-value="currentModelValue"
    @update:opened="updateOpened"
    @update:activated="updateActivated"
    @update:selected="$emit('update:selected', $event)"
    @update:model-value="updateModelValue"
    @click:open="$emit('click:open', $event)"
    @click:select="$emit('click:select', $event)"
  >
    <template v-if="$slots.prepend" #prepend="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="prepend" />
    </template>

    <template v-if="$slots.label" #title="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="label" />
    </template>

    <template v-if="$slots.append" #append="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="append" />
    </template>

    <template v-if="$slots.item" #item="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="item" />
    </template>

    <template v-if="$slots.header" #header="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="header" />
    </template>

    <template v-if="$slots.footer" #footer="slotProps">
      <slot v-bind="getLegacySlotProps(slotProps)" name="footer" />
    </template>

    <template v-if="$slots['no-data']" #no-data>
      <slot name="no-data" />
    </template>
  </v-treeview>
</template>

<script>
import { VTreeview } from 'vuetify/components';

export default {
  components: { VTreeview },

  inheritAttrs: false,

  props: {
    items: {
      type: Array,
      default: () => [],
    },
    itemKey: {
      type: [String, Array, Function],
      default: 'id',
    },
    itemText: {
      type: [String, Array, Function],
      default: undefined,
    },
    itemTitle: {
      type: [String, Array, Function],
      default: undefined,
    },
    itemValue: {
      type: [String, Array, Function],
      default: undefined,
    },
    itemChildren: {
      type: [String, Array, Function],
      default: 'children',
    },
    openOnClick: {
      type: Boolean,
      default: true,
    },
    open: {
      type: Array,
      default: undefined,
    },
    opened: {
      type: Array,
      default: undefined,
    },
    active: {
      type: Array,
      default: undefined,
    },
    activated: {
      type: Array,
      default: undefined,
    },
    selected: {
      type: Array,
      default: undefined,
    },
    value: {
      type: Array,
      default: undefined,
    },
    modelValue: {
      type: Array,
      default: undefined,
    },
  },

  data() {
    return {
      internalOpened: [],
      internalActivated: [],
      internalModelValue: [],
    };
  },

  emits: [
    'click:open',
    'click:select',
    'input',
    'update:active',
    'update:activated',
    'update:modelValue',
    'update:open',
    'update:opened',
    'update:selected',
  ],

  computed: {
    resolvedItemTitle() {
      return this.itemTitle ?? this.itemText ?? 'title';
    },

    resolvedItemValue() {
      return this.itemValue ?? this.itemKey;
    },

    currentOpened() {
      return this.opened ?? this.open ?? this.internalOpened;
    },

    currentActivated() {
      return this.activated ?? this.active ?? this.internalActivated;
    },

    currentModelValue() {
      return this.modelValue ?? this.value ?? this.internalModelValue;
    },
  },

  methods: {
    getItemChildren(item) {
      if (typeof this.itemChildren === 'function') {
        return this.itemChildren(item);
      }

      if (Array.isArray(this.itemChildren)) {
        return this.itemChildren.reduce((acc, key) => acc?.[key], item);
      }

      return item?.[this.itemChildren];
    },

    isLeaf(item) {
      if (!item || typeof item !== 'object') {
        return true;
      }

      if (typeof this.itemChildren === 'string' && Object.hasOwn(item, this.itemChildren)) {
        return false;
      }

      return !this.getItemChildren(item);
    },

    getLegacySlotProps(slotProps) {
      return {
        ...slotProps,
        leaf: this.isLeaf(slotProps.item),
      };
    },

    updateOpened(opened) {
      if (!this.opened && !this.open) {
        this.internalOpened = opened;
      }

      this.$emit('update:opened', opened);
      this.$emit('update:open', opened);
    },

    updateActivated(activated) {
      if (!this.activated && !this.active) {
        this.internalActivated = activated;
      }

      this.$emit('update:activated', activated);
      this.$emit('update:active', activated);
    },

    updateModelValue(value) {
      if (!this.modelValue && !this.value) {
        this.internalModelValue = value;
      }

      this.$emit('update:modelValue', value);
      this.$emit('input', value);
    },
  },
};
</script>
