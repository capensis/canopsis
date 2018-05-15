import alarmModule from './alarm';
import alarmConventionModule from './alarmConvention';

export default {
  namespaced: true,
  modules: {
    alarm: alarmModule,
    alarmConvention: alarmConventionModule,
  },
};
