import dataModule from './data';
import testModule from './test';
import varsModule from './vars';
import validationModule from './validation';

export default {
  namespaced: true,
  modules: {
    data: dataModule,
    test: testModule,
    vars: varsModule,
    validation: validationModule,
  },
};
