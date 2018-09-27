import { DsDaysView } from 'dayspan-vuetify/src/components';

export default {
  extends: DsDaysView,
  data() {
    return {
      hours: [
        '   ', '01h', '02h', '03h', '04h', '05h', '06h', '07h', '08h', '09h', '10h', '11h',
        '12h', '13h', '14h', '15h', '16h', '17h', '18h', '19h', '20h', '21h', '22h', '23h',
      ],
    };
  },
  render(...args) {
    return DsDaysView.render.apply(this, args);
  },
};
