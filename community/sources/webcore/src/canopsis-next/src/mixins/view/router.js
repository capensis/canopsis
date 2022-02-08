import { ROUTES_NAMES } from '@/constants';

export const viewRouterMixin = {
  methods: {
    /**
     * Redirect to first view tab if exists
     *
     * @return {Promise<unknown>}
     */
    redirectToFirstTab() {
      return new Promise((resolve, reject) => {
        const { tabId } = this.$route.query;

        if (tabId || !this.view?.tabs?.length) {
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
