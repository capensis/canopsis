const GridItem = () => import(/* webpackChunkName: "View" */ './components/grid-item.vue');
const GridLayout = () => import(/* webpackChunkName: "View" */ './components/grid-layout.vue');

export default {
  install(Vue) {
    Vue.component('grid-item', GridItem);
    Vue.component('grid-layout', GridLayout);
  },
};
