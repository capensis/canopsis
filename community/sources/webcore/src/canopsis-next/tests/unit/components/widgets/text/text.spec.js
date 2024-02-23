import AxiosMockAdapter from 'axios-mock-adapter';
import axios from 'axios';

import { flushPromises, generateRenderer } from '@unit/utils/vue';

import TextWidget from '@/components/widgets/text/text.vue';
import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';

const stubs = {
  'c-runtime-template': CRuntimeTemplate,
  'c-compiled-template': CCompiledTemplate,
};

describe('text', () => {
  const axiosMockAdapter = new AxiosMockAdapter(axios);

  const snapshotFactory = generateRenderer(TextWidget, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  test('Renders `text` with default template', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          parameters: {
            template: '<div><b>bold template</b></div>',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `text` with request template', async () => {
    const url = 'https://jsonplaceholder.typicode.com/todos';

    axiosMockAdapter
      .onPost(url)
      .reply(201, {
        userId: '1',
        title: 'test',
        completed: false,
        id: 201,
      });

    const wrapper = snapshotFactory({
      propsData: {
        widget: {
          parameters: {
            template: `
              {{#request
               method="post"
               url="${url}"
               variable="post"
               headers='{ "Content-Type": "application/json" }'
               data='{ "userId": "1", "title": "test", "completed": false }'}}
                {{#each post}}
                    <li><strong>{{@key}}</strong>: {{this}}</li>
                {{/each}}
              {{/request}}
            `,
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
  });
});
