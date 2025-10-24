import { omit } from 'lodash';
import { useRouter, useRoute } from 'vue-router/composables';

import { ROUTES_NAMES } from '@/constants';

/**
 * Hook for view routing functionality
 * Replaces the viewRouterMixin with Composition API
 *
 * @param {Object} options - Options object
 * @param {import('vue').Ref<Object>} options.view - Reactive reference to the current view
 * @param {import('vue').Ref<Object>} options.tab - Reactive reference to the current tab
 * @returns {Object} An object containing view router methods
 */
export const useViewRouter = ({ view, tab } = {}) => {
  const router = useRouter();
  const route = useRoute();

  /**
   * Redirect to home page if we are surfing on this view at the moment
   */
  const redirectToHomeIfCurrentRoute = () => {
    const { name, params = {} } = route;

    if (name === ROUTES_NAMES.view && params.id === view?.value?._id) {
      router.push({ name: ROUTES_NAMES.home });
    }
  };

  /**
   * Redirect to first view tab if exists
   *
   * @return {Promise<unknown>}
   */
  const redirectToFirstTab = () => new Promise((resolve, reject) => {
    if (!view?.value?.tabs?.length) {
      return resolve();
    }

    return router.replace({
      query: {
        tabId: view.value.tabs[0]._id,
      },
    }, resolve, reject);
  });

  /**
   * Redirect to view root route (without tabId)
   *
   * @return {Promise<unknown>}
   */
  const redirectToViewRoot = () => new Promise((resolve, reject) => router.replace({
    query: omit(route.query, 'tabId'),
  }, resolve, reject));

  /**
   * Redirect to selected view and tab, if it's different then the view/tab we're actually on
   *
   * @param {string} tabId
   * @param {string} viewId
   * @return {Promise<unknown>}
   */
  const redirectToSelectedViewAndTab = ({ tabId, viewId }) => new Promise((resolve, reject) => {
    if (tab?.value?._id === tabId) {
      return resolve();
    }

    return router.push({
      name: ROUTES_NAMES.view,
      params: { id: viewId },
      query: { tabId },
    }, resolve, reject);
  });

  return {
    router,
    route,
    redirectToHomeIfCurrentRoute,
    redirectToFirstTab,
    redirectToViewRoot,
    redirectToSelectedViewAndTab,
  };
};
