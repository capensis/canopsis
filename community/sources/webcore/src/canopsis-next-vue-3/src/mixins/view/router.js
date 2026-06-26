import { ROUTES_NAMES } from '@/constants';

export const viewRouterMixin = {
  methods: {
    redirectToHomeIfCurrentRoute() {
      if (this.$route.name === ROUTES_NAMES.view && this.$route.params.id === this.view._id) {
        this.$router.push({ name: ROUTES_NAMES.home });
      }
    },
  },
};
