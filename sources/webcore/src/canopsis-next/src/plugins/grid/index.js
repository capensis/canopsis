import GridItem from './components/grid-item.vue';
import GridLayout from './components/grid-layout.vue';

export default {
  install(Vue) {
    Vue.component('grid-item', GridItem);
    Vue.component('grid-layout', GridLayout);
  },
};
