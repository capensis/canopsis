import { generateRenderer } from '@unit/utils/vue';

import ThirdPartyLogin from '@/components/other/login/partials/third-party-login.vue';

describe('cas-login', () => {
  const snapshotFactory = generateRenderer(ThirdPartyLogin);

  it('Renders `cas-login` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `cas-login` with title and links', () => {
    const wrapper = snapshotFactory({
      propsData: {
        title: 'Title',
        links: [
          { title: 'Link title', href: 'Link href' },
          { title: 'Link title 2', href: 'Link href 2' },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
