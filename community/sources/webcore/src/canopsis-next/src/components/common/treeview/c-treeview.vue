<template>
  <v-defaults-provider :defaults="buttonDefaults">
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
      class="c-treeview"
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
  </v-defaults-provider>
</template>

<script>
import { isString } from 'lodash';
import { computed, ref } from 'vue';
import { VTreeview } from 'vuetify/components';

const buttonDefaults = {
  VBtn: {
    size: 'x-small',
  },
};

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

  setup(props, { emit }) {
    const internalOpened = ref([]);
    const internalActivated = ref([]);
    const internalModelValue = ref([]);

    const resolvedItemTitle = computed(() => props.itemTitle ?? props.itemText ?? 'title');
    const resolvedItemValue = computed(() => props.itemValue ?? props.itemKey);
    const currentOpened = computed(() => props.opened ?? props.open ?? internalOpened.value);
    const currentActivated = computed(() => props.activated ?? props.active ?? internalActivated.value);
    const currentModelValue = computed(() => props.modelValue ?? props.value ?? internalModelValue.value);

    const getItemChildren = (item) => {
      if (typeof props.itemChildren === 'function') {
        return props.itemChildren(item);
      }

      if (Array.isArray(props.itemChildren)) {
        return props.itemChildren.reduce((acc, key) => acc?.[key], item);
      }

      return item?.[props.itemChildren];
    };

    const isLeaf = (item) => {
      if (!item || typeof item !== 'object') {
        return true;
      }

      if (isString(props.itemChildren) && Object.hasOwn(item, props.itemChildren)) {
        return false;
      }

      return !getItemChildren(item);
    };

    const getLegacySlotProps = slotProps => ({
      ...slotProps,
      leaf: isLeaf(slotProps.item),
    });

    const updateOpened = (opened) => {
      if (!props.opened && !props.open) {
        internalOpened.value = opened;
      }

      emit('update:opened', opened);
      emit('update:open', opened);
    };

    const updateActivated = (activated) => {
      if (!props.activated && !props.active) {
        internalActivated.value = activated;
      }

      emit('update:activated', activated);
      emit('update:active', activated);
    };

    const updateModelValue = (value) => {
      if (!props.modelValue && !props.value) {
        internalModelValue.value = value;
      }

      emit('update:modelValue', value);
      emit('input', value);
    };

    return {
      buttonDefaults,
      resolvedItemTitle,
      resolvedItemValue,
      currentOpened,
      currentActivated,
      currentModelValue,
      getLegacySlotProps,
      updateOpened,
      updateActivated,
      updateModelValue,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-treeview {
  :deep(.v-list-item) {
    min-height: 32px;
    padding-block: 0;
  }
}
</style>
