import Vue from 'vue';
import VueI18n from 'vue-i18n';
import Vuetify from 'vuetify';
import { shallow } from '@vue/test-utils';

import i18n from '@/i18n';
import HelloWorld from '@/components/hello-world.vue';

describe('HelloWorld.vue', () => {
  beforeAll(() => {
    Vue.use(VueI18n);
    Vue.use(Vuetify);
  });

  it('renders props.msg when passed', () => {
    const msg = 'new message';
    const wrapper = shallow(HelloWorld, {
      propsData: { msg },
      i18n,
    });

    expect(wrapper.text()).toMatch(msg);
  });
});
