/**
 * This store is for the naming/color/icon convention we give to the status:states of an alarm
 * Like that we have the same convention for all the app
 */
import getProp from 'lodash/get';

export default {
  namespaced: true,
  state: {
    statesConvention: {
      0: {
        color: 'green',
        text: 'ok',
        icon: 'assistant_photo',
      },
      1: {
        color: 'yellow darken-1',
        text: 'minor',
        icon: 'assistant_photo',
      },
      2: {
        color: 'orange',
        text: 'major',
        icon: 'assistant_photo',
      },
      3: {
        color: 'red',
        text: 'critical',
        icon: 'assistant_photo',
      },
    },
    statusConvention: {
      0: {
        color: 'black',
        text: 'off',
        icon: 'keyboard_arrow_up',
      },
      1: {
        color: 'grey',
        text: 'ongoing',
        icon: 'keyboard_arrow_up',
      },
      2: {
        color: 'yellow darken-1',
        text: 'stealthy',
        icon: 'keyboard_arrow_up',
      },
      3: {
        color: 'orange',
        text: 'flapping',
        icon: 'keyboard_arrow_up',
      },
      4: {
        color: 'red',
        text: 'cancelled',
        icon: 'keyboard_arrow_up',
      },
    },
  },
  getters: {
    getStateAlarmConvention: state => value => conventionProp =>
      getProp(state.statesConvention, `${value}.${conventionProp}`),

    getStatusAlarmConvention: state => value => conventionProp =>
      getProp(state.statusConvention, `${value}.${conventionProp}`),
  },
};
