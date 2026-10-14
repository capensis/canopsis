import { omit } from 'lodash';

import { ROUTES_NAMES } from '@/constants';

export const viewRouterMixin = {
  methods: {
    /**
     * Redirect to home page if we are surfing on this view at the moment
     */
    redirectToHomeIfCurrentRoute() {
      const { name, params = {} } = this.$route;

      if (name === ROUTES_NAMES.view && params.id === this.view._id) {
        this.$router.push({ name: ROUTES_NAMES.home });
      }
    },

    /**
     * Redirect to first view tab if exists
     *
     * @return {Promise<unknown>}
     */
    redirectToFirstTab() {
      return new Promise((resolve, reject) => {
        if (!this.view?.tabs?.length) {
          return resolve();
        }

        return this.$router.replace({
          query: {
            tabId: this.view.tabs[0]._id,
          },
        }, resolve, reject);
      });
    },

    /**
     * Redirect to view root route (without tabId)
     *
     * @return {Promise<unknown>}
     */
    redirectToViewRoot() {
      return new Promise((resolve, reject) => this.$router.replace({
        query: omit(this.$route.query, 'tabId'),
      }, resolve, reject));
    },

    /**
     * Redirect to selected view and tab, if it's different then the view/tab we're actually on
     *
     * @param {string} tabId
     * @param {string} viewId
     * @return {Promise<unknown>}
     */
    redirectToSelectedViewAndTab({ tabId, viewId }) {
      return new Promise((resolve, reject) => {
        if (this.tab._id === tabId) {
          return resolve();
        }

        return this.$router.push({
          name: ROUTES_NAMES.view,
          params: { id: viewId },
          query: { tabId },
        }, resolve, reject);
      });
    },
  },
};
