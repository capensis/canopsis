import { generateShallowRenderer } from '@unit/utils/vue';

import CRequestTextInformation from '@/components/common/request/c-request-text-information.vue';
import RequestInformationContent from '@/components/common/request/partials/request-information-content.vue';

const selectInformationContent = wrapper => wrapper.findAllComponents(RequestInformationContent);

describe('c-request-text-information', () => {
  const factory = generateShallowRenderer(CRequestTextInformation);

  it('renders request information', () => {
    const request = 'POST /tickets HTTP/1.1\r\nContent-Type: application/json';
    const wrapper = factory({ propsData: { value: request } });
    const contents = selectInformationContent(wrapper);

    expect(contents).toHaveLength(1);
    expect(contents[0].props('text')).toBe(request);
  });

  it('renders request and response information', () => {
    const request = 'POST /tickets HTTP/1.1';
    const response = 'HTTP/1.1 201 Created';
    const wrapper = factory({ propsData: { value: `${request}\n\n${response}` } });
    const contents = selectInformationContent(wrapper);

    expect(contents).toHaveLength(2);
    expect(contents[0].props('text')).toBe(request);
    expect(contents[1].props('text')).toBe(response);
  });
});

describe('request-information-content', () => {
  const factory = generateShallowRenderer(RequestInformationContent, {
    stubs: {
      'c-json-treeview': {
        props: ['json'],
        template: '<div class="c-json-treeview-stub" />',
      },
    },
  });

  it('renders HTTP information and a JSON body', () => {
    const body = '{"ticket":"42"}';
    const wrapper = factory({
      propsData: {
        text: `HTTP/1.1 201 Created\r\nContent-Type: application/json\r\n\r\n${body}`,
      },
    });

    expect(wrapper.text()).toContain('HTTP/1.1 201 Created');
    expect(wrapper.text()).toContain('Content-Type: application/json');
    expect(wrapper.vm.request.body).toBe(body);
    expect(wrapper.vm.request.isJsonBody).toBe(true);
    expect(wrapper.find('.c-json-treeview-stub').exists()).toBe(true);
  });

  it('renders a plain-text body', () => {
    const wrapper = factory({
      propsData: {
        text: 'HTTP/1.1 500 Internal Server Error\r\n\r\nConnection failed',
      },
    });

    expect(wrapper.text()).toContain('Connection failed');
    expect(wrapper.find('.c-json-treeview-stub').exists()).toBe(false);
  });
});
