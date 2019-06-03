/* eslint-disable global-require */

export default function featureRegister(Vue) {
  if (process.env.VUE_APP_SHOW_ALARM_SNOW_FEATURE) {
    const featureModule = require('@/test.js');

    Vue.use(featureModule.default);
  }

  console.log('required');
}
