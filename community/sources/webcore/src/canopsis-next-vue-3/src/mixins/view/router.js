import { omit } from 'lodash';

import { ROUTES_NAMES } from '@/constants';

export const viewRouterMixin = {
  methods: {
    /**
     * Redirect to home page if we are surfing on this view at the moment.
     *
     * @returns {Promise<unknown>|undefined}
     */
    redirectToHomeIfCurrentRoute() {
      const { name, params = {} } = this.$route;

      if (name === ROUTES_NAMES.view && params.id === this.view?._id) {
        return this.$router.push({ name: ROUTES_NAMES.home });
      }

      return undefined;
    },

    /**
     * Redirect to first view tab if exists.
     *
     * @returns {Promise<unknown>}
     */
    redirectToFirstTab() {
      if (!this.view?.tabs?.length) {
        return Promise.resolve();
      }

      return this.$router.replace({
        query: {
          tabId: this.view.tabs[0]._id,
        },
      });
    },

    /**
     * Redirect to view root route (without tabId).
     *
     * @returns {Promise<unknown>}
     */
    redirectToViewRoot() {
      return this.$router.replace({
        query: omit(this.$route.query, 'tabId'),
      });
    },

    /**
     * Redirect to selected view and tab, if it's different then the view/tab we're actually on.
     *
     * @param {Object} payload - Redirect payload.
     * @param {string} payload.tabId - Selected tab id.
     * @param {string} payload.viewId - Selected view id.
     * @returns {Promise<unknown>}
     */
    redirectToSelectedViewAndTab({ tabId, viewId }) {
      if (this.tab?._id === tabId) {
        return Promise.resolve();
      }

      return this.$router.push({
        name: ROUTES_NAMES.view,
        params: { id: viewId },
        query: { tabId },
      });
    },
  },
};
