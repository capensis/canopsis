import { ROUTES_NAMES } from '@/constants';

/**
 * TODO: use it
 */
export default {
  methods: {
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
