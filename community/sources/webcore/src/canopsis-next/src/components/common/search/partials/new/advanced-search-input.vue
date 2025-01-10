<template>
  <v-menu v-bind="menuProps" bottom>
    <template #activator="{ on }">
      <input
        v-if="active"
        v-model="inputValue"
        ref="inputEl"
        type="text"
        class="ml-1"
        autocomplete="off"
        v-on="on"
        @keydown="keydown"
        @focusout="focusout"
      >
    </template>
    <variables-list
      :items="items"
      children-key="items"
      return-object
      @input="select"
    />
  </v-menu>
</template>

<script>
import { computed, ref, watch, nextTick } from 'vue';

import VariablesList from '@/components/common/text-editor/variables-list.vue';

export default {
  components: { VariablesList },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    items: {
      type: Array,
      default: () => [],
    },
    closable: {
      type: Boolean,
      default: false,
    },
    itemText: {
      type: String,
      required: false,
    },
    itemValue: {
      type: String,
      required: false,
    },
    active: {
      type: Boolean,
      default: false,
    },
    allowText: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const inputEl = ref(null);
    const inputValue = ref('');

    const menuProps = computed(() => ({
      // openOnClick: false,
      disableKeys: true,
      closeOnContentClick: false,
      ignoreClickOutsideOnActivator: true,
      maxHeight: 304,
      nudgeBottom: 1,
      bottom: true,
      offsetY: true,
      transition: false,
    }));

    const chipText = computed(() => (props.itemText ? props.value[props.itemText] : props.value));

    const apply = (event) => {
      console.log(event);
    };
    const select = item => emit('input', item);

    const keydown = (event) => {
      console.log(event);
    };

    const mousedown = () => emit('mousedown');
    const close = () => emit('remove');
    const focusout = () => emit('focusout');

    watch(() => props.value, (newValue) => {
      if (!newValue.value) {
        return;
      }

      inputValue.value = newValue[props.itemText] ?? '';
    }, { immediate: true });

    watch(() => props.active, (active) => {
      if (!active) {
        return;
      }

      nextTick(() => inputEl.value?.focus());
    }, { immediate: true });

    return {
      inputEl,
      inputValue,
      menuProps,
      chipText,

      apply,
      select,
      keydown,
      mousedown,
      focusout,
      close,
    };
  },
};
</script>
