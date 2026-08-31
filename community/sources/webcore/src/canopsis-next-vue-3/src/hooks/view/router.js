import { omit } from 'lodash';
import { useRoute, useRouter } from 'vue-router';

import { ROUTES_NAMES } from '@/constants';

/**
 * Hook for view routing functionality.
 *
 * @param {Object} options - Options object.
 * @param {import('vue').Ref<Object>} options.view - Reactive reference to the current view.
 * @param {import('vue').Ref<Object>} options.tab - Reactive reference to the current tab.
 * @returns {Object} View router methods.
 */
export const useViewRouter = ({ view, tab } = {}) => {
  const router = useRouter();
  const route = useRoute();

  /**
   * Redirect to home page if we are surfing on this view at the moment.
   *
   * @returns {Promise<unknown>|undefined}
   */
  const redirectToHomeIfCurrentRoute = () => {
    const { name, params = {} } = route;

    if (name === ROUTES_NAMES.view && params.id === view?.value?._id) {
      return router.push({ name: ROUTES_NAMES.home });
    }

    return undefined;
  };

  /**
   * Redirect to first view tab if exists.
   *
   * @returns {Promise<unknown>|undefined}
   */
  const redirectToFirstTab = () => {
    if (!view?.value?.tabs?.length) {
      return Promise.resolve();
    }

    return router.replace({
      query: {
        tabId: view.value.tabs[0]._id,
      },
    });
  };

  /**
   * Redirect to view root route (without tabId).
   *
   * @returns {Promise<unknown>}
   */
  const redirectToViewRoot = () => router.replace({
      query: omit(route.query, 'tabId'),
    });

  /**
   * Redirect to selected view and tab, if it's different then the view/tab we're actually on.
   *
   * @param {Object} payload - Redirect payload.
   * @param {string} payload.tabId - Selected tab id.
   * @param {string} payload.viewId - Selected view id.
   * @returns {Promise<unknown>|undefined}
   */
  const redirectToSelectedViewAndTab = ({ tabId, viewId }) => {
    if (tab?.value?._id === tabId) {
      return Promise.resolve();
    }

    return router.push({
      name: ROUTES_NAMES.view,
      params: { id: viewId },
      query: { tabId },
    });
  };

  return {
    router,
    route,
    redirectToHomeIfCurrentRoute,
    redirectToFirstTab,
    redirectToViewRoot,
    redirectToSelectedViewAndTab,
  };
};
