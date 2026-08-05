import { nextTick, unref, watch } from 'vue';
import { useRouter } from 'vue-router';

/**
 * Watches for changes in the active notification ID and the items list.
 * When the active item is found in the items list, triggers the provided action and resets the router query.
 *
 * @param {Object} params
 * @param {import('vue').Ref<string|null>} params.activeId - Reactive reference to the active notification ID.
 * @param {import('vue').Ref<Array<Object>>} params.items - Reactive reference to the list of notification items.
 * @param {import('vue').Ref<Function|null>|Function|null} params.action -
 *   Action to execute when the active item is found.
 */
export const useNotificationActiveId = ({ activeId, items, action }) => {
  const router = useRouter();
  let unwatchItems = null;

  const watchItems = (unwrappedActiveId) => {
    unwatchItems?.();

    unwatchItems = watch(
      items,
      (newItems) => {
        const activeItem = newItems.find(item => item._id === unwrappedActiveId);

        if (activeItem) {
          unref(action)?.(activeItem);

          nextTick(() => {
            unwatchItems?.();
            unwatchItems = null;

            router.replace({ query: {} });
          });
        }
      },
      { immediate: true },
    );
  };

  watch(
    activeId,
    (newActiveId) => {
      if (newActiveId) {
        watchItems(newActiveId);
      } else {
        unwatchItems?.();
      }
    },
    { immediate: true },
  );
};
