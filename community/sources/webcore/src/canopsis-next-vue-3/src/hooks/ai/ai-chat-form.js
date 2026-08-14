import { computed, provide, ref } from 'vue';

/**
 * Temporary AI chat migration placeholder.
 *
 * The real AI chat layer is intentionally deferred until the end of the migration.
 * Keep this hook surface available for already migrated forms that register expand callbacks.
 *
 * @deprecated Replace with the full `../canopsis-next/src/hooks/ai/ai-chat-form.js` port when AI chat is migrated.
 * @param {Object} [options={}] Hook options kept for compatibility with the old API.
 * @returns {Object}
 */
export const useAiChatExpand = () => {};

/**
 * Temporary AI chat form placeholder.
 *
 * It preserves the public shape used by migrated modals while keeping the sidebar hidden and inactive.
 *
 * @deprecated Replace with the full AI chat form hook during the final AI chat migration pass.
 * @returns {{
 *   shown: import('vue').ComputedRef<boolean>,
 *   options: import('vue').ComputedRef<{
 *     bind: Object,
 *     on: Object,
 *   }>,
 *   llms: import('vue').Ref<Array>,
 *   showChat: Function,
 *   hideChat: Function,
 * }}
 */
export const useAiChatForm = () => {
  const llms = ref([]);
  const pending = ref(false);

  const noop = () => {};

  provide('$aiChat', {
    llms,
    updateJsonString: noop,
    registerExpandFunction: noop,
    unregisterExpandFunction: noop,
  });

  const options = computed(() => ({
    bind: {
      pending: pending.value,
      llms: llms.value,
    },
    on: {},
  }));

  return {
    shown: computed(() => false),
    options,
    llms,
    showChat: noop,
    hideChat: noop,
  };
};
