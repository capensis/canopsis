import FeatureModal from '@/feature-modal.vue';

export default {
  install(Vue) {
    console.log('INSTALLED');
    Vue.prototype.$test = 'TEST321';
    Vue.component('fFeatureModal', FeatureModal);
  },
};
