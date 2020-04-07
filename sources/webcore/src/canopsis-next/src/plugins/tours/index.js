import VStep from 'vue-tour/src/components/VStep.vue';

import VTour from './components/v-tour.vue';

export default {
  install(Vue) {
    Vue.component(VTour.name, VTour);
    Vue.component(VStep.name, VStep);

    Vue.prototype.$tours = {};
  },
};
