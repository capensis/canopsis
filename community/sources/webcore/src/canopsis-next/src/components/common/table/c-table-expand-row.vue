<template>
  <tr v-if="shown" :class="rowClass">
    <td
      :colspan="colspan"
      :class="[
        cellClass,
        { 'c-table-expand-row__cell--closing': closing },
      ]"
      class="c-table-expand-row__cell"
    >
      <div
        ref="content"
        :class="{
          'c-table-expand-row__content-wrapper--opening': opening,
          'c-table-expand-row__content-wrapper--closing': closing,
        }"
        class="c-table-expand-row__content-wrapper"
      >
        <div :class="contentClass" class="c-table-expand-row__content">
          <slot />
        </div>
      </div>
    </td>
  </tr>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue';

const EXPAND_TRANSITION_DURATION = 260;

const props = defineProps({
  expanded: {
    type: Boolean,
    default: false,
  },
  colspan: {
    type: Number,
    required: true,
  },
  rowClass: {
    type: [String, Array, Object],
    required: false,
  },
  cellClass: {
    type: [String, Array, Object],
    required: false,
  },
  contentClass: {
    type: [String, Array, Object],
    required: false,
  },
});

const emit = defineEmits(['after-close']);

const shown = ref(props.expanded);
const opening = ref(false);
const closing = ref(false);
const timeout = ref(null);
const content = ref(null);

/**
 * Clears the pending expand/collapse timeout, if any.
 */
const clearExpandTimeout = () => {
  if (!timeout.value) {
    return;
  }

  clearTimeout(timeout.value);
  timeout.value = null;
};

/**
 * Shows the expand row and runs the open height animation.
 */
const show = () => {
  clearExpandTimeout();

  shown.value = true;
  closing.value = false;
  opening.value = true;

  nextTick(() => {
    const element = content.value;

    if (!element) {
      return;
    }

    element.style.height = '0px';

    requestAnimationFrame(() => {
      element.style.height = `${element.scrollHeight}px`;
    });
  });

  timeout.value = window.setTimeout(() => {
    const element = content.value;

    if (element) {
      element.style.height = 'auto';
    }

    opening.value = false;
    timeout.value = null;
  }, EXPAND_TRANSITION_DURATION);
};

/**
 * Hides the expand row and runs the close height animation.
 * Emits `after-close` once the collapse transition finishes.
 */
const hide = () => {
  if (!shown.value || closing.value) {
    return;
  }

  clearExpandTimeout();
  opening.value = false;
  closing.value = true;

  nextTick(() => {
    const element = content.value;

    if (!element) {
      return;
    }

    element.style.height = `${element.scrollHeight}px`;

    requestAnimationFrame(() => {
      element.style.height = '0px';
    });
  });

  timeout.value = window.setTimeout(() => {
    shown.value = false;
    closing.value = false;
    timeout.value = null;

    emit('after-close');
  }, EXPAND_TRANSITION_DURATION);
};

watch(
  () => props.expanded,
  (expanded) => {
    if (expanded) {
      show();

      return;
    }

    hide();
  },
);

onMounted(() => {
  if (props.expanded) {
    show();
  }
});

onBeforeUnmount(clearExpandTimeout);

defineExpose({
  show,
  hide,
});
</script>

<style lang="scss" scoped>
.c-table-expand-row {
  &__cell {
    height: 0;
    padding: 0 !important;
    border-bottom: 0 !important;
    line-height: 0;
  }

  &__content-wrapper {
    height: auto;
    overflow: hidden;
    opacity: 1;
    line-height: normal;
    transition:
      height 0.2s ease-out,
      opacity 0.2s ease-out;
  }

  &__content-wrapper--opening,
  &__content-wrapper--closing {
    will-change: height, opacity;
  }

  &__content-wrapper--closing {
    opacity: 0;
  }

  &__cell--closing {
    padding: 0 !important;
    border-bottom: 0 !important;
  }

  &__content {
    overflow: hidden;
  }
}
</style>
